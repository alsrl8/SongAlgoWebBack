package user

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

	decodedBytes, err := base64.StdEncoding.DecodeString(params.Code)
	if err != nil {
		fmt.Println("Error decoding Base64 string:", err)
		return
	}
	decodedString := string(decodedBytes)
	fmt.Println("Decoded String:", decodedString)

	fmt.Printf("Code: %+v\n", params)
	c.JSON(http.StatusOK, gin.H{"Response": "Ok"})
}
