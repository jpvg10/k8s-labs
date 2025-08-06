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
	Description string
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

		c.HTML(http.StatusOK, "index.html", gin.H{
			"imagePath": "/image.jpg",
			"todos":     todos.Todos,
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	fmt.Printf("Server started in port %s\n", port)
	router.Run(":" + port)
}
