package main

import (
	tinymce "github.com/eyedeekay/go-htmleditor"
)

func main() {
	if err := tinymce.Serve(); err != nil {
		panic(err)
	}
}
