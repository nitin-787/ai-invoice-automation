package models

import "time"

type Invoice struct {
	ID               string     `json:"id"`
	InvoiceNumber    string     `json:"invoice_number"`
	VendorName       string     `json:"vendor_name"`
	VendorEmail      *string    `json:"vendor_email,omitempty"`
	InvoiceDate      *time.Time `json:"invoice_date,omitempty"`
	DueDate          *time.Time `json:"due_date,omitempty"`
	Amount           float64    `json:"amount"`
	Currency         string     `json:"currency"`
	Description      *string    `json:"description,omitempty"`
	Status           string     `json:"status"`
	Source           string     `json:"source"`
	SourceReference  *string    `json:"source_reference,omitempty"`
	ExtractedData    any        `json:"extracted_data,omitempty"`
	ValidationErrors any        `json:"validation_errors,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
