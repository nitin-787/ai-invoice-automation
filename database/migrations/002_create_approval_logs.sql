CREATE TABLE approval_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,

    action VARCHAR(30) NOT NULL,

    actor VARCHAR(255),
    reason TEXT,

    metadata JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT approval_action_check
        CHECK (
            action IN (
                'SUBMITTED',
                'AUTO_APPROVED',
                'APPROVED',
                'REJECTED',
                'ESCALATED'
            )
        )
);

CREATE INDEX idx_approval_logs_invoice
    ON approval_logs(invoice_id);

CREATE INDEX idx_approval_logs_created_at
    ON approval_logs(created_at);