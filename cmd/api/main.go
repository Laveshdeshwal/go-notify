package main

import "github.com/gin-gonic/gin"

func main() {
	app := gin.Default()

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Server": "go-notify-server",
			"Status": "running",
		})
	})

	app.Run(":8080")
}
