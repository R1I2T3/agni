package inapp

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
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

func (h *Hub) Register(UserID string, conn *websocket.Conn) *Client {
	fmt.Println(UserID)
	c := &Client{conn: conn, Write: make(chan []byte, 32), Read: make(chan []byte, 32), userID: UserID}
	h.mu.Lock()
	if _, ok := h.clients[UserID]; !ok {
		h.clients[UserID] = make(map[*Client]bool)
	}
	h.clients[UserID][c] = true
	h.mu.Unlock()
	go c.writePump()
	fmt.Printf("started ws with someone\n")
	return c
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		fmt.Printf("close 1\n")
	}()
	for {
		select {
		case msg, ok := <-c.Write:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				fmt.Printf("%s", msg)
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Printf("error", "%s", msg)
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {

				return
			}
		}
	}
}

func (c *Client) ReadPump() {
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
			fmt.Printf("close 2 \n")
			return
		}
		fmt.Println(string(message))

		select {
		case c.Read <- message:
		default:

			DefaultHub.Unregister(c)
			c.conn.Close()
			fmt.Printf("close 3\n")
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
	fmt.Printf("close 4\n")
}

func (h *Hub) BroadcastToUser(userID string, payload interface{}) {
	h.mu.RLock()
	fmt.Println(userID)
	clients := h.clients[userID]
	h.mu.RUnlock()
	if clients == nil {
		fmt.Println("No clients found by UserId", userID)
		return
	}
	b, err := json.Marshal(payload)
	if err != nil {
		log.Printf("inapp: marshal payload: %v", err)
		return
	}
	for c := range clients {
		select {
		case c.Write <- b:
		default:
			h.Unregister(c)
		}
	}
}
