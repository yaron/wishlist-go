package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/golang-jwt/jwt"
)

func rsaFile() string {
	return path.Join(os.Getenv("WISH_PATH"), "rsa.pem")
}

// GetRsa gets the rsa pem from file or generates a new one and writes it to file
func GetRsa() ([]byte, error) {
	var r []byte
	if _, err := os.Stat(rsaFile()); os.IsNotExist(err) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return r, err
		}
		var r []byte = x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyBlock := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: r,
		}

		privatePem, err := os.Create(rsaFile())
		if err != nil {
			fmt.Printf("error when creating rsa.pem: %s \n", err)
			os.Exit(1)
		}
		err = pem.Encode(privatePem, privateKeyBlock)
		if err != nil {
			fmt.Printf("error when encoding rsa.pem: %s \n", err)
			os.Exit(1)
		}
		return r, nil
	}
	r, err := ioutil.ReadFile(rsaFile())
	if err != nil {
		return r, err
	}
	return r, nil
}

// TestToken can be used to test if a token is valid and get the userID from it
func TestToken(t string) (int, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		pem, err := GetRsa()
		if err != nil {
			return nil, err
		}
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
		if err != nil {
			return nil, err
		}

		return &privateKey.PublicKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if _, ok := claims["exp"]; !ok {
			return 0, fmt.Errorf("Old tokens without expiration are no longer valid")
		}
		return int(claims["userID"].(float64)), nil
	}

	return 0, fmt.Errorf("Unable to parse JWT")
}
