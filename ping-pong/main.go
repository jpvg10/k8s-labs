package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	visits := 0

	router.GET("/ping", func(c *gin.Context) {
		visits++

		file, err := os.Create("/usr/src/app/files/visits.txt")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = file.WriteString(strconv.Itoa(visits))
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"pong": visits,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server started in port %s\n", port)
	router.Run(":" + port)
}
