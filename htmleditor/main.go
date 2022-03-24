package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	tinymce "github.com/eyedeekay/go-htmleditor"
)

func main() {
	host := flag.String("host", "127.0.0.1", "Hostname to listen on")
	port := flag.Int("port", 8082, "Port to listen on")
	dir := flag.String("dir", "./www", "Directory to serve files from")
	file := flag.String("file", "index.html", "File to serve")
	flag.Parse()
	directory, err := filepath.Abs(*dir)
	if err != nil {
		panic(err)
	}
	index := filepath.Join(directory, *file)
	if _, err := os.Stat(index); os.IsNotExist(err) {
		if err := ioutil.WriteFile(index, []byte(tinymce.MinHtmlDoc), 0644); err != nil {
			panic(err)
		}
	}
	if err := tinymce.Serve(*host, *dir, *file, *port); err != nil {
		panic(err)
	}
}
