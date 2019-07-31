package main

import (
	"flag"
	"fmt"
	"github.com/nettyrnp/url-shortener/service"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	addr := flag.String("addr", "", "Address to bind server on")
	flag.Parse()

	http.Handle("/", &indexHandler{filename: "index.html"})
	http.Handle("/v1/", &mainHandler{})
	http.Handle("/app", service.NewApp())

	log.Printf("Listening on %v\n", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type mainHandler struct{}

func (t *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.URL.Path, "/", 5)
	action := parts[2]
	username := ""
	tail := ""
	if len(parts) > 3 {
		username = parts[3]
		if len(parts) > 4 {
			tail = "/" + parts[4]
		}
		if full, ok := service.ShortToFullMap[action]; ok {
			shortURL := "/v1/" + full + "/" + username + tail
			w.Header().Set("Location", shortURL)
			w.WriteHeader(http.StatusFound)
		} else if _, ok := service.FullToShortMap[action]; ok {
			log.Printf("User %v tried to %v", username, action)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Page '%v' not found", r.URL.Path)
		}
	}
}

type indexHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}
