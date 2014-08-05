package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func App(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/app.html")
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a websocket
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c := NewConnection(ws)
	Register(c)
	defer Unregister(c)
	go c.Writer()
	c.Reader()
}
