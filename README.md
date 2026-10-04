# AI Invoice Automation

An AI-powered invoice processing and approval system built with Go, PostgreSQL, AI extraction, and n8n automation.

The system automates invoice ingestion, validation, duplicate detection, approval decisions, and human approval workflows while maintaining an audit trail.

## Features

- Invoice creation through REST API
- PostgreSQL persistence
- Invoice validation and normalization
- Duplicate invoice detection
- Automatic approval for invoices up to ₹50,000
- Human approval workflow for invoices above ₹50,000
- Invoice approval/rejection
- Reviewer identity tracking
- Approval audit logs
- AI-based invoice data extraction
- PDF invoice extraction
- n8n workflow automation
- Invoice retrieval and status tracking

## Architecture

```text
                 ┌─────────────────┐
                 │   Invoice Input │
                 │   API / PDF     │
                 └────────┬────────┘
                          │
                          ▼
                 ┌─────────────────┐
                 │  AI Extraction  │
                 └────────┬────────┘
                          │
                          ▼
                 ┌─────────────────┐
                 │   Validation    │
                 └────────┬────────┘
                          │
                          ▼
                 ┌─────────────────┐
                 │ Duplicate Check │
                 └────────┬────────┘
                          │
                 ┌────────┴────────┐
                 │                 │
             ≤ ₹50,000          > ₹50,000
                 │                 │
                 ▼                 ▼
          AUTO APPROVED     PENDING APPROVAL
                                   │
                                   ▼
                            Human Reviewer
                              /        \
                         APPROVE       REJECT
                              \        /
                               ▼      ▼
                              Audit Log
```

## Tech Stack

- **Backend:** Go
- **Web Framework:** Gin
- **Database:** PostgreSQL
- **Database Driver:** pgx
- **AI:** AI-powered invoice extraction
- **Automation:** n8n
- **API:** REST
- **Containerization:** Docker

## API Endpoints

### Health

```http
GET /health
```

### Invoices

```http
POST /api/v1/invoices
GET  /api/v1/invoices
GET  /api/v1/invoices/:id
```

### AI Extraction

```http
POST /api/v1/invoices/extract
POST /api/v1/invoices/extract/file
```

### Approval

```http
POST /api/v1/invoices/:id/approve
POST /api/v1/invoices/:id/reject
```

## Approval Logic

Invoices are automatically classified based on their amount.

| Invoice Amount | Status |
|---|---|
| ≤ ₹50,000 | `APPROVED` |
| > ₹50,000 | `PENDING_APPROVAL` |

Invoices requiring human review can subsequently be:

```text
PENDING_APPROVAL → APPROVED
PENDING_APPROVAL → REJECTED
```

Every submission and approval decision is recorded in the approval audit log.

## Duplicate Detection

An invoice is considered a duplicate when the combination of:

- Invoice number
- Vendor name

already exists in the database.

Duplicate submissions return:

```http
409 Conflict
```

Example:

```json
{
  "error": "duplicate invoice: DUP-TEST-001 from Duplicate Test Pvt Ltd"
}
```

## Project Structure

```text
backend/
├── db/
├── handlers/
├── middleware/
├── models/
├── repository/
├── services/
├── tests/
├── main.go
├── go.mod
└── go.sum
```

## Running the Backend

### Prerequisites

- Go
- PostgreSQL
- Docker
- n8n (for automation workflows)

### Start PostgreSQL

The project uses PostgreSQL for invoice persistence.

Make sure the PostgreSQL container/database is running before starting the API.

### Start the API

```powershell
cd backend

go build -o invoice-api.exe .

.\invoice-api.exe
```

The API runs on:

```text
http://localhost:8080
```

### Health Check

```powershell
curl.exe http://localhost:8080/health
```

### Get Invoices

```powershell
curl.exe http://localhost:8080/api/v1/invoices
```

## Example Invoice

```json
{
  "invoice_number": "INV-2026-001",
  "vendor_name": "ABC Construction Pvt Ltd",
  "amount": 85000,
  "currency": "INR"
}
```

A ₹85,000 invoice will enter:

```text
PENDING_APPROVAL
```

A ₹42,000 invoice will be:

```text
APPROVED
```

## Testing

Run the complete Go test suite:

```powershell
go test ./...
```

Format the code:

```powershell
gofmt -w .
```

## Workflow Automation

n8n is used to automate invoice processing and human approval workflows.

The automation can:

1. Receive invoice data
2. Trigger invoice processing
3. Perform AI extraction
4. Create the invoice
5. Check the approval status
6. Wait for human approval when required
7. Approve or reject the invoice
8. Preserve the approval history

## Audit Trail

Invoice actions are recorded in the approval log, including:

- Invoice submission
- Automatic approval
- Human approval
- Human rejection
- Reviewer identity
- Reason for the action

Example:

```text
SUBMITTED
    ↓
APPROVED
    ↓
actor: api-reviewer
```

## Status

The core invoice automation backend and approval workflow are implemented and tested.