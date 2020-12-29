package pages

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yaron/wishlist-go/src/utils"
)

// List returns a list of all items in the wishlist as JSON
func List(c *gin.Context) {
	items, err := utils.FetchItems()
	if err != nil {
		log.Println("Warning: " + err.Error())
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, items)
}
