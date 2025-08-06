package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Description string `json:"description" binding:"required"`
}

func main() {
	router := gin.Default()

	todos := []Todo{}

	router.GET("/todos", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"todos": todos,
		})
	})

	router.POST("/todos", func(c *gin.Context) {
		var newTodo Todo

		if err := c.ShouldBindJSON(&newTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todos = append(todos, newTodo)

		c.JSON(http.StatusCreated, newTodo)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "6000"
	}

	fmt.Printf("Server started in port %s\n", port)
	router.Run(":" + port)
}
