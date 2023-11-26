package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Server struct {
	port      int
	fs        http.FileSystem
	maxMemory int64
}

func NewServer(port int, fs http.FileSystem) *Server {
	return &Server{
		port:      port,
		fs:        fs,
		maxMemory: 1024 * 1024 * 1024,
	}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/api/download" {
		s.download(writer, request)
		return
	}

	if request.URL.Path == "/" {
		request.URL.Path = "/front"
	}

	http.FileServer(s.fs).ServeHTTP(writer, request)
}

func (s *Server) ListenAndServe() error {
	url := fmt.Sprintf("http://localhost:%d", s.port)
	fmt.Println(url)

	address := fmt.Sprintf(":%d", s.port)
	return http.ListenAndServe(address, s)
}

func (s *Server) download(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(s.maxMemory)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(writer, err)
		return
	}

	file, _, err := request.FormFile("file")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(writer, err)
		return
	}

	value := request.FormValue("corpInfo")
	crop := &Corp{}
	err = json.Unmarshal([]byte(value), crop)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(writer, err)
		return
	}

	filename := time.Now().Format("20060102150405")
	tempFile, err := os.CreateTemp("", filename+".webm")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(writer, err)
		return
	}
	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name())
	}()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(writer, err)
		return
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(writer, err)
		return
	}

	outFilename := fmt.Sprintf("%s_out.mp4", tempFile.Name())
	err = crop.Crop(tempFile, outFilename)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(writer, err)
		return
	}

	http.ServeFile(writer, request, outFilename)
}

func writeError(writer http.ResponseWriter, err error) {
	writer.WriteHeader(http.StatusInternalServerError)
	_, _ = fmt.Fprintln(writer, err)
}
