package inapp

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	Write  chan []byte
	Read   chan []byte
	userID string
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]bool
}

var DefaultHub = &Hub{
	clients: make(map[string]map[*Client]bool),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// implement origin check in prod
		return true
	},
}

func (h *Hub) Register(UserID string, conn *websocket.Conn) *Client {
	c := &Client{conn: conn, Write: make(chan []byte, 32), Read: make(chan []byte, 32), userID: UserID}
	h.mu.Lock()
	if _, ok := h.clients[UserID]; !ok {
		h.clients[UserID] = make(map[*Client]bool)
	}
	h.clients[UserID][c] = true
	h.mu.Unlock()
	go c.writePump()
	go c.readPump()
	return c
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Write:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	const (
		pongWait       = 60 * time.Second
		maxMessageSize = 512
	)

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			DefaultHub.Unregister(c)
			c.conn.Close()
			return
		}

		// Example handling: echo message back to this client.
		// Adjust behavior as needed (broadcast, process, etc.).
		select {
		case c.Read <- message:
		default:
			// send buffer full, unregister client
			DefaultHub.Unregister(c)
			c.conn.Close()
			return
		}
	}
}
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	if clients, ok := h.clients[c.userID]; ok {
		delete(clients, c)
		if len(clients) == 0 {
			delete(h.clients, c.userID)
		}
	}
	h.mu.Unlock()
	close(c.Write)
	c.conn.Close()
}
