package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	port := flag.Int("p", 8080, "port")
	flag.Parse()

	fs := http.FS(ff)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/api/download" {
			download(writer, request)
			return
		}

		if request.URL.Path == "/" {
			request.URL.Path = "/front"
		}
		http.FileServer(fs).ServeHTTP(writer, request)
	})

	fmt.Println(fmt.Sprintf("Please visit http://localhost:%d", *port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		panic(err)
	}
}

type CorpInfo struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func (i CorpInfo) isEmpty() bool {
	return i.X == 0 && i.Y == 0 && i.Width == 0 && i.Height == 0
}

func download(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(1024 * 1024 * 1024)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(writer, err)
		return
	}

	file, _, err := request.FormFile("file")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(writer, err)
		return
	}

	value := request.FormValue("corpInfo")
	corpInfo := CorpInfo{}
	err = json.Unmarshal([]byte(value), &corpInfo)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(writer, err)
		return
	}

	filename := time.Now().Format("20060102150405")
	tempFile, err := os.CreateTemp("", filename+".webm")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(writer, err)
		return
	}
	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name())
	}()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(writer, err)
		return
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(writer, err)
		return
	}

	outFilename := fmt.Sprintf("%s_out.mp4", tempFile.Name())
	err = crop(tempFile, outFilename, corpInfo)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(writer, err)
		return
	}

	http.ServeFile(writer, request, outFilename)
}

func crop(tempFile *os.File, outFilename string, corpInfo CorpInfo) error {
	args := []string{"-i", tempFile.Name()}
	if corpInfo.isEmpty() == false {
		args = append(args, "-filter:v", fmt.Sprintf("crop=%d:%d:%d:%d", int(corpInfo.Width), int(corpInfo.Height), int(corpInfo.X), int(corpInfo.Y)))
	}
	args = append(args, outFilename)
	return exec.Command("ffmpeg", args...).Run()
}
