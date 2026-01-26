package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Id          int    `json:"id"`
	Description string `json:"description" binding:"required"`
	Done        bool   `json:"done"`
}

type TodoArray struct {
	Todos []Todo
}

func getImage(imagePath string) (err error) {
	url := "https://picsum.photos/200"

	out, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer out.Close()

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	todoBackendUrl := os.Getenv("TODO_BACKEND_URL")
	imagePath := "/usr/src/app/files/image.jpg"
	var timestamp time.Time

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.StaticFile("/image.jpg", imagePath)

	router.GET("/", func(c *gin.Context) {
		current := time.Now()

		if current.Sub(timestamp).Seconds() > 10*60 {
			fmt.Println("Getting new image")
			err := getImage(imagePath)
			if err != nil {
				fmt.Println("Error fetching the image")
			}
			timestamp = current
		} else {
			fmt.Println("Using cached image")
		}

		var todos TodoArray
		res, err := http.Get(fmt.Sprintf("%s/todos", todoBackendUrl))
		if err != nil {
			fmt.Println("Error fetching the Todos")
		} else {
			defer res.Body.Close()
			decoder := json.NewDecoder(res.Body)
			err = decoder.Decode(&todos)
			if err != nil {
				fmt.Println("Error decoding the Todos")
			}
		}

		pending := []Todo{}
		done := []Todo{}
		for _, t := range todos.Todos {
			if t.Done {
				done = append(done, t)
			} else {
				pending = append(pending, t)
			}
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"imagePath": "/image.jpg",
			"pending":   pending,
			"done":      done,
		})
	})

	router.POST("/createTodo", func(c *gin.Context) {
		desc := c.PostForm("description")
		jsonValue, _ := json.Marshal(gin.H{"description": desc})
		res, err := http.Post(fmt.Sprintf("%s/todos", todoBackendUrl), "application/json", bytes.NewBuffer(jsonValue))
		if err != nil || res.StatusCode >= 400 {
			fmt.Println("Error posting the Todo")
		}
		c.Redirect(http.StatusFound, "/")
	})

	router.POST("/markTodoAsDone", func(c *gin.Context) {
		id := c.PostForm("id")
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/todos/%s", todoBackendUrl, id), nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode >= 400 {
			fmt.Println("Error marking the Todo as done")
		}
		c.Redirect(http.StatusFound, "/")
	})

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		_, err := http.Get(fmt.Sprintf("%s/healthz", todoBackendUrl))
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Backend service not available",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "ready",
			})
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	fmt.Printf("Server started in port %s\n", port)
	router.Run(":" + port)
}
