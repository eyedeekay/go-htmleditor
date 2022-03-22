package main

import (
	"flag"

	tinymce "github.com/eyedeekay/go-htmleditor"
)

func main() {
	host := flag.String("host", "127.0.0.1", "Hostname to listen on")
	port := flag.Int("port", 8082, "Port to listen on")
	dir := flag.String("dir", "./www", "Directory to serve files from")
	flag.Parse()
	if err := tinymce.Serve(*host, *dir, *port); err != nil {
		panic(err)
	}
}
