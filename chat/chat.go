package chat

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

// upgrader upgrade http request to socket connection.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

var clients = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{}

// RequestHandler handles incoming WebSocket connections
func RequestHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error upgrading to WebSocket: %s", err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			// Error handling needed

		}
	}(ws)

	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	for {
		_, messageBytes, err := ws.ReadMessage()
		if err != nil {
			return
		}

		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			fmt.Printf("Error unmarshalling message: %v\n", err)
			continue
		}

		fmt.Printf("%+v\n", msg)

		switch msg.Type {
		case Default:
			broadcast(msg.Text, ws)
		case TypingNotification:
			broadcastTypingNotification(msg.User, ws)
		}
	}
}

// broadcast sends the message to all clients except the sender
func broadcast(text string, sender *websocket.Conn) {
	message := Message{
		Type: Default,
		User: "user",
		Text: text,
	}
	messageBytes, _ := json.Marshal(message)

	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		if client == sender {
			continue
		}
		if err := client.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			// Handle error...
		}
	}
}

// broadcastTypingNotification sends a typing notification to all clients except the sender
func broadcastTypingNotification(user string, sender *websocket.Conn) {
	notification := Message{
		Type: TypingNotification,
		User: user,
	}
	notificationBytes, _ := json.Marshal(notification)

	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		if client == sender {
			continue
		}
		if err := client.WriteMessage(websocket.TextMessage, notificationBytes); err != nil {
			// Handle error...
		}
	}
}
