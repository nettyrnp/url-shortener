package router

import (
	"fmt"
	"github.com/nettyrnp/url-shortener/log"
	"github.com/nettyrnp/url-shortener/storage"
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

		if full, ok := storage_service.ShortToFullMap[action]; ok { // When action is a short version of an allowed action
			newURL := "/v1/" + full + "/" + username + tail
			http.Redirect(w, r, newURL, http.StatusFound)

		} else if _, ok := storage_service.FullToShortMap[action]; ok { // When action is a full allowed action
			logger.Infof("User %v tried to %v", username, action)

		} else { // When action is unknown
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
	})
	t.templ.Execute(w, r)
}
