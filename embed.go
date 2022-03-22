package tinymce

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/justinas/nosurf"
)

//go:embed www/*
var Content embed.FS

//go:generate go run --tags generate ./gen/gen.go

//templatestring is index.html from the embedded www directory
var templateFile, _ = Content.Open("www/index.html")
var templateString, _ = ioutil.ReadAll(templateFile)

var templ = template.Must(template.New("t1").Parse(string(templateString)))

type EditorView struct {
	Hostname string
	Port     int
	WorkDir  string
	File     string
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

func (e EditorView) Origin() string {
	return fmt.Sprintf("http://%s:%d", e.Hostname, e.Port)
}

func (e EditorView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// sanitize the path
	cleanpath := "www" + path.Clean(r.URL.Path)
	if (cleanpath == "www/") || (cleanpath == "www") {
		cleanpath = "www/index.html"
	}

	//w.Header().Set("Access-Control-Allow-Headers", "content-type")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Origin", e.Origin())
	w.Header().Add("Vary", "Origin")
	w.Header().Add("Vary", "Access-Control-Request-Method")
	w.Header().Add("Vary", "Access-Control-Request-Headers")
	w.Header().Add("Access-Control-Allow-Origin", e.Origin()) //"*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	token := nosurf.Token(r)
	w.Header().Add("Set-Cookie", fmt.Sprintf("token=%s; Path=/", token))

	// if we encounter the path in the embedded Content FS, serve it
	if file, err := Content.Open(cleanpath); err == nil {
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", ContentType(cleanpath, bytes))
		if cleanpath == "www/index.html" {
			lb, err := ioutil.ReadFile(filepath.Join(e.WorkDir, e.File))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			context := make(map[string]string)
			context["token"] = token
			context["filename"] = filepath.Join(e.WorkDir, e.File)
			context["content"] = string(lb)
			log.Println("Serving:", cleanpath)
			if err := templ.Execute(w, context); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//w.Write(bytes)
			return
		}
		w.Write(bytes)
	} else {
		apipath := strings.TrimPrefix(cleanpath, "www/")
		log.Println("API:", apipath)
		context := make(map[string]string)
		context["token"] = token
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch apipath {
		case "save":
			context["body"] = string(bytes)
			log.Println("Save:", context["body"])
			data := make(map[string]interface{})
			err = json.Unmarshal(bytes, &data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fp := data["path"].(string)
			fb := []byte(data["text"].(string))
			log.Println("Save:\n\t", string(fb), "\n\t\tas", fp)
			// redirect home
			err = SaveFileOnDisk(fp, fb)
			http.Redirect(w, r, "/", http.StatusFound)
		case "download":
			context["body"] = string(bytes)
			log.Println("Download:", context["body"])
			w.Write(bytes)
		case "load":
			// load the content and refresh the page
			context["body"] = string(bytes)
			log.Println("Load:", context["body"])
			data := make(map[string]interface{})
			err = json.Unmarshal(bytes, &data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Println("Load:", data)
			fp := filepath.Join(e.WorkDir, data["path"].(string))
			//fb := []byte(data["text"].(string))
			fb, err := LoadFileOnDisk(fp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// redirect home
			http.Redirect(w, r, "/", http.StatusFound)
			w.Write(fb)
		default:
			// serve any files we find in the WorkDir on the filesystem
			if file, err := os.Open(filepath.Join(e.WorkDir, apipath)); err == nil {
				bytes, err := ioutil.ReadAll(file)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", ContentType(cleanpath, bytes))
				w.Write(bytes)
			} else {
				http.NotFound(w, r)
			}

		}
	}
	return

}

func Serve(host, dir, file string, port int) error {
	log.Println("Serving:", host, port)
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	return http.ListenAndServe(addr, EditorView{Hostname: host, Port: port, WorkDir: dir, File: file})
	//return http.ListenAndServe(addr, nosurf.New(EditorView{Hostname: host, Port: port, WorkDir: dir, File: file}))
}
