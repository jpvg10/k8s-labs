package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	Description string `json:"description" binding:"required"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getTodos(db *sql.DB) ([]Todo, error) {
	todos := []Todo{}
	rows, err := db.Query(`SELECT "description" FROM "todos"`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var description string
		rows.Scan(&description)
		todos = append(todos, Todo{Description: description})
	}
	return todos, err
}

func addTodo(db *sql.DB, todo Todo) error {
	_, err := db.Exec(`insert into "todos"("description") values($1)`, todo.Description)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := 5432
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "example")
	dbName := getEnv("DB_NAME", "postgres")

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println("Database connection failed")
		fmt.Println(err)
		return
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/todos", func(c *gin.Context) {
		todos, err := getTodos(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read todos from the database",
			})
			return
		}
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

		if len(newTodo.Description) > 140 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Description must be maximum 140 characters",
			})
			fmt.Println("Failed to create todo due to exceeding 140 character limit.")
			return
		}

		err := addTodo(db, newTodo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to add new todo to the database",
			})
			return
		}

		c.JSON(http.StatusCreated, newTodo)
	})

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Database not available",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	port := getEnv("PORT", "6000")

	fmt.Printf("Server started in port %s\n", port)
	router.Run(":" + port)
}
