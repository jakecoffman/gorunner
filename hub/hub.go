package hub

import (
	"fmt"
)

var h Hub

type Hub struct {
	connections map[*Connection]bool
	register    chan *Connection
	unregister  chan *Connection
	refresh     chan bool
	onRefresh   func() []byte
}

func NewHub(f func() []byte) {
	h = Hub{
		refresh:     make(chan bool),
		onRefresh:   f,
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
	}
}

func Register(c *Connection) {
	h.register <- c
}

func Unregister(c *Connection) {
	h.unregister <- c
}

func Refresh() {
	h.refresh <- true
}

func Run() {
	for {
		select {
		case c := <-h.register:
			fmt.Println("Connect")
			h.connections[c] = true
			bytes := h.onRefresh()
			c.send <- bytes
		case c := <-h.unregister:
			fmt.Println("Disconnect")
			delete(h.connections, c)
			close(c.send)
		case <-h.refresh:
			fmt.Println("Refreshing")
			bytes := h.onRefresh()
			for c := range h.connections {
				select {
				case c.send <- bytes:
				default:
					delete(h.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
		}
	}
}
