package main

import (
	"SongAlgoWeb/chat"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/chat", chat.RequestHandler)

	err := router.Run("localhost:8080")
	if err != nil {
		fmt.Printf("Error running gin router: %s", err)
		return
	}
}
