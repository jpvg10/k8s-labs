package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getVisits(db *sql.DB) (int, error) {
	var visits int
	err := db.QueryRow(`SELECT "visits" FROM "pingpong" WHERE "id"=$1`, 1).Scan(&visits)
	return visits, err
}

func incrementVisits(db *sql.DB) (int, error) {
	var visits int
	err := db.QueryRow(`
        UPDATE "pingpong"
        SET "visits" = "visits" + 1
        WHERE "id" = $1
        RETURNING "visits"
    `, 1).Scan(&visits)
	return visits, err
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

	router.GET("/ping", func(c *gin.Context) {
		visits, err := incrementVisits(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update ping counter",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"pong": visits,
		})
	})

	router.GET("/ping-count", func(c *gin.Context) {
		visits, err := getVisits(db)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read pings from the database",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"pings": visits,
			})
		}
	})

	port := getEnv("PORT", "4000")

	fmt.Printf("Server started in port %s\n", port)
	router.Run(":" + port)
}
