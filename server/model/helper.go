package model

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	jsoniter "github.com/json-iterator/go"
)

var jsonIter = jsoniter.ConfigDefault
var secret string

func GetDecryptionSecret() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secret = os.Getenv("DECRYPTION_SECRET")
}

func WriteOutput(w http.ResponseWriter, data interface{}, status int, err error) {
	response := Response{
		StatusCode: status,
		Data:       data,
	}
	w.WriteHeader(status)
	if err != nil {
		byteData, _ := jsonIter.Marshal(err)
		w.Write(byteData)
		return
	}
	byteData, marshalErr := jsonIter.Marshal(response)
	if marshalErr != nil {
		fmt.Println(marshalErr)
	}
	w.Write(byteData)
}

func GenerateToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"username":  userId,
		"createdAt": time.Now().Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Error generating token:", err)
		return "", err
	}
	return tokenString, nil
}

func DecryptToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure the token's signing method is HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used to sign the token
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["username"].(string), nil
	}
	return "", errors.New("error while parsing jwt token")
}
