package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"simpleMessenger/internal/model"
	"simpleMessenger/internal/service"
	"time"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID uint
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024
)

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: %v", err)
			}
			break
		}
		c.handleMessage(message)
	}
}

func (c *Client) handleMessage(data []byte) {
	message := &model.Message{}
	if err := json.Unmarshal(data, &message); err != nil {
		log.Println("invalid message:", err)
		return
	}

	message.UserID = c.userID

	req := &service.SendMessageRequest{
		Text:   message.Text,
		UserID: message.UserID,
		ChatID: message.ChatID,
	}

	msgResp, err := c.hub.messageService.SendMessage(req)
	if err != nil {
		log.Printf("send message error: %v", err)
		return
	}

	msgBytes, err := json.Marshal(msgResp)
	if err != nil {
		log.Printf("failed to marshal message: %v", err)
		return
	}

	if err := c.hub.SendToChat(message.ChatID, message.UserID, msgBytes); err != nil {
		log.Printf("failed to send to chat: %v", err)
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
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
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
