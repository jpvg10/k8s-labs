package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	visits := 0

	router.GET("/ping", func(c *gin.Context) {
		visits++

		c.JSON(http.StatusOK, gin.H{
			"pong": visits,
		})
	})

	router.GET("/ping-count", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"pings": visits,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	fmt.Printf("Server started in port %s\n", port)
	router.Run(":" + port)
}
