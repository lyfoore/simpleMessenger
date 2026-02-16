package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf("ServeWs called, headers: %v", r.Header)
	val := r.Context().Value("user_id")
	log.Printf("user_id in context: value=%v, type=%T", val, val)

	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		log.Printf("user_id not found or invalid type: %v", val)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error while upgrading connection: %v", err)
		return
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}

	log.Printf("hello\n")

	client.hub.register <- client

	go client.writePump()
	log.Printf("writing gorutine for %v is running", userID)
	go client.readPump()
	log.Printf("writing gorutine for %v is running", userID)
}
