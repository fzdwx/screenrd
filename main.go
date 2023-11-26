package main

import (
	"flag"
)

func main() {
	port := flag.Int("p", 8080, "port")
	flag.Parse()

	fs := GetFrontFs()
	err := NewServer(*port, fs).
		ListenAndServe()
	if err != nil {
		panic(err)
	}
}
