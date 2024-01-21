package main

import (
	"SongAlgoWeb/chat"
	"SongAlgoWeb/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/chat", chat.RequestHandler)

	var err error
	if utils.IsDevelopmentMode() {
		err = router.Run("0.0.0.0:8080")
	} else {
		certFile := "/home/mingi4754song/cert/fullchain.pem"
		keyFile := "/home/mingi4754song/cert/privkey.pem"
		err = router.RunTLS("0.0.0.0:443", certFile, keyFile)
	}

	if err != nil {
		fmt.Printf("Error running gin router: %s", err)
		return
	}
}
