CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    invoice_number VARCHAR(100) NOT NULL,
    vendor_name VARCHAR(255) NOT NULL,
    vendor_email VARCHAR(255),

    invoice_date DATE,
    due_date DATE,

    amount NUMERIC(15, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'INR',

    description TEXT,

    status VARCHAR(30) NOT NULL DEFAULT 'RECEIVED',

    source VARCHAR(50) NOT NULL DEFAULT 'email',
    source_reference TEXT,

    extracted_data JSONB,

    validation_errors JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT invoices_status_check
        CHECK (
            status IN (
                'RECEIVED',
                'PROCESSING',
                'PENDING_APPROVAL',
                'APPROVED',
                'REJECTED'
            )
        )
);

CREATE INDEX idx_invoices_status
    ON invoices(status);

CREATE INDEX idx_invoices_vendor
    ON invoices(vendor_name);

CREATE INDEX idx_invoices_invoice_number
    ON invoices(invoice_number);