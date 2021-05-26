package pages

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/yaron/wishlist-go/src/utils"
)

// Login checks provded credentials with the db and returns a JWT.
func Login(c *gin.Context) {
	var json utils.Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := utils.FetchUser(json.User, json.Password)

	if err != nil {
		log.Println("Unable to find user: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), jwt.MapClaims{
		"userID": user.ID,
		"nbf":    time.Now().Unix(),
		"exp":    time.Now().Unix() + 60*60*24, // Expire after 24 hours.
	})

	rsa, err := utils.GetRsa()
	if err != nil {
		log.Println("Unable to get rsa key: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(rsa)
	if err != nil {
		log.Println("Unable to get private key from keyfile: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	tokenString, err := token.SignedString(rsaKey)
	if err != nil {
		log.Println("Unable to sign token: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Welcome " + user.Username, "token": tokenString})
}

func JwtKey(c *gin.Context) {
	rsa, err := utils.GetRsa()
	if err != nil {
		log.Println("Unable to get rsa key: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(rsa)
	if err != nil {
		log.Println("Unable get private key: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	publickeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	pemPubKey := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: publickeyBytes}))

	c.JSON(http.StatusOK, gin.H{"key": pemPubKey})
}
