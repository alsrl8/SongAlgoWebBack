package main

import (
	"SongAlgoWeb/chat"
	"SongAlgoWeb/user"
	"SongAlgoWeb/utils"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust this to be more restrictive if needed
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/auth", user.AuthHandler)
	router.GET("/chat", chat.RequestHandler)

	if utils.IsDevelopmentMode() {
		runDevMode(router)
	} else {
		runProdMode(router)
	}

	fmt.Println("SongAlgoWeb Server is Running... ")
}

func runDevMode(router *gin.Engine) {
	fmt.Println("Running Dev Mode")
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Error running gin router in development mode: %v", err)
	}
}

func runProdMode(router *gin.Engine) {
	fmt.Println("Running Prod Mode")

	certFile, err := utils.GetEnv("TLS_CERT_FULLCHAIN_PATH")
	if err != nil {
		log.Fatalln("Error getting TLS certificate full chain path")
	}

	keyFile, err := utils.GetEnv("TLS_CERT_PRIVKEY_PATH")
	if err != nil {
		log.Fatalln("Error getting TLS private key path")
	}

	if err := router.RunTLS("0.0.0.0:443", certFile, keyFile); err != nil {
		log.Fatalf("Error running gin router in production mode: %v", err)
	}
}
