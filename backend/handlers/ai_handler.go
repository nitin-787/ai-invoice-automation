package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nitin-787/ai-invoice-automation/models"
	"github.com/nitin-787/ai-invoice-automation/services"
)

type AIHandler struct {
	aiService      *services.AIService
	invoiceService *services.InvoiceService
}

func NewAIHandler(
	aiService *services.AIService,
	invoiceService *services.InvoiceService,
) *AIHandler {
	return &AIHandler{
		aiService:      aiService,
		invoiceService: invoiceService,
	}
}

type ExtractInvoiceRequest struct {
	Input string `json:"input" binding:"required"`
}

func (h *AIHandler) ExtractInvoice(c *gin.Context) {
	var request ExtractInvoiceRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	extracted, err := h.aiService.ExtractInvoice(
		c.Request.Context(),
		request.Input,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	invoice := &models.Invoice{
		InvoiceNumber: extracted.InvoiceNumber,
		VendorName:    extracted.VendorName,
		VendorEmail:   extracted.VendorEmail,
		Amount:        extracted.Amount,
		Currency:      extracted.Currency,
		Description:   extracted.Description,
		Source:        "ai-extraction",
		ExtractedData: extracted,
	}

	if extracted.InvoiceDate != nil {
		if parsed, err := time.Parse("2006-01-02", *extracted.InvoiceDate); err == nil {
			invoice.InvoiceDate = &parsed
		}
	}

	if extracted.DueDate != nil {
		if parsed, err := time.Parse("2006-01-02", *extracted.DueDate); err == nil {
			invoice.DueDate = &parsed
		}
	}

	if err := h.invoiceService.ProcessInvoice(
		c.Request.Context(),
		invoice,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

func (h *AIHandler) ExtractInvoiceFile(c *gin.Context) {
	file, err := c.FormFile("invoice")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invoice file is required",
		})
		return
	}

	allowedTypes := map[string]bool{
		"application/pdf": true,
		"image/jpeg":      true,
		"image/png":       true,
		"image/webp":      true,
	}

	contentType := file.Header.Get("Content-Type")

	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unsupported file type; use PDF, JPEG, PNG, or WEBP",
		})
		return
	}

	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file size must be less than 10MB",
		})
		return
	}

	uploadedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to open invoice file",
		})
		return
	}
	defer uploadedFile.Close()

	fileData, err := io.ReadAll(uploadedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to read invoice file",
		})
		return
	}

	if contentType != "application/pdf" {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "image extraction will be connected next",
		})
		return
	}

	extracted, err := h.aiService.ExtractInvoicePDF(
		c.Request.Context(),
		fileData,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	invoice := &models.Invoice{
		InvoiceNumber: extracted.InvoiceNumber,
		VendorName:    extracted.VendorName,
		VendorEmail:   extracted.VendorEmail,
		Amount:        extracted.Amount,
		Currency:      extracted.Currency,
		Description:   extracted.Description,
		Source:        "ai-pdf-extraction",
		ExtractedData: extracted,
	}

	if extracted.InvoiceDate != nil {
		if parsed, err := time.Parse("2006-01-02", *extracted.InvoiceDate); err == nil {
			invoice.InvoiceDate = &parsed
		}
	}

	if extracted.DueDate != nil {
		if parsed, err := time.Parse("2006-01-02", *extracted.DueDate); err == nil {
			invoice.DueDate = &parsed
		}
	}

	if err := h.invoiceService.ProcessInvoice(
		c.Request.Context(),
		invoice,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}
