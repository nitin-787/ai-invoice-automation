package services

import (
	"context"
	"testing"

	"github.com/nitin-787/ai-invoice-automation/db"
	"github.com/nitin-787/ai-invoice-automation/models"
	"github.com/nitin-787/ai-invoice-automation/repository"
)

func setupTestService(t *testing.T) (*InvoiceService, context.Context) {
	t.Helper()

	ctx := context.Background()

	pool, err := db.NewPostgresPool()
	if err != nil {
		t.Fatalf("database connection failed: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	repo := repository.NewInvoiceRepository(pool)
	service := NewInvoiceService(repo)

	return service, ctx
}

func TestProcessInvoiceAutoApproval(t *testing.T) {
	service, ctx := setupTestService(t)

	invoice := &models.Invoice{
		InvoiceNumber: "TEST-AUTO-001",
		VendorName:    "Test Vendor Auto",
		Amount:        42000,
		Currency:      "INR",
		Source:        "test",
	}

	err := service.ProcessInvoice(ctx, invoice)
	if err != nil {
		t.Fatalf("ProcessInvoice() error = %v", err)
	}

	if invoice.Status != "APPROVED" {
		t.Fatalf(
			"expected status APPROVED, got %s",
			invoice.Status,
		)
	}
}

func TestProcessInvoicePendingApproval(t *testing.T) {
	service, ctx := setupTestService(t)

	invoice := &models.Invoice{
		InvoiceNumber: "TEST-PENDING-001",
		VendorName:    "Test Vendor Pending",
		Amount:        85000,
		Currency:      "INR",
		Source:        "test",
	}

	err := service.ProcessInvoice(ctx, invoice)
	if err != nil {
		t.Fatalf("ProcessInvoice() error = %v", err)
	}

	if invoice.Status != "PENDING_APPROVAL" {
		t.Fatalf(
			"expected status PENDING_APPROVAL, got %s",
			invoice.Status,
		)
	}
}

func TestProcessInvoiceRejectsMissingVendor(t *testing.T) {
	service, ctx := setupTestService(t)

	invoice := &models.Invoice{
		InvoiceNumber: "TEST-NO-VENDOR-001",
		Amount:        42000,
		Currency:      "INR",
	}

	err := service.ProcessInvoice(ctx, invoice)

	if err == nil {
		t.Fatal("expected error for missing vendor")
	}
}

func TestProcessInvoiceRejectsMissingInvoiceNumber(t *testing.T) {
	service, ctx := setupTestService(t)

	invoice := &models.Invoice{
		VendorName: "Test Vendor",
		Amount:     42000,
		Currency:   "INR",
	}

	err := service.ProcessInvoice(ctx, invoice)

	if err == nil {
		t.Fatal("expected error for missing invoice number")
	}
}

func TestProcessInvoiceRejectsInvalidAmount(t *testing.T) {
	service, ctx := setupTestService(t)

	invoice := &models.Invoice{
		InvoiceNumber: "TEST-INVALID-AMOUNT-001",
		VendorName:    "Test Vendor",
		Amount:        0,
		Currency:      "INR",
	}

	err := service.ProcessInvoice(ctx, invoice)

	if err == nil {
		t.Fatal("expected error for invalid amount")
	}
}
