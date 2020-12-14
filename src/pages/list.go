package pages

import (
	"fmt"
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

	var itemList []string

	for _, item := range items {
		itemList = append(itemList, fmt.Sprintf("Name: %s, price: %.2f", item.Name, float64(item.Price)/100))
	}

	c.JSON(200, itemList)
}
