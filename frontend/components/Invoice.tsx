"use client";

import Link from "next/link";
import { ArrowUpRight, MoreHorizontal } from "lucide-react";
import type { Invoice } from "@/lib/api";
import StatusBadge from "./StatusBadge";

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
      <div className="empty-state">
        <div className="empty-icon">
          <span>∅</span>
        </div>

        <h3>No invoices found</h3>
        <p>There are no invoices matching your criteria.</p>
      </div>
    );
  }

  return (
    <div className="table-wrapper">
      <table className="invoice-table">
        <thead>
          <tr>
            <th>Invoice</th>
            <th>Vendor</th>
            <th>Amount</th>
            <th>Status</th>
            <th>Source</th>
            <th>Date</th>
            <th />
          </tr>
        </thead>

        <tbody>
          {invoices.map((invoice) => (
            <tr key={invoice.id}>
              <td>
                <Link
                  href={`/invoices/${invoice.id}`}
                  className="invoice-number"
                >
                  {invoice.invoice_number}
                  <ArrowUpRight size={13} />
                </Link>
              </td>

              <td>
                <div className="vendor-cell">
                  <div className="vendor-avatar">
                    {invoice.vendor_name.charAt(0).toUpperCase()}
                  </div>

                  <div>
                    <strong>{invoice.vendor_name}</strong>
                    {invoice.vendor_email && (
                      <span>{invoice.vendor_email}</span>
                    )}
                  </div>
                </div>
              </td>

              <td className="amount">
                {formatCurrency(invoice.amount, invoice.currency)}
              </td>

              <td>
                <StatusBadge status={invoice.status} />
              </td>

              <td>
                <span className="source-badge">{invoice.source || "api"}</span>
              </td>

              <td className="date">{formatDate(invoice.created_at)}</td>

              <td>
                <Link href={`/invoices/${invoice.id}`} className="table-action">
                  <MoreHorizontal size={18} />
                </Link>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
