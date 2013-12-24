package hub

import (
	"fmt"
	"github.com/gorilla/websocket"
)

// Goroutine wrapper for a websocket connection. Anything sent on the `send` channel
// will be written to the websocket.
type Connection struct {
	ws   *websocket.Conn
	send chan []byte
}

// Creates and returns a new Connection object
func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{send: make(chan []byte, 256), ws: ws}
}

// Listens forever on the websocket, performing actions as needed.
func (c *Connection) Reader() {
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			fmt.Printf("Error in websocket read: %s\n", err.Error())
			break
		}
		fmt.Printf("Message received: %s\n", msg)
		// TODO: Do something with the msg
	}
	c.ws.Close()
}

// Writes anything on the send channel to the websocket.
func (c *Connection) Writer() {
	for msg := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Printf("Error in websocket write: %s\n", err.Error())
			break
		}
	}
	c.ws.Close()
}
