package handlers

import (
	"errors"
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
			"error":   "validation failed",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.ProcessInvoice(
		c.Request.Context(),
		&invoice,
	); err != nil {
		if errors.Is(err, services.ErrDuplicateInvoice) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

func (h *InvoiceHandler) ApproveInvoice(c *gin.Context) {
	invoiceID := c.Param("id")

	reviewer, exists := c.Get("reviewer")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "reviewer identity not found",
		})
		return
	}

	reviewerName, ok := reviewer.(string)
	if !ok || reviewerName == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid reviewer identity",
		})
		return
	}

	if err := h.service.ApproveInvoice(
		c.Request.Context(),
		invoiceID,
		reviewerName,
	); err != nil {
		if errors.Is(err, services.ErrInvoiceNotPending) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, services.ErrInvoiceIDRequired) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  "APPROVED",
		"message": "Invoice approved successfully",
	})
}

func (h *InvoiceHandler) RejectInvoice(c *gin.Context) {
	invoiceID := c.Param("id")

	reviewer, exists := c.Get("reviewer")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "reviewer identity not found",
		})
		return
	}

	reviewerName, ok := reviewer.(string)
	if !ok || reviewerName == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid reviewer identity",
		})
		return
	}

	if err := h.service.RejectInvoice(
		c.Request.Context(),
		invoiceID,
		reviewerName,
	); err != nil {
		if errors.Is(err, services.ErrInvoiceNotPending) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, services.ErrInvoiceIDRequired) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  "REJECTED",
		"message": "Invoice rejected successfully",
	})
}

func (h *InvoiceHandler) GetInvoice(c *gin.Context) {
	invoiceID := c.Param("id")

	if invoiceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invoice ID is required",
		})
		return
	}

	invoice, err := h.service.GetInvoiceByID(
		c.Request.Context(),
		invoiceID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

func (h *InvoiceHandler) GetInvoices(c *gin.Context) {
	invoices, err := h.service.GetAllInvoices(
		c.Request.Context(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"invoices": invoices,
		"count":    len(invoices),
	})
}
