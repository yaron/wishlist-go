package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yaron/wishlist-go/src/utils"
)

// Add creates new wishlist items in the db
func Add(c *gin.Context) {
	var json utils.WishlistItem
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := utils.AddItem(json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "Item added"})
}
