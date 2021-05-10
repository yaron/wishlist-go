package pages

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"text/template"
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
	item, err := utils.FetchItem(json.ID)
	if !item.Claimable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item is not claimable"})
		return
	}

	err = utils.ClaimItem(json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(json.Mail) > 0 {
		domain := os.Getenv("WISH_MAILGUN_DOMAIN")
		key := os.Getenv("WISH_MAILGUN_KEY")
		eu := os.Getenv("WISH_MAILGUN_EU")
		from := os.Getenv("WISH_MAIL_FROM")
		subjectTemplate := os.Getenv("WISH_MAIL_SUBJECT")
		bodyTemplate := os.Getenv("WISH_MAIL_BODY")

		if err != nil {
			c.JSON(200, gin.H{"status": "Item claimed", "error": fmt.Sprintf("Unable to send mail %s", err.Error())})
			return
		}

		claimKey, err := utils.ClaimKey(json.ID)
		if err != nil {
			c.JSON(200, gin.H{"status": "Item claimed", "error": fmt.Sprintf("Unable to create claimkey %s", err.Error())})
			return
		}

		data := struct {
			ItemName string
			ItemID   int
			Key      string
			Mail     string
		}{
			ItemName: item.Name,
			ItemID:   item.ID,
			Key:      claimKey,
			Mail:     json.Mail,
		}

		mg := mailgun.NewMailgun(domain, key)
		if eu == "1" {
			mg.SetAPIBase(mailgun.APIBaseEU)
		}
		var body bytes.Buffer
		bodyTmpl := template.Must(template.New("body").Parse(bodyTemplate))
		err = bodyTmpl.Execute(&body, data)
		if err != nil {
			c.JSON(200, gin.H{"status": "Item claimed", "error": fmt.Sprintf("Unable to send mail %s", err.Error())})
			return
		}
		var subject bytes.Buffer
		subjectTmpl := template.Must(template.New("subject").Parse(subjectTemplate))
		err = subjectTmpl.Execute(&subject, data)
		if err != nil {
			c.JSON(200, gin.H{"status": "Item claimed", "error": fmt.Sprintf("Unable to send mail %s", err.Error())})
			return
		}
		message := mg.NewMessage(from, subject.String(), body.String(), json.Mail)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// Send the message with a 10 second timeout
		_, _, err = mg.Send(ctx, message)

		if err != nil {
			c.JSON(200, gin.H{"status": "Item claimed", "error": fmt.Sprintf("Unable to send mail %s", err.Error())})
			return
		}
	}
	c.JSON(200, gin.H{"status": "Item claimed"})
}

// Unclaim marks an item as unclaimed
func Unclaim(c *gin.Context) {
	var json utils.Unclaim
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := utils.UnclaimItem(json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "Item unclaimed"})
}
