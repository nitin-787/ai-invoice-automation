package services

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nitin-787/ai-invoice-automation/db"
	"github.com/nitin-787/ai-invoice-automation/models"
	"github.com/nitin-787/ai-invoice-automation/repository"
)

func setupTestService(t *testing.T) (*InvoiceService, context.Context, *pgxpool.Pool) {
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

	return service, ctx, pool
}

func cleanupTestInvoice(t *testing.T, ctx context.Context, pool *pgxpool.Pool, invoiceID string) {
	t.Helper()

	_, err := pool.Exec(
		ctx,
		"DELETE FROM invoices WHERE id = $1",
		invoiceID,
	)
	if err != nil {
		t.Fatalf("failed to cleanup test invoice: %v", err)
	}
}

func TestProcessInvoiceAutoApproval(t *testing.T) {
	service, ctx, pool := setupTestService(t)

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

	t.Cleanup(func() {
		cleanupTestInvoice(t, ctx, pool, invoice.ID)
	})

	if invoice.Status != "APPROVED" {
		t.Fatalf("expected status APPROVED, got %s", invoice.Status)
	}
}

func TestProcessInvoicePendingApproval(t *testing.T) {
	service, ctx, pool := setupTestService(t)

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

	t.Cleanup(func() {
		cleanupTestInvoice(t, ctx, pool, invoice.ID)
	})

	if invoice.Status != "PENDING_APPROVAL" {
		t.Fatalf(
			"expected status PENDING_APPROVAL, got %s",
			invoice.Status,
		)
	}
}

func TestProcessInvoiceRejectsMissingVendor(t *testing.T) {
	service, _, _ := setupTestService(t)

	invoice := &models.Invoice{
		InvoiceNumber: "TEST-NO-VENDOR-002",
		Amount:        42000,
		Currency:      "INR",
	}

	err := service.ProcessInvoice(context.Background(), invoice)

	if err == nil {
		t.Fatal("expected error for missing vendor")
	}
}

func TestProcessInvoiceRejectsMissingInvoiceNumber(t *testing.T) {
	service, _, _ := setupTestService(t)

	invoice := &models.Invoice{
		VendorName: "Test Vendor",
		Amount:     42000,
		Currency:   "INR",
	}

	err := service.ProcessInvoice(context.Background(), invoice)

	if err == nil {
		t.Fatal("expected error for missing invoice number")
	}
}

func TestProcessInvoiceRejectsInvalidAmount(t *testing.T) {
	service, _, _ := setupTestService(t)

	invoice := &models.Invoice{
		InvoiceNumber: "TEST-INVALID-AMOUNT-002",
		VendorName:    "Test Vendor",
		Amount:        0,
		Currency:      "INR",
	}

	err := service.ProcessInvoice(context.Background(), invoice)

	if err == nil {
		t.Fatal("expected error for invalid amount")
	}
}

func TestApproveInvoice(t *testing.T) {
	service, ctx, pool := setupTestService(t)

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

	t.Cleanup(func() {
		cleanupTestInvoice(t, ctx, pool, invoice.ID)
	})

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
	service, ctx, pool := setupTestService(t)

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

	t.Cleanup(func() {
		cleanupTestInvoice(t, ctx, pool, invoice.ID)
	})

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

func TestProcessInvoiceRejectsDuplicate(t *testing.T) {
	service, ctx, pool := setupTestService(t)

	invoice := &models.Invoice{
		InvoiceNumber: "TEST-DUPLICATE-001",
		VendorName:    "Test Duplicate Vendor",
		Amount:        42000,
		Currency:      "INR",
		Source:        "test",
	}

	if err := service.ProcessInvoice(ctx, invoice); err != nil {
		t.Fatalf("first ProcessInvoice() error = %v", err)
	}

	t.Cleanup(func() {
		cleanupTestInvoice(t, ctx, pool, invoice.ID)
	})

	duplicate := &models.Invoice{
		InvoiceNumber: "TEST-DUPLICATE-001",
		VendorName:    "Test Duplicate Vendor",
		Amount:        42000,
		Currency:      "INR",
		Source:        "test",
	}

	err := service.ProcessInvoice(ctx, duplicate)

	if err == nil {
		t.Fatal("expected duplicate invoice error")
	}
}
