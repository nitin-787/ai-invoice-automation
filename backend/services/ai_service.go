package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/nitin-787/ai-invoice-automation/models"
	"google.golang.org/genai"
)

type AIService struct {
	client *genai.Client
}

func NewAIService(ctx context.Context) (*AIService, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not configured")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &AIService{
		client: client,
	}, nil
}

func (s *AIService) ExtractInvoice(
	ctx context.Context,
	input string,
) (*models.ExtractedInvoice, error) {
	if input == "" {
		return nil, fmt.Errorf("invoice input is empty")
	}

	prompt := `Extract invoice information from the following text.

Return ONLY valid JSON matching this structure:

{
  "invoice_number": "string",
  "vendor_name": "string",
  "vendor_email": "string or null",
  "invoice_date": "YYYY-MM-DD or null",
  "due_date": "YYYY-MM-DD or null",
  "amount": 0,
  "currency": "string",
  "description": "string or null"
}

Rules:
- Do not invent information.
- Use null when a field is unavailable.
- amount must be a number.
- Return JSON only.

Invoice text:
` + input

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	}

	result, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("Gemini extraction failed: %w", err)
	}

	if len(result.Candidates) == 0 ||
		result.Candidates[0].Content == nil ||
		len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("Gemini returned an empty response")
	}

	responseText := result.Candidates[0].Content.Parts[0].Text

	var invoice models.ExtractedInvoice

	if err := json.Unmarshal([]byte(responseText), &invoice); err != nil {
		return nil, fmt.Errorf(
			"failed to parse Gemini response: %w; response: %s",
			err,
			responseText,
		)
	}

	if invoice.InvoiceNumber == "" {
		return nil, fmt.Errorf("invoice number was not extracted")
	}

	if invoice.VendorName == "" {
		return nil, fmt.Errorf("vendor name was not extracted")
	}

	if invoice.Amount <= 0 {
		return nil, fmt.Errorf("invoice amount was not extracted")
	}

	return &invoice, nil
}

func (s *AIService) ExtractInvoicePDF(
	ctx context.Context,
	pdfData []byte,
) (*models.ExtractedInvoice, error) {
	if len(pdfData) == 0 {
		return nil, fmt.Errorf("invoice PDF is empty")
	}

	prompt := `Extract invoice information from this PDF.

Return ONLY valid JSON matching this structure:

{
  "invoice_number": "string",
  "vendor_name": "string",
  "vendor_email": "string or null",
  "invoice_date": "YYYY-MM-DD or null",
  "due_date": "YYYY-MM-DD or null",
  "amount": 0,
  "currency": "string",
  "description": "string or null"
}

Rules:
- Do not invent information.
- Use null when a field is unavailable.
- amount must be a number.
- Return JSON only.`

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	}

	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				genai.NewPartFromText(prompt),
				genai.NewPartFromBytes(pdfData, "application/pdf"),
			},
		},
	}

	result, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		config,
	)

	if err != nil {
		return nil, fmt.Errorf("Gemini PDF extraction failed: %w", err)
	}

	if len(result.Candidates) == 0 ||
		result.Candidates[0].Content == nil ||
		len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("Gemini returned an empty response")
	}

	responseText := result.Candidates[0].Content.Parts[0].Text

	var invoice models.ExtractedInvoice

	if err := json.Unmarshal([]byte(responseText), &invoice); err != nil {
		return nil, fmt.Errorf(
			"failed to parse Gemini PDF response: %w; response: %s",
			err,
			responseText,
		)
	}

	if invoice.InvoiceNumber == "" {
		return nil, fmt.Errorf("invoice number was not extracted")
	}

	if invoice.VendorName == "" {
		return nil, fmt.Errorf("vendor name was not extracted")
	}

	if invoice.Amount <= 0 {
		return nil, fmt.Errorf("invoice amount was not extracted")
	}

	return &invoice, nil
}
