package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
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

// CreateWithApprovalLog creates the invoice and its initial approval log
// inside a single database transaction.
func (r *InvoiceRepository) CreateWithApprovalLog(
	ctx context.Context,
	invoice *models.Invoice,
	action string,
	actor *string,
	reason *string,
) error {
	extractedData, err := json.Marshal(invoice.ExtractedData)
	if err != nil {
		return fmt.Errorf("marshal extracted data: %w", err)
	}

	validationErrors, err := json.Marshal(invoice.ValidationErrors)
	if err != nil {
		return fmt.Errorf("marshal validation errors: %w", err)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin invoice transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	invoiceQuery := `
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

	err = tx.QueryRow(
		ctx,
		invoiceQuery,
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
		return fmt.Errorf("create invoice in transaction: %w", err)
	}

	logQuery := `
		INSERT INTO approval_logs (
			invoice_id,
			action,
			actor,
			reason
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(
		ctx,
		logQuery,
		invoice.ID,
		action,
		actor,
		reason,
	)

	if err != nil {
		return fmt.Errorf("create approval log in transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit invoice transaction: %w", err)
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

func (r *InvoiceRepository) GetAll(
	ctx context.Context,
) ([]models.Invoice, error) {
	query := `
		SELECT
			id,
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
			validation_errors,
			created_at,
			updated_at
		FROM invoices
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get invoices: %w", err)
	}
	defer rows.Close()

	var invoices []models.Invoice

	for rows.Next() {
		var invoice models.Invoice
		var extractedData []byte
		var validationErrors []byte

		err := rows.Scan(
			&invoice.ID,
			&invoice.InvoiceNumber,
			&invoice.VendorName,
			&invoice.VendorEmail,
			&invoice.InvoiceDate,
			&invoice.DueDate,
			&invoice.Amount,
			&invoice.Currency,
			&invoice.Description,
			&invoice.Status,
			&invoice.Source,
			&invoice.SourceReference,
			&extractedData,
			&validationErrors,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan invoice: %w", err)
		}

		if len(extractedData) > 0 {
			if err := json.Unmarshal(extractedData, &invoice.ExtractedData); err != nil {
				return nil, fmt.Errorf("unmarshal extracted data: %w", err)
			}
		}

		if len(validationErrors) > 0 {
			if err := json.Unmarshal(validationErrors, &invoice.ValidationErrors); err != nil {
				return nil, fmt.Errorf("unmarshal validation errors: %w", err)
			}
		}

		invoices = append(invoices, invoice)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate invoices: %w", err)
	}

	return invoices, nil
}

func (r *InvoiceRepository) GetByID(
	ctx context.Context,
	invoiceID string,
) (*models.Invoice, error) {
	query := `
		SELECT
			id,
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
			validation_errors,
			created_at,
			updated_at
		FROM invoices
		WHERE id = $1
	`

	var invoice models.Invoice
	var extractedData []byte
	var validationErrors []byte

	err := r.db.QueryRow(ctx, query, invoiceID).Scan(
		&invoice.ID,
		&invoice.InvoiceNumber,
		&invoice.VendorName,
		&invoice.VendorEmail,
		&invoice.InvoiceDate,
		&invoice.DueDate,
		&invoice.Amount,
		&invoice.Currency,
		&invoice.Description,
		&invoice.Status,
		&invoice.Source,
		&invoice.SourceReference,
		&extractedData,
		&validationErrors,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invoice not found")
		}

		return nil, fmt.Errorf("get invoice: %w", err)
	}

	if len(extractedData) > 0 {
		if err := json.Unmarshal(extractedData, &invoice.ExtractedData); err != nil {
			return nil, fmt.Errorf("unmarshal extracted data: %w", err)
		}
	}

	if len(validationErrors) > 0 {
		if err := json.Unmarshal(validationErrors, &invoice.ValidationErrors); err != nil {
			return nil, fmt.Errorf("unmarshal validation errors: %w", err)
		}
	}

	return &invoice, nil
}
