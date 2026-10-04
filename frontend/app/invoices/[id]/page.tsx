import Link from "next/link";
import { ArrowLeft } from "lucide-react";

import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import StatusBadge from "@/components/StatusBadge";
import { getInvoice } from "@/lib/api";
import InvoiceActions from "@/components/InvoiceActions";

function formatCurrency(amount: number, currency: string) {
  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: currency || "INR",
    maximumFractionDigits: 2,
  }).format(amount);
}

function formatDate(date?: string) {
  if (!date) return "—";

  return new Date(date).toLocaleDateString("en-IN", {
    day: "2-digit",
    month: "long",
    year: "numeric",
  });
}

export default async function InvoiceDetailsPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  let invoice;

  try {
    invoice = await getInvoice(id);
  } catch {
    return (
      <div className="app-shell">
        <Sidebar />

        <main className="main">
          <Header title="Invoice not found" />

          <div className="page-content">
            <div className="content-card">
              <div className="empty-state">
                <h3>Invoice could not be found</h3>
                <p>The invoice may have been deleted or the ID is invalid.</p>

                <Link href="/invoices" className="button button-primary">
                  Back to invoices
                </Link>
              </div>
            </div>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="app-shell">
      <Sidebar />

      <main className="main">
        <Header
          title="Invoice details"
          description="Review invoice information and processing status."
        />

        <div className="page-content">
          <Link href="/invoices" className="back-link">
            <ArrowLeft size={16} />
            Back to invoices
          </Link>

          <div className="invoice-detail-header">
            <div>
              <div className="detail-eyebrow">INVOICE</div>

              <h2>{invoice.invoice_number}</h2>

              <div className="detail-meta">
                <StatusBadge status={invoice.status} />
                <span>Created {formatDate(invoice.created_at)}</span>
              </div>
            </div>
          </div>

          <InvoiceActions id={invoice.id} />

          <div className="detail-grid">
            <section className="content-card">
              <div className="section-header">
                <div>
                  <h3>Invoice information</h3>
                  <p>Core invoice and vendor information.</p>
                </div>
              </div>

              <div className="info-grid">
                <div className="info-item">
                  <div>
                    <span>Invoice number</span>
                    <strong>{invoice.invoice_number}</strong>
                  </div>
                </div>

                <div className="info-item">
                  <div>
                    <span>Vendor</span>
                    <strong>{invoice.vendor_name}</strong>
                  </div>
                </div>

                <div className="info-item">
                  <div>
                    <span>Vendor email</span>
                    <strong>{invoice.vendor_email || "—"}</strong>
                  </div>
                </div>

                <div className="info-item">
                  <div>
                    <span>Amount</span>
                    <strong className="amount-large">
                      {formatCurrency(invoice.amount, invoice.currency)}
                    </strong>
                  </div>
                </div>

                <div className="info-item">
                  <div>
                    <span>Invoice date</span>
                    <strong>{formatDate(invoice.invoice_date)}</strong>
                  </div>
                </div>

                <div className="info-item">
                  <div>
                    <span>Due date</span>
                    <strong>{formatDate(invoice.due_date)}</strong>
                  </div>
                </div>
              </div>
            </section>

            <section className="content-card">
              <div className="section-header">
                <div>
                  <h3>Processing</h3>
                  <p>Automation metadata.</p>
                </div>
              </div>

              <div className="processing-list">
                <div>
                  <span>Source</span>
                  <strong>{invoice.source || "API"}</strong>
                </div>

                <div>
                  <span>Status</span>
                  <StatusBadge status={invoice.status} />
                </div>

                <div>
                  <span>Created</span>
                  <strong>{formatDate(invoice.created_at)}</strong>
                </div>

                <div>
                  <span>Last updated</span>
                  <strong>{formatDate(invoice.updated_at)}</strong>
                </div>
              </div>
            </section>

            <section className="content-card detail-full">
              <div className="section-header">
                <div>
                  <h3>Description</h3>
                  <p>Invoice description and extracted information.</p>
                </div>
              </div>

              <div className="description-box">
                <span>{invoice.description || "No description provided."}</span>
              </div>
            </section>

            {invoice.extracted_data && (
              <section className="content-card detail-full">
                <div className="section-header">
                  <div>
                    <h3>AI extracted data</h3>
                    <p>Data extracted during invoice processing.</p>
                  </div>
                </div>

                <pre className="json-box">
                  {JSON.stringify(invoice.extracted_data, null, 2)}
                </pre>
              </section>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
