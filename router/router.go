package router

import (
	"fmt"
	http_service "github.com/nettyrnp/url-shortener/http_service"
	"github.com/nettyrnp/url-shortener/log"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

var (
	logger = log.GetLogger()
)

type SvcHandler struct{}

func (t *SvcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.URL.Path, "/", 5)
	action := parts[2]
	username := ""
	tail := ""
	if len(parts) > 3 {
		username = parts[3]
		if len(parts) > 4 {
			tail = "/" + parts[4]
		}
		if full, ok := http_service.ShortToFullMap[action]; ok {
			shortURL := "/v1/" + full + "/" + username + tail
			w.Header().Set("Location", shortURL)
			w.WriteHeader(http.StatusFound)
		} else if _, ok := http_service.FullToShortMap[action]; ok {
			logger.Infof("User %v tried to %v", username, action)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Page '%v' not found", r.URL.Path)
		}
	}
}

type RootHandler struct {
	once     sync.Once
	Filename string
	templ    *template.Template
}

func (t *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("", t.Filename)))
		//t.templ = template.Must(template.ParseFiles(filepath.Join("http_service/templates", t.Filename)))
	})
	t.templ.Execute(w, r)
}
