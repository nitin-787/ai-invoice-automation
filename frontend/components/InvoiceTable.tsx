"use client";

import Link from "next/link";
import { ArrowUpRight, FileText } from "lucide-react";
import type { Invoice } from "@/lib/api";

interface InvoiceTableProps {
  invoices: Invoice[];
}

export default function InvoiceTable({ invoices }: InvoiceTableProps) {
  if (!invoices || invoices.length === 0) {
    return (
      <div className="rounded-2xl border border-slate-200 bg-white p-10 text-center">
        <FileText className="mx-auto mb-3 h-8 w-8 text-slate-400" />

        <h3 className="text-sm font-semibold text-slate-900">
          No invoices found
        </h3>

        <p className="mt-1 text-sm text-slate-500">
          Invoices will appear here once they are created.
        </p>
      </div>
    );
  }

  return (
    <div className="w-full overflow-hidden rounded-2xl border border-slate-200 bg-white">
      <div className="overflow-x-auto">
        <table className="w-full min-w-[950px] border-collapse">
          <thead>
            <tr className="border-b border-slate-200 bg-slate-50">
              <th className="px-5 py-4 text-left text-xs font-semibold uppercase tracking-wide text-slate-500">
                Invoice
              </th>

              <th className="px-5 py-4 text-left text-xs font-semibold uppercase tracking-wide text-slate-500">
                Vendor
              </th>

              <th className="px-5 py-4 text-right text-xs font-semibold uppercase tracking-wide text-slate-500">
                Amount
              </th>

              <th className="px-5 py-4 text-left text-xs font-semibold uppercase tracking-wide text-slate-500">
                Status
              </th>

              <th className="px-5 py-4 text-left text-xs font-semibold uppercase tracking-wide text-slate-500">
                Source
              </th>

              <th className="px-5 py-4 text-left text-xs font-semibold uppercase tracking-wide text-slate-500">
                Created
              </th>

              <th className="w-12 px-3 py-4" />
            </tr>
          </thead>

          <tbody className="divide-y divide-slate-100">
            {invoices.map((invoice) => (
              <tr
                key={invoice.id}
                className="group transition-colors hover:bg-slate-50"
              >
                {/* Invoice */}
                <td className="px-5 py-4">
                  <Link
                    href={`/invoices/${invoice.id}`}
                    className="block min-w-0"
                  >
                    <div
                      className="max-w-[220px] truncate text-sm font-semibold text-slate-900"
                      title={invoice.invoice_number}
                    >
                      {invoice.invoice_number}
                    </div>

                    <div
                      className="mt-1 max-w-[220px] truncate font-mono text-[11px] text-slate-400"
                      title={invoice.id}
                    >
                      {invoice.id}
                    </div>
                  </Link>
                </td>

                {/* Vendor */}
                <td className="px-5 py-4">
                  <div
                    className="max-w-[220px] truncate text-sm font-medium text-slate-700"
                    title={invoice.vendor_name}
                  >
                    {invoice.vendor_name}
                  </div>
                </td>

                {/* Amount */}
                <td className="whitespace-nowrap px-5 py-4 text-right">
                  <span className="text-sm font-semibold text-slate-900">
                    {formatCurrency(invoice.amount, invoice.currency)}
                  </span>
                </td>

                {/* Status */}
                <td className="px-5 py-4">
                  <StatusPill status={invoice.status} />
                </td>

                {/* Source */}
                <td className="px-5 py-4">
                  <span className="inline-flex rounded-md bg-slate-100 px-2.5 py-1 text-xs font-medium text-slate-600">
                    {invoice.source || "api"}
                  </span>
                </td>

                {/* Created */}
                <td className="whitespace-nowrap px-5 py-4 text-sm text-slate-500">
                  {formatDate(invoice.created_at)}
                </td>

                {/* View */}
                <td className="px-3 py-4 text-right">
                  <Link
                    href={`/invoices/${invoice.id}`}
                    aria-label={`View ${invoice.invoice_number}`}
                    className="inline-flex h-8 w-8 items-center justify-center rounded-lg text-slate-400 transition hover:bg-slate-100 hover:text-slate-900"
                  >
                    <ArrowUpRight className="h-4 w-4" />
                  </Link>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function StatusPill({ status }: { status: string }) {
  const normalized = status?.toUpperCase();

  let className = "inline-flex rounded-full px-2.5 py-1 text-xs font-semibold";

  if (normalized === "APPROVED") {
    className += " bg-emerald-50 text-emerald-700";
  } else if (normalized === "PENDING_APPROVAL") {
    className += " bg-amber-50 text-amber-700";
  } else if (normalized === "REJECTED") {
    className += " bg-red-50 text-red-700";
  } else if (normalized === "RECEIVED") {
    className += " bg-blue-50 text-blue-700";
  } else {
    className += " bg-slate-100 text-slate-600";
  }

  const label =
    normalized === "PENDING_APPROVAL"
      ? "Pending Review"
      : normalized || "Unknown";

  return <span className={className}>{label}</span>;
}

function formatCurrency(amount: number, currency?: string): string {
  try {
    return new Intl.NumberFormat("en-IN", {
      style: "currency",
      currency: currency || "INR",
      maximumFractionDigits: 2,
    }).format(amount);
  } catch {
    return `${currency || "INR"} ${amount.toLocaleString("en-IN")}`;
  }
}

function formatDate(value: string): string {
  if (!value) {
    return "—";
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat("en-IN", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  }).format(date);
}
