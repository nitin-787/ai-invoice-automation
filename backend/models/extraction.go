package models

type ExtractedInvoice struct {
	InvoiceNumber string  `json:"invoice_number"`
	VendorName    string  `json:"vendor_name"`
	VendorEmail   *string `json:"vendor_email,omitempty"`
	InvoiceDate   *string `json:"invoice_date,omitempty"`
	DueDate       *string `json:"due_date,omitempty"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Description   *string `json:"description,omitempty"`
}
