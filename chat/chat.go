package chat

import (
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
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			return
		}

		fmt.Printf("%s sent: %s\n", ws.RemoteAddr(), string(p))

		mutex.Lock()
		for client := range clients {
			if err := client.WriteMessage(messageType, p); err != nil {
				err := client.Close()
				if err != nil {
					// Error handling needed
				}
				delete(clients, client)
			}
		}
		mutex.Unlock()

		if err := ws.WriteMessage(messageType, p); err != nil {
			mutex.Lock()
			delete(clients, ws)
			mutex.Unlock()
			err := ws.Close()
			if err != nil {
				// Error handling needed
			}
			break
		}
	}
}
