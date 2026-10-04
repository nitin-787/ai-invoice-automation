package services

import (
	"context"
	"fmt"

	"github.com/nitin-787/ai-invoice-automation/models"
	"github.com/nitin-787/ai-invoice-automation/repository"
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
			"duplicate invoice: %s from %s",
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

	if err := s.repository.Create(ctx, invoice); err != nil {
		return err
	}

	action := "SUBMITTED"

	if invoice.Status == "APPROVED" {
		action = "AUTO_APPROVED"
	}

	reason := fmt.Sprintf(
		"Invoice processed with amount %.2f",
		invoice.Amount,
	)

	if err := s.repository.CreateApprovalLog(
		ctx,
		invoice.ID,
		action,
		nil,
		&reason,
	); err != nil {
		return err
	}

	return nil
}
