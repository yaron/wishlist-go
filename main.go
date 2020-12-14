package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yaron/wishlist-go/src/pages"
)

func main() {
	r := gin.Default()
	r.GET("/list", pages.List)
	r.Run()
}
