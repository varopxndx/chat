package handler

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 10000
)

var newline = []byte{'\n'}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client to handle websocket connection
type Client struct {
	conn      *websocket.Conn
	webSocket *Socket
	send      chan []byte
}

// ProcessMessage prepare the client to process message
func (h Handler) ProcessMessage(g *gin.Context) {
	conn, err := wsUpgrader.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		log.Default().Printf("websocket connection error: %v", err)
		return
	}

	client := newClient(conn, h.webSocket)
	go client.writePump()
	go client.readPump()

	h.webSocket.register <- client
}

func newClient(conn *websocket.Conn, webSocket *Socket) *Client {
	return &Client{
		conn:      conn,
		webSocket: webSocket,
		send:      make(chan []byte, 256),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.disconnect()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		c.webSocket.broadcast <- jsonMessage
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) disconnect() {
	c.webSocket.unregister <- c
	close(c.send)
	c.conn.Close()
}
