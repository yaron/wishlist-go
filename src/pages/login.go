package pages

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
		log.Println("Warning: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"nbf":    time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	hmac, err := utils.GetHmac()
	if err != nil {
		log.Println("Warning: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	tokenString, err := token.SignedString(hmac)
	if err != nil {
		log.Println("Warning: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Welcome " + user.Username, "token": tokenString})
}
