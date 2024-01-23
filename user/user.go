package user

import (
	"SongAlgoWeb/utils"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type AuthParameters struct {
	Code string `json:"code"`
}

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}

type ClientInfo struct {
	Name string `json:"name"`
}

func getGoogleAuthClientId() (string, error) {
	clientId := os.Getenv("GOOGLE_AUTH_CLIENT_ID")
	if clientId == "" {
		return "", errors.New("no client id for requesting google authentication")
	}
	return clientId, nil
}

func getGoogleAuthClientSecret() (string, error) {
	clientSecret := os.Getenv("GOOGLE_AUTH_CLIENT_SECRET")
	if clientSecret == "" {
		return "", errors.New("no client secret for requesting google authentication")
	}
	return clientSecret, nil
}

func getRedirectUri() string {
	if utils.IsDevelopmentMode() {
		return "http://localhost:3000"
	} else {
		return "https://songmingi.com"
	}
}

func getTokenResponseFromResponseBody(body []byte) (*GoogleTokenResponse, error) {
	var tokenResponse GoogleTokenResponse
	err := json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, errors.New("failed to unmarshal response body")
	}
	return &tokenResponse, nil
}

func getGoogleAuthClaimsFromTokenResponse(tokenResponse *GoogleTokenResponse) (*jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified((*tokenResponse).IDToken, jwt.MapClaims{})
	if err != nil {
		log.Printf("Error parsing JWT token: %v", err)
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("Error asserting type of JWT claims")
		return nil, errors.New("invalid JWT claims type")
	}

	return &claims, nil
}

func googleAuthenticate(code string) ([]byte, error) {
	clientId, err := getGoogleAuthClientId()
	if err != nil {
		return nil, err
	}
	clientSecret, err := getGoogleAuthClientSecret()
	if err != nil {
		return nil, err
	}
	redirectUri := getRedirectUri()

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectUri)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return nil, errors.New("failed to request token")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Can't close response body. %s", err.Error())
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body")
	}

	return body, nil
}

func getClientInfoFromClaims(claims *jwt.MapClaims) (*ClientInfo, error) {
	nameInterface, ok := (*claims)["name"]
	if !ok {
		return nil, errors.New("name claim not found")
	}

	name, ok := nameInterface.(string)
	if !ok {
		return nil, errors.New("name claim is not a string")
	}

	clientInfo := ClientInfo{
		Name: name,
	}
	return &clientInfo, nil
}

func AuthHandler(c *gin.Context) {
	var params AuthParameters
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, err := googleAuthenticate(params.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("body: %s\n", string(body))
	tokenResponse, err := getTokenResponseFromResponseBody(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	claims, err := getGoogleAuthClaimsFromTokenResponse(tokenResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	clientInfo, err := getClientInfoFromClaims(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Response": "Ok", "ClientInfo": clientInfo})
}
