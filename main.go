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
	r.Use(CORSMiddleware())
	r.GET("/list", pages.List)
	r.POST("/claim", pages.Claim)
	r.POST("/unclaim", pages.Unclaim)
	r.POST("/login", pages.Login)
	authorized := r.Group("/admin", jWTAuth)
	authorized.POST("/add", pages.Add)
	authorized.POST("/delete/:id", pages.Delete)
	authorized.POST("/edit/:id", pages.Edit)
	r.Run()
}


func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {

        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func jWTAuth(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	token := strings.Replace(header, "Bearer ", "", 1)
	userID, err := utils.TestToken(token)
	if err != nil {
		log.Println("Warning: " + err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set(gin.AuthUserKey, userID)

}
