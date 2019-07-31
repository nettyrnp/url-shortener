package service

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var (
	// Allowed actions and their short versions
	FullToShortMap = map[string]string{
		"login":        "lgn",
		"authenticate": "thntct",
		"register":     "rgstr",
		"discover":     "dscvr",
	}
	ShortToFullMap = reverse(FullToShortMap)

	vowels = strings.Split("aeoiu", "")

	upgrader = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
	}
)

func reverse(m1 map[string]string) map[string]string {
	m2 := map[string]string{}
	for key, value := range m1 {
		m2[value] = key
	}
	return m2
}

type app struct {
	socket *websocket.Conn
	sendCh chan []byte
}

func (a *app) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

func (a *app) read() {
	defer a.socket.Close()
	for {
		_, msg, err := a.socket.ReadMessage()
		if err != nil {
			return
		}
		a.sendCh <- shorten(msg)
	}
}

func (a *app) write() {
	defer a.socket.Close()
	for msg := range a.sendCh {
		err := a.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

func shorten(msg []byte) []byte {
	parts := strings.SplitN(string(msg), "/", 2)
	head := parts[0]
	tail := parts[1]
	short := deleteVowels(head)
	return []byte(short + "/" + tail)
}

func deleteVowels(str string) string {
	for _, s := range vowels {
		str = strings.ReplaceAll(str, s, "")
	}
	return str
}

func NewApp() *app {
	return &app{}
}
