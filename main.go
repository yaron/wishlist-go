package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yaron/wishlist-go/src/pages"
	"github.com/yaron/wishlist-go/src/utils"
)

func main() {
	r := gin.Default()
	r.GET("/list", pages.List)
	r.POST("/claim", pages.Claim)
	r.POST("/unclaim", pages.Unclaim)
	r.POST("/login", pages.Login)
	authorized := r.Group("/admin", jWTAuth)
	authorized.POST("/add", pages.Add)
	authorized.POST("/edit/:id", pages.Edit)
	r.Run()
}

func jWTAuth(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	token := strings.Replace(header, "Bearer ", "", 1)
	username, err := utils.TestToken(token)
	if err != nil {
		log.Println("Warning: " + err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set(gin.AuthUserKey, username)

}
