package service

import (
	"encoding/json"
	"log"
	"github.com/Order-Payment-Go-Microservice/message-service/internal/model"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[string][]*websocket.Conn
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string][]*websocket.Conn),
	}
}

func (h *Hub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = append(h.clients[userID], conn)
	log.Printf("User %s registered. Total connections for user: %d", userID, len(h.clients[userID]))
}

func (h *Hub) Unregister(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, ok := h.clients[userID]; ok {
		for i, c := range conns {
			if c == conn {
				h.clients[userID] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		if len(h.clients[userID]) == 0 {
			delete(h.clients, userID)
		}
	}
}

func (h *Hub) BroadcastMessage(msg *model.Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	receiverID := msg.ReceiverID.String()
	if conns, ok := h.clients[receiverID]; ok {
		h.sendToConnections(conns, msg)
	}

	senderID := msg.SenderID.String()
	if conns, ok := h.clients[senderID]; ok {
		h.sendToConnections(conns, msg)
	}
}

func (h *Hub) sendToConnections(conns []*websocket.Conn, msg *model.Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Error sending message to connection: %v", err)
		}
	}
}
