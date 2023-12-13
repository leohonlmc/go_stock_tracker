package main

import "github.com/gin-gonic/gin"

func setupRoutes(router *gin.Engine) {
	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
}
