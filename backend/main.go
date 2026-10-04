package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nitin-787/ai-invoice-automation/db"
	"github.com/nitin-787/ai-invoice-automation/handlers"
	"github.com/nitin-787/ai-invoice-automation/middleware"
	"github.com/nitin-787/ai-invoice-automation/repository"
	"github.com/nitin-787/ai-invoice-automation/services"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Origin, Content-Type, Accept, Authorization, X-API-Key, X-Reviewer",
		)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env file not found")
	}

	pool, err := db.NewPostgresPool()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	invoiceRepository := repository.NewInvoiceRepository(pool)
	invoiceService := services.NewInvoiceService(invoiceRepository)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)

	aiService, err := services.NewAIService(context.Background())
	if err != nil {
		log.Fatalf("AI service initialization failed: %v", err)
	}

	aiHandler := handlers.NewAIHandler(aiService, invoiceService)

	router := gin.Default()

	// CORS
	router.Use(corsMiddleware())

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
		api.GET("/invoices", invoiceHandler.GetInvoices)
		api.GET("/invoices/:id", invoiceHandler.GetInvoice)

		api.POST("/invoices/extract", aiHandler.ExtractInvoice)
		api.POST("/invoices/extract/file", aiHandler.ExtractInvoiceFile)

		approval := api.Group("/invoices/:id")
		approval.Use(middleware.APIKeyAuth())
		{
			approval.POST("/approve", invoiceHandler.ApproveInvoice)
			approval.POST("/reject", invoiceHandler.RejectInvoice)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
    	port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
