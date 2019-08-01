package http_service

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/nettyrnp/url-shortener/config"
	"log"
	"net/http"
	"strings"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var (
	vowels = strings.Split("aeoiu", "")

	upgrader = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
	}
)

func NewHTTPService(ctx context.Context, conf config.HTTPConfig) (*HTTPService, error) {
	return &HTTPService{HTTPConfig: conf}, nil
}

type HTTPService struct {
	config.HTTPConfig
	socket *websocket.Conn
	sendCh chan []byte
}

func (a *HTTPService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	sock, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	a.socket = sock
	a.sendCh = make(chan []byte, messageBufferSize)
	defer close(a.sendCh)

	go a.read()

	a.write()
}

func (a *HTTPService) read() {
	defer a.socket.Close()
	for {
		_, msg, err := a.socket.ReadMessage()
		if err != nil {
			return
		}
		a.sendCh <- shorten(msg)
	}
}

func (a *HTTPService) write() {
	defer a.socket.Close()
	for msg := range a.sendCh {
		err := a.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

func shorten(url []byte) []byte {
	parts := strings.SplitN(string(url), "/", 2)
	head := parts[0]
	tail := parts[1]
	return []byte(deleteVowels(head) + "/" + tail)
}

func deleteVowels(str string) string {
	for _, s := range vowels {
		str = strings.ReplaceAll(str, s, "")
	}
	return str
}
