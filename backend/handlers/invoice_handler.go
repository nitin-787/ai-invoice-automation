package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nitin-787/ai-invoice-automation/models"
	"github.com/nitin-787/ai-invoice-automation/services"
)

type InvoiceHandler struct {
	service *services.InvoiceService
}

func NewInvoiceHandler(
	service *services.InvoiceService,
) *InvoiceHandler {
	return &InvoiceHandler{
		service: service,
	}
}

func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var invoice models.Invoice

	if err := c.ShouldBindJSON(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.ProcessInvoice(
		c.Request.Context(),
		&invoice,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}
