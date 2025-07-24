package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func getImage(imagePath string) (err error) {
	url := "https://picsum.photos/200"

	out, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var timestamp time.Time
	imagePath := "/usr/src/app/files/image.jpg"

	router := gin.Default()
	router.LoadHTMLFiles("index.html")
	router.StaticFile("/image.jpg", imagePath)

	router.GET("/", func(c *gin.Context) {
		current := time.Now()

		if current.Sub(timestamp).Seconds() > 10*60 {
			fmt.Println("Getting new image")
			getImage(imagePath)
			timestamp = current
		} else {
			fmt.Println("Using cached image")
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"imagePath": "/image.jpg",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server started in port %s\n", port)
	router.Run(":" + port)
}
