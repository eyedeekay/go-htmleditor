package tinymce

import (
	"embed"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

//go:embed www/*
var Content embed.FS

//go:generate go run --tags generate ./gen/gen.go

type EditorView struct {
}

func (e EditorView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// sanitize the path
	path := path.Clean(r.URL.Path)
	log.Println("Serving:", path)
	// if we encounter the path in the embedded Content FS, serve it
	if file, err := Content.Open(path); err == nil {
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// write the content-type header
		w.Header().Set("Content-Type", http.DetectContentType(bytes))
		// write the content
		w.Write(bytes)
		return
	}

}

func Serve() error {
	return http.ListenAndServe("127.0.0.1:8081", EditorView{})
}
