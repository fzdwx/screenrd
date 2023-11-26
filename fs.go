package main

import (
	"embed"
	"net/http"
)

//go:embed front
var frontFs embed.FS

func GetFrontFs() http.FileSystem {
	return http.FS(frontFs)
}
