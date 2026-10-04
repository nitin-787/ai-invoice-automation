package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nitin-787/ai-invoice-automation/models"
	"github.com/nitin-787/ai-invoice-automation/repository"
)

var (
	ErrDuplicateInvoice  = errors.New("duplicate invoice")
	ErrInvoiceNotPending = errors.New("invoice is not pending approval")
	ErrInvoiceIDRequired = errors.New("invoice ID is required")
)

const autoApprovalLimit = 50000.00

type InvoiceService struct {
	repository *repository.InvoiceRepository
}

func NewInvoiceService(
	repository *repository.InvoiceRepository,
) *InvoiceService {
	return &InvoiceService{
		repository: repository,
	}
}

func (s *InvoiceService) ProcessInvoice(
	ctx context.Context,
	invoice *models.Invoice,
) error {
	invoice.InvoiceNumber = strings.TrimSpace(invoice.InvoiceNumber)
	invoice.VendorName = strings.TrimSpace(invoice.VendorName)
	invoice.Currency = strings.TrimSpace(invoice.Currency)
	invoice.Source = strings.TrimSpace(invoice.Source)

	if invoice.VendorEmail != nil {
		email := strings.TrimSpace(*invoice.VendorEmail)
		invoice.VendorEmail = &email
	}

	if invoice.Description != nil {
		description := strings.TrimSpace(*invoice.Description)
		invoice.Description = &description
	}

	if invoice.VendorName == "" {
		return fmt.Errorf("vendor name is required")
	}

	if invoice.InvoiceNumber == "" {
		return fmt.Errorf("invoice number is required")
	}

	if invoice.Amount <= 0 {
		return fmt.Errorf("invoice amount must be greater than zero")
	}

	duplicate, err := s.repository.FindDuplicate(
		ctx,
		invoice.InvoiceNumber,
		invoice.VendorName,
	)
	if err != nil {
		return err
	}

	if duplicate {
		return fmt.Errorf(
			"%w: %s from %s",
			ErrDuplicateInvoice,
			invoice.InvoiceNumber,
			invoice.VendorName,
		)
	}

	invoice.Status = "RECEIVED"

	if invoice.Amount <= autoApprovalLimit {
		invoice.Status = "APPROVED"
	} else {
		invoice.Status = "PENDING_APPROVAL"
	}

	if invoice.Currency == "" {
		invoice.Currency = "INR"
	}

	if invoice.Source == "" {
		invoice.Source = "api"
	}

	action := "SUBMITTED"

	if invoice.Status == "APPROVED" {
		action = "AUTO_APPROVED"
	}

	reason := fmt.Sprintf(
		"Invoice processed with amount %.2f",
		invoice.Amount,
	)

	if err := s.repository.CreateWithApprovalLog(
		ctx,
		invoice,
		action,
		nil,
		&reason,
	); err != nil {
		return err
	}

	return nil
}

func (s *InvoiceService) ApproveInvoice(
	ctx context.Context,
	invoiceID string,
	actor string,
) error {
	invoiceID = strings.TrimSpace(invoiceID)
	actor = strings.TrimSpace(actor)

	if invoiceID == "" {
		return ErrInvoiceIDRequired
	}

	if actor == "" {
		return fmt.Errorf("reviewer identity is required")
	}

	if err := s.repository.UpdateStatus(
		ctx,
		invoiceID,
		"APPROVED",
	); err != nil {
		if strings.Contains(err.Error(), "not pending approval") {
			return fmt.Errorf(
				"%w: %s",
				ErrInvoiceNotPending,
				invoiceID,
			)
		}

		return err
	}

	reason := "Invoice approved by human reviewer"

	if err := s.repository.CreateApprovalLog(
		ctx,
		invoiceID,
		"APPROVED",
		&actor,
		&reason,
	); err != nil {
		return err
	}

	return nil
}

func (s *InvoiceService) RejectInvoice(
	ctx context.Context,
	invoiceID string,
	actor string,
) error {
	invoiceID = strings.TrimSpace(invoiceID)
	actor = strings.TrimSpace(actor)

	if invoiceID == "" {
		return ErrInvoiceIDRequired
	}

	if actor == "" {
		return fmt.Errorf("reviewer identity is required")
	}

	if err := s.repository.UpdateStatus(
		ctx,
		invoiceID,
		"REJECTED",
	); err != nil {
		if strings.Contains(err.Error(), "not pending approval") {
			return fmt.Errorf(
				"%w: %s",
				ErrInvoiceNotPending,
				invoiceID,
			)
		}

		return err
	}

	reason := "Invoice rejected by human reviewer"

	if err := s.repository.CreateApprovalLog(
		ctx,
		invoiceID,
		"REJECTED",
		&actor,
		&reason,
	); err != nil {
		return err
	}

	return nil
}

func (s *InvoiceService) GetAllInvoices(
	ctx context.Context,
) ([]models.Invoice, error) {
	return s.repository.GetAll(ctx)
}

func (s *InvoiceService) GetInvoiceByID(
	ctx context.Context,
	invoiceID string,
) (*models.Invoice, error) {
	invoiceID = strings.TrimSpace(invoiceID)

	if invoiceID == "" {
		return nil, ErrInvoiceIDRequired
	}

	return s.repository.GetByID(ctx, invoiceID)
}
