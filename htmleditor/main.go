package main

import (
	"flag"

	tinymce "github.com/eyedeekay/go-htmleditor"
)

func main() {
	host := flag.String("host", "localhost", "Hostname to listen on")
	port := flag.Int("port", 8081, "Port to listen on")
	flag.Parse()
	if err := tinymce.Serve(*host, *port); err != nil {
		panic(err)
	}
}
