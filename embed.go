package tinymce

import (
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/justinas/nosurf"
)

//go:embed www/*
var Content embed.FS

//go:generate go run --tags generate ./gen/gen.go

type EditorView struct {
}

func ContentType(path string, bytes []byte) string {
	// get the file extension
	ext := strings.ToLower(filepath.Ext(path))
	// map the extension to the content type
	switch ext {
	case ".js":
		return "text/javascript"
	case ".css":
		return "text/css"
	case ".html":
		return "text/html"
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".jpg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	default:
		return http.DetectContentType(bytes)
	}
}

func (e EditorView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// sanitize the path
	filepath := "www" + path.Clean(r.URL.Path)
	if (filepath == "www/") || (filepath == "www") {
		filepath = "www/index.html"
	}
	// if we encounter the path in the embedded Content FS, serve it
	if file, err := Content.Open(filepath); err == nil {
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// write the content-type header
		w.Header().Set("Content-Type", ContentType(filepath, bytes))
		// write the content
		w.Write(bytes)
	} else {
		apipath := strings.TrimPrefix(filepath, "www/")
		// if we encounter the path in the embedded Content FS, serve it
		switch apipath {
		case "save":
			// save the content
		case "load":
			// load the content and refresh the page
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	return

}

func Serve(host string, port int) error {
	log.Println("Serving:", host, port)
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	return http.ListenAndServe(addr, nosurf.New(EditorView{}))
}
