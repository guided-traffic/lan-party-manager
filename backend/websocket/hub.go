package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// MessageType defines the type of WebSocket message
type MessageType string

const (
	// MessageTypeVoteReceived is sent when a user receives a vote
	MessageTypeVoteReceived MessageType = "vote_received"
	// MessageTypeNewVote is sent to all clients when any vote is created (for timeline)
	MessageTypeNewVote MessageType = "new_vote"
	// MessageTypeUserJoined is sent when a new user joins
	MessageTypeUserJoined MessageType = "user_joined"
	// MessageTypeSettingsUpdate is sent when admin changes settings
	MessageTypeSettingsUpdate MessageType = "settings_update"
	// MessageTypeError is sent when an error occurs
	MessageTypeError MessageType = "error"
)

// Message represents a WebSocket message
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// VotePayload contains vote information for notifications
type VotePayload struct {
	VoteID        uint64 `json:"vote_id"`
	FromUserID    uint64 `json:"from_user_id"`
	FromUsername  string `json:"from_username"`
	FromAvatar    string `json:"from_avatar"`
	ToUserID      uint64 `json:"to_user_id"`
	ToUsername    string `json:"to_username"`
	ToAvatar      string `json:"to_avatar"`
	AchievementID string `json:"achievement_id"`
	Achievement   string `json:"achievement_name"`
	IsPositive    bool   `json:"is_positive"`
	CreatedAt     string `json:"created_at"`
}

// SettingsPayload contains settings information for broadcasts
type SettingsPayload struct {
	CreditIntervalMinutes int `json:"credit_interval_minutes"`
	CreditMax             int `json:"credit_max"`
}

// Client represents a connected WebSocket client
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   uint64
	steamID  string
	username string
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients by user ID
	clients map[uint64]*Client

	// All clients for broadcast
	allClients map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast to all clients
	broadcast chan []byte

	// Send to specific user
	sendToUser chan *UserMessage

	mutex sync.RWMutex
}

// UserMessage is a message targeted at a specific user
type UserMessage struct {
	UserID  uint64
	Message []byte
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint64]*Client),
		allClients: make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		sendToUser: make(chan *UserMessage),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.userID] = client
			h.allClients[client] = true
			h.mutex.Unlock()
			log.Printf("WebSocket: Client connected - User %d (%s)", client.userID, client.username)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.allClients[client]; ok {
				delete(h.allClients, client)
				delete(h.clients, client.userID)
				close(client.send)
				log.Printf("WebSocket: Client disconnected - User %d (%s)", client.userID, client.username)
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.allClients {
				select {
				case client.send <- message:
				default:
					// Client send buffer full, close connection
					close(client.send)
					delete(h.allClients, client)
					delete(h.clients, client.userID)
				}
			}
			h.mutex.RUnlock()

		case userMsg := <-h.sendToUser:
			h.mutex.RLock()
			if client, ok := h.clients[userMsg.UserID]; ok {
				select {
				case client.send <- userMsg.Message:
				default:
					// Client send buffer full
					close(client.send)
					delete(h.allClients, client)
					delete(h.clients, client.userID)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// BroadcastVote sends a new vote notification to all clients
func (h *Hub) BroadcastVote(payload *VotePayload) {
	msg := Message{
		Type:    MessageTypeNewVote,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("WebSocket: Failed to marshal broadcast message: %v", err)
		return
	}

	h.broadcast <- data
}

// NotifyVoteReceived sends a notification to the user who received a vote
func (h *Hub) NotifyVoteReceived(toUserID uint64, payload *VotePayload) {
	msg := Message{
		Type:    MessageTypeVoteReceived,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("WebSocket: Failed to marshal notification message: %v", err)
		return
	}

	h.sendToUser <- &UserMessage{
		UserID:  toUserID,
		Message: data,
	}
}

// GetConnectedUserCount returns the number of connected users
func (h *Hub) GetConnectedUserCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.allClients)
}

// IsUserConnected checks if a specific user is connected
func (h *Hub) IsUserConnected(userID uint64) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// BroadcastSettingsUpdate sends settings update to all connected clients
func (h *Hub) BroadcastSettingsUpdate(payload *SettingsPayload) {
	msg := Message{
		Type:    MessageTypeSettingsUpdate,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("WebSocket: Failed to marshal settings message: %v", err)
		return
	}

	h.broadcast <- data
	log.Printf("WebSocket: Broadcasted settings update to all clients")
}
