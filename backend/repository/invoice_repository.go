package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nitin-787/ai-invoice-automation/models"
)

type InvoiceRepository struct {
	db *pgxpool.Pool
}

func NewInvoiceRepository(db *pgxpool.Pool) *InvoiceRepository {
	return &InvoiceRepository{
		db: db,
	}
}

func (r *InvoiceRepository) Create(
	ctx context.Context,
	invoice *models.Invoice,
) error {
	extractedData, err := json.Marshal(invoice.ExtractedData)
	if err != nil {
		return fmt.Errorf("marshal extracted data: %w", err)
	}

	validationErrors, err := json.Marshal(invoice.ValidationErrors)
	if err != nil {
		return fmt.Errorf("marshal validation errors: %w", err)
	}

	query := `
		INSERT INTO invoices (
			invoice_number,
			vendor_name,
			vendor_email,
			invoice_date,
			due_date,
			amount,
			currency,
			description,
			status,
			source,
			source_reference,
			extracted_data,
			validation_errors
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13
		)
		RETURNING id, created_at, updated_at
	`

	err = r.db.QueryRow(
		ctx,
		query,
		invoice.InvoiceNumber,
		invoice.VendorName,
		invoice.VendorEmail,
		invoice.InvoiceDate,
		invoice.DueDate,
		invoice.Amount,
		invoice.Currency,
		invoice.Description,
		invoice.Status,
		invoice.Source,
		invoice.SourceReference,
		extractedData,
		validationErrors,
	).Scan(
		&invoice.ID,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create invoice: %w", err)
	}

	return nil
}

func (r *InvoiceRepository) FindDuplicate(
	ctx context.Context,
	invoiceNumber string,
	vendorName string,
) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM invoices
			WHERE invoice_number = $1
			  AND vendor_name = $2
		)
	`

	var exists bool

	err := r.db.QueryRow(
		ctx,
		query,
		invoiceNumber,
		vendorName,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check duplicate invoice: %w", err)
	}

	return exists, nil
}

func (r *InvoiceRepository) CreateApprovalLog(
	ctx context.Context,
	invoiceID string,
	action string,
	actor *string,
	reason *string,
) error {
	query := `
		INSERT INTO approval_logs (
			invoice_id,
			action,
			actor,
			reason
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		invoiceID,
		action,
		actor,
		reason,
	)

	if err != nil {
		return fmt.Errorf("create approval log: %w", err)
	}

	return nil
}

func (r *InvoiceRepository) UpdateStatus(
	ctx context.Context,
	invoiceID string,
	status string,
) error {
	query := `
		UPDATE invoices
		SET status = $1,
		    updated_at = NOW()
		WHERE id = $2
		  AND status = 'PENDING_APPROVAL'
	`

	result, err := r.db.Exec(
		ctx,
		query,
		status,
		invoiceID,
	)
	if err != nil {
		return fmt.Errorf("update invoice status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("invoice not found or not pending approval")
	}

	return nil
}
