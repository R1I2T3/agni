package inapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	conn   *websocket.Conn
	Write  chan []byte
	Read   chan []byte
	userID string
}

type Hub struct {
	mu          sync.RWMutex
	clients     map[string]map[*Client]bool
	rdb         *redis.Client
	ctx         context.Context
	subscribers map[string]*redis.PubSub
	cancelFuncs map[string]context.CancelFunc // Store cancel functions
	subMutex    sync.RWMutex
}

var DefaultHub *Hub

// InitializeHub sets up the hub with Redis client
func InitializeHub(ctx context.Context, rdb *redis.Client) {
	DefaultHub = &Hub{
		clients:     make(map[string]map[*Client]bool),
		rdb:         rdb,
		ctx:         ctx,
		subscribers: make(map[string]*redis.PubSub),
		cancelFuncs: make(map[string]context.CancelFunc),
	}
}

func (h *Hub) Register(userID string, conn *websocket.Conn) *Client {
	c := &Client{
		conn:   conn,
		Write:  make(chan []byte, 32),
		Read:   make(chan []byte, 32),
		userID: userID,
	}

	h.mu.Lock()
	isFirstClient := false
	if _, ok := h.clients[userID]; !ok {
		h.clients[userID] = make(map[*Client]bool)
		isFirstClient = true
	}
	h.clients[userID][c] = true
	h.mu.Unlock()

	// If this is the first client for this user, subscribe to their channel
	if isFirstClient {
		h.subscribeToUser(userID)
	}

	go c.writePump()
	log.Printf("✅ Client registered: %s (total clients for user: %d)", userID, len(h.clients[userID]))
	return c
}

// subscribeToUser creates a Redis subscription for this specific user
func (h *Hub) subscribeToUser(userID string) {
	h.subMutex.Lock()
	defer h.subMutex.Unlock()

	// Check if already subscribed
	if _, exists := h.subscribers[userID]; exists {
		return
	}

	// Subscribe to this user's specific channel: inapp:broadcast:app_id:user_id
	channel := fmt.Sprintf("%s%s", BroadcastChannelPrefix, userID)
	pubsub := h.rdb.Subscribe(h.ctx, channel)

	// Create a cancellable context for this subscription
	ctx, cancel := context.WithCancel(h.ctx)

	h.subscribers[userID] = pubsub
	h.cancelFuncs[userID] = cancel

	log.Printf("📡 Subscribed to: %s", channel)

	// Start listener with cancellable context
	go h.listenToChannel(ctx, userID, pubsub)
}

func (h *Hub) listenToChannel(ctx context.Context, userID string, pubsub *redis.PubSub) {
	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			// Context cancelled - stop this goroutine
			log.Printf("🛑 Stopping listener for: %s", userID)
			return

		case msg, ok := <-ch:
			if !ok {
				// Channel closed
				log.Printf("🔌 Channel closed for: %s", userID)
				return
			}

			var payload map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				log.Printf("❌ Invalid payload: %v", err)
				continue
			}

			h.BroadcastToUser(userID, payload)
		}
	}
}

// unsubscribeFromUser removes Redis subscription when no clients remain
func (h *Hub) unsubscribeFromUser(userID string) {
	h.subMutex.Lock()
	defer h.subMutex.Unlock()

	if pubsub, exists := h.subscribers[userID]; exists {
		pubsub.Close()
		delete(h.subscribers, userID)
		log.Printf("🔕 Unsubscribed from Redis channel for: %s", userID)
	}

	if cancel, exists := h.cancelFuncs[userID]; exists {
		cancel()
		delete(h.cancelFuncs, userID)
	}
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
				log.Printf("❌ Write error: %v", err)
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
			return
		}

		log.Printf("📥 Message from client: %s", string(message))

		select {
		case c.Read <- message:
		default:
			DefaultHub.Unregister(c)
			c.conn.Close()
			return
		}
	}
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	clientsRemaining := 0
	if clients, ok := h.clients[c.userID]; ok {
		delete(clients, c)
		clientsRemaining = len(clients)
		if clientsRemaining == 0 {
			delete(h.clients, c.userID)
		}
	}
	h.mu.Unlock()

	close(c.Write)
	c.conn.Close()

	// If no more clients for this user, unsubscribe from Redis
	if clientsRemaining == 0 {
		h.unsubscribeFromUser(c.userID)
	}

	log.Printf("❌ Client unregistered: %s (remaining: %d)", c.userID, clientsRemaining)
}

func (h *Hub) BroadcastToUser(userID string, payload interface{}) {
	h.mu.RLock()
	clients := h.clients[userID]
	h.mu.RUnlock()

	if clients == nil {
		log.Printf("⚠️  No clients found for: %s", userID)
		return
	}

	b, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ Marshal error: %v", err)
		return
	}

	sent := 0
	for c := range clients {
		select {
		case c.Write <- b:
			sent++
		default:
			h.Unregister(c)
		}
	}

	log.Printf("📨 Sent to %d/%d clients for user: %s", sent, len(clients), userID)
}
