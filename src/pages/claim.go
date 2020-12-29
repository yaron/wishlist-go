package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	c.JSON(200, gin.H{"status": "Item claimed"})
}
