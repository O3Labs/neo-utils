package neoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Encrypt string to base64 format using AES
func Encrypt(key []byte, text string) string {
	plaintext := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	// Convert to base64
	return base64.URLEncoding.EncodeToString(cipherText)
}

// Decrypt AES encrypted string in base64 format to decrypted string
func Decrypt(key []byte, encryptedText string) string {
	cipherText, _ := base64.URLEncoding.DecodeString(encryptedText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	if len(cipherText) < aes.BlockSize {
		return ""
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	return fmt.Sprintf("%s", cipherText)
}
