package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nitin-787/ai-invoice-automation/db"
	"github.com/nitin-787/ai-invoice-automation/handlers"
	"github.com/nitin-787/ai-invoice-automation/repository"
	"github.com/nitin-787/ai-invoice-automation/services"
)

func main() {
	pool, err := db.NewPostgresPool()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	invoiceRepository := repository.NewInvoiceRepository(pool)
	invoiceService := services.NewInvoiceService(invoiceRepository)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"service":  "invoice-automation-api",
			"database": "connected",
		})
	})

	api := router.Group("/api/v1")
	{
		api.POST("/invoices", invoiceHandler.CreateInvoice)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
