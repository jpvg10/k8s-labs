package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

type Todo struct {
	Id          int    `json:"id"`
	Description string `json:"description" binding:"required"`
	Done        bool   `json:"done"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getTodos(db *sql.DB) ([]Todo, error) {
	todos := []Todo{}
	rows, err := db.Query(`SELECT "id", "description", "done" FROM "todos"`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var description string
		var done bool
		rows.Scan(&id, &description, &done)
		todos = append(todos, Todo{Id: id, Description: description, Done: done})
	}
	return todos, err
}

func addTodo(db *sql.DB, todo Todo, nc *nats.Conn) (Todo, error) {
	var insertedTodo Todo
	err := db.QueryRow(
		`insert into "todos"("description", "done") values($1, $2) returning "id", "description", "done"`,
		todo.Description,
		false,
	).Scan(&insertedTodo.Id, &insertedTodo.Description, &insertedTodo.Done)
	if err != nil {
		fmt.Println(err)
		return insertedTodo, err
	}

	message, err := json.Marshal(insertedTodo)
	sendNatsMessage(nc, "todos.created", message)

	return insertedTodo, nil
}

func markTodoAsDone(db *sql.DB, id int, nc *nats.Conn) (Todo, error) {
	var updatedTodo Todo
	err := db.QueryRow(
		`update "todos" set "done"=$2 where "id"=$1 returning "id", "description", "done"`,
		id,
		true,
	).Scan(&updatedTodo.Id, &updatedTodo.Description, &updatedTodo.Done)
	if err != nil {
		fmt.Println(err)
		return updatedTodo, err
	}

	message, err := json.Marshal(updatedTodo)
	sendNatsMessage(nc, "todos.updated", message)

	return updatedTodo, nil
}

func sendNatsMessage(nc *nats.Conn, topic string, message []byte) {
	if err := nc.Publish(topic, []byte(message)); err != nil {
		fmt.Printf("Failed to publish NATS message: %v\n", err)
	}
}

// TODO: if it starts without NATS it won't connect when NATS becomes available

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := 5432
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "example")
	dbName := getEnv("DB_NAME", "postgres")
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println("Database connection failed")
		fmt.Println(err)
		return
	}
	defer db.Close()

	var nc *nats.Conn
	nc, err = nats.Connect(natsURL)
	if err != nil {
		fmt.Printf("NATS connection failed: %v\n", err)
		fmt.Println("Continuing without NATS...")
	} else {
		defer nc.Close()
		fmt.Println("Connected to NATS successfully")
	}

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

		insertedTodo, err := addTodo(db, newTodo, nc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to add new todo to the database",
			})
			return
		}

		c.JSON(http.StatusCreated, insertedTodo)
	})

	router.PUT("/todos/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		updatedTodo, err := markTodoAsDone(db, id, nc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to mark the todo as done",
			})
			return
		}
		c.JSON(http.StatusOK, updatedTodo)
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
