package websocket

import (
	"simpleMessenger/internal/service"
	"sync"
)

type Hub struct {
	clients        map[*Client]bool
	register       chan *Client
	unregister     chan *Client
	mu             sync.RWMutex
	messageService *service.MessageService
	chatService    *service.ChatService
}

func NewHub(messageService *service.MessageService, chatService *service.ChatService) *Hub {
	return &Hub{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		messageService: messageService,
		chatService:    chatService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}

func (h *Hub) SendToUser(userID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.userID == userID {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
			break
		}
	}
}

func (h *Hub) SendToChat(chatID, senderID uint, message []byte) error {
	participants, err := h.chatService.GetUsersInChat(chatID)
	if err != nil {
		return err
	}

	for _, userID := range participants {
		if userID == senderID {
			//continue
		}
		h.SendToUser(userID, message)
	}
	return nil
}
