package main

import (
	"SongAlgoWeb/chat"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/chat", chat.RequestHandler)

	//err := router.Run("0.0.0.0:8080")

	certFile := "/etc/letsencrypt/live/server.songmingi.com/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/server.songmingi.com/privkey.pem"
	err := router.RunTLS("0.0.0.0:8080", certFile, keyFile) // RunTLS Test

	if err != nil {
		fmt.Printf("Error running gin router: %s", err)
		return
	}
}
