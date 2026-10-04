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

func cleanupTestInvoice(t *testing.T, ctx context.Context, invoiceID string) {
	t.Helper()

	pool, err := db.NewPostgresPool()
	if err != nil {
		t.Fatalf("database connection failed during cleanup: %v", err)
	}
	defer pool.Close()

	_, err = pool.Exec(
		ctx,
		"DELETE FROM invoices WHERE id = $1",
		invoiceID,
	)
	if err != nil {
		t.Fatalf("failed to cleanup test invoice: %v", err)
	}
}

func TestProcessInvoiceAutoApproval(t *testing.T) {
	service, ctx := setupTestService(t)

	invoice := &models.Invoice{
		InvoiceNumber: "TEST-AUTO-002",
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
		InvoiceNumber: "TEST-PENDING-002",
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
		InvoiceNumber: "TEST-NO-VENDOR-002",
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
		InvoiceNumber: "TEST-INVALID-AMOUNT-002",
		VendorName:    "Test Vendor",
		Amount:        0,
		Currency:      "INR",
	}

	err := service.ProcessInvoice(ctx, invoice)

	if err == nil {
		t.Fatal("expected error for invalid amount")
	}
}

func TestProcessInvoiceTrimsWhitespace(t *testing.T) {
	service, ctx := setupTestService(t)

	invoice := &models.Invoice{
		InvoiceNumber: "  TEST-TRIM-001  ",
		VendorName:    "  Test Vendor Trim  ",
		Amount:        42000,
		Currency:      " INR ",
		Source:        " test ",
	}

	err := service.ProcessInvoice(ctx, invoice)
	if err != nil {
		t.Fatalf("ProcessInvoice() error = %v", err)
	}

	if invoice.InvoiceNumber != "TEST-TRIM-001" {
		t.Fatalf(
			"expected trimmed invoice number, got %q",
			invoice.InvoiceNumber,
		)
	}

	if invoice.VendorName != "Test Vendor Trim" {
		t.Fatalf(
			"expected trimmed vendor name, got %q",
			invoice.VendorName,
		)
	}

	if invoice.Currency != "INR" {
		t.Fatalf(
			"expected trimmed currency, got %q",
			invoice.Currency,
		)
	}

	if invoice.Source != "test" {
		t.Fatalf(
			"expected trimmed source, got %q",
			invoice.Source,
		)
	}
}

func TestApproveInvoice(t *testing.T) {
	service, ctx := setupTestService(t)

	invoice := &models.Invoice{
		InvoiceNumber: "TEST-APPROVE-002",
		VendorName:    "Test Vendor Approve",
		Amount:        85000,
		Currency:      "INR",
		Source:        "test",
	}

	if err := service.ProcessInvoice(ctx, invoice); err != nil {
		t.Fatalf("ProcessInvoice() error = %v", err)
	}

	if invoice.Status != "PENDING_APPROVAL" {
		t.Fatalf("expected PENDING_APPROVAL, got %s", invoice.Status)
	}

	err := service.ApproveInvoice(
		ctx,
		invoice.ID,
		"test-reviewer",
	)
	if err != nil {
		t.Fatalf("ApproveInvoice() error = %v", err)
	}
}

func TestRejectInvoice(t *testing.T) {
	service, ctx := setupTestService(t)

	invoice := &models.Invoice{
		InvoiceNumber: "TEST-REJECT-002",
		VendorName:    "Test Vendor Reject",
		Amount:        85000,
		Currency:      "INR",
		Source:        "test",
	}

	if err := service.ProcessInvoice(ctx, invoice); err != nil {
		t.Fatalf("ProcessInvoice() error = %v", err)
	}

	if invoice.Status != "PENDING_APPROVAL" {
		t.Fatalf("expected PENDING_APPROVAL, got %s", invoice.Status)
	}

	err := service.RejectInvoice(
		ctx,
		invoice.ID,
		"test-reviewer",
	)
	if err != nil {
		t.Fatalf("RejectInvoice() error = %v", err)
	}
}
