package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nitin-787/ai-invoice-automation/db"
)

func main() {
	pool, err := db.NewPostgresPool()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "invoice-automation-api",
			"database": "connected",
		})
	})

	router.Run(":8080")
}