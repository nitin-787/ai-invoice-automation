"use client";

import Link from "next/link";
import { ArrowUpRight, FileText } from "lucide-react";

import StatusBadge from "@/components/StatusBadge";
import type { Invoice } from "@/lib/api";

function formatCurrency(amount: number, currency: string) {
  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: currency || "INR",
    maximumFractionDigits: 0,
  }).format(amount);
}

function formatDate(date?: string) {
  if (!date) return "—";

  return new Date(date).toLocaleDateString("en-IN", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
}

export default function InvoiceTable({ invoices }: { invoices: Invoice[] }) {
  if (!invoices.length) {
    return (
      <div className="invoice-empty">
        <div className="invoice-empty-icon">
          <FileText size={22} />
        </div>

        <h3>No invoices found</h3>
        <p>Try changing your search or filter.</p>
      </div>
    );
  }

  return (
    <div className="invoice-table-wrapper">
      <table className="invoice-table">
        <thead>
          <tr>
            <th>Invoice</th>
            <th>Vendor</th>
            <th>Amount</th>
            <th>Status</th>
            <th>Source</th>
            <th>Created</th>
            <th></th>
          </tr>
        </thead>

        <tbody>
          {invoices.map((invoice) => (
            <tr key={invoice.id}>
              <td>
                <Link
                  href={`/invoices/${invoice.id}`}
                  className="invoice-main-link"
                >
                  <span className="invoice-number">
                    {invoice.invoice_number}
                  </span>
                  <span className="invoice-subtext">View invoice details</span>
                </Link>
              </td>

              <td>
                <span className="vendor-name">{invoice.vendor_name}</span>
              </td>

              <td>
                <span className="invoice-amount">
                  {formatCurrency(invoice.amount, invoice.currency)}
                </span>
              </td>

              <td>
                <StatusBadge status={invoice.status} />
              </td>

              <td>
                <span className="source-badge">{invoice.source || "API"}</span>
              </td>

              <td>
                <span className="created-date">
                  {formatDate(invoice.created_at)}
                </span>
              </td>

              <td>
                <Link
                  href={`/invoices/${invoice.id}`}
                  className="invoice-open"
                  aria-label={`Open ${invoice.invoice_number}`}
                >
                  <ArrowUpRight size={17} />
                </Link>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
