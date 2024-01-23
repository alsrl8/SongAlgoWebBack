package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
)

type AuthParameters struct {
	Code string `json:"code"`
}

func AuthHandler(c *gin.Context) {
	var params AuthParameters
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("%s\n", params.Code)

	clientId := "297952994958-e2mg151ea9g89lqm339cqo2hlljfup6j.apps.googleusercontent.com"
	clientSecret := "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"

	data := url.Values{}
	data.Set("code", params.Code)
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", "http://localhost:3000")
	data.Set("grant_type", "authorization_code")

	// Make the POST request
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to request token"})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	// Here, 'body' contains the access token response
	// You can unmarshal this response into a struct and use it as needed
	fmt.Println("Response from Google:", string(body))

	c.JSON(http.StatusOK, gin.H{"Response": "Ok", "TokenResponse": string(body)})
}
