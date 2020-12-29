package pages

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/yaron/wishlist-go/src/utils"
)

// Claim marks an item as claimed
func Claim(c *gin.Context) {
	var json utils.Claim
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := utils.ClaimItem(json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domain := os.Getenv("WISH_MAILGUN_DOMAIN")
	key := os.Getenv("WISH_MAILGUN_KEY")
	url := os.Getenv("WISH_URL")
	siteName := os.Getenv("WISH_NAME")

	item, err := utils.FetchItem(json.ID)
	if err != nil {
		c.JSON(200, gin.H{"status": "Item claimed", "error": fmt.Sprintf("Unable to send mail %s", err.Error())})
		return
	}

	mg := mailgun.NewMailgun(domain, key)
	from := fmt.Sprintf("%s <mail@%s>", siteName, url)
	subject := fmt.Sprintf("Cadeau afgestreept %s", siteName)
	body := fmt.Sprintf(`Hallo,
	Je hebt op %s %s afgestreept. Hierbij de link om dit
	ongedaan te maken mocht het toch niet zijn gelukt:
	https://%s/unclaim`, siteName, item.Name, siteName)
	message := mg.NewMessage(from, subject, body, json.Mail)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	_, _, err = mg.Send(ctx, message)

	if err != nil {
		c.JSON(200, gin.H{"status": "Item claimed", "error": fmt.Sprintf("Unable to send mail %s", err.Error())})
		return
	}

	c.JSON(200, gin.H{"status": "Item claimed"})
}
