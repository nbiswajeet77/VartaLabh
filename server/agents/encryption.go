package agents

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Encrypts the struct into a byte slice
func EncryptStruct(s interface{}) (string, error) {
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(jsonBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], jsonBytes)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypts a byte slice into a struct
func DecryptBytes(encodedData string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, err
	}
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext is too short")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}
