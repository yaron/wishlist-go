package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dgrijalva/jwt-go"
)

const hmacFile = "./hmac"

// GetHmac gets the Hmac from file or generates a new one and writes it to file
func GetHmac() ([]byte, error) {
	var r []byte
	if _, err := os.Stat(hmacFile); os.IsNotExist(err) {
		randomBytes := make([]byte, 16)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return r, err
		}
		secret := base64.URLEncoding.EncodeToString(randomBytes)
		randomBytes = make([]byte, 64)
		_, err = rand.Read(randomBytes)
		if err != nil {
			return r, err
		}
		data := base64.URLEncoding.EncodeToString(randomBytes)
		h := hmac.New(sha256.New, []byte(secret))
		h.Write([]byte(data))
		sha := []byte(hex.EncodeToString(h.Sum(nil)))
		err = ioutil.WriteFile(hmacFile, sha, 0600)
		if err != nil {
			return r, err
		}
		return sha, nil
	}
	data, err := ioutil.ReadFile(hmacFile)
	if err != nil {
		return r, err
	}
	return data, nil
}

// TestToken can be used to test if a token is valid and get the username from it
func TestToken(t string) (string, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		hmac, err := GetHmac()
		if err != nil {
			return nil, err
		}
		return hmac, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["username"].(string), nil
	}

	return "", fmt.Errorf("Unable to parse JWT")
}