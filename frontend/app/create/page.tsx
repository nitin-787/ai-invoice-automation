"use client";

import { FormEvent, useState } from "react";
import { ArrowLeft, Save, Loader2 } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import { createInvoice } from "@/lib/api";

export default function CreateInvoicePage() {
  const router = useRouter();

  const [form, setForm] = useState({
    invoice_number: "",
    vendor_name: "",
    vendor_email: "",
    invoice_date: "",
    due_date: "",
    amount: "",
    currency: "INR",
    description: "",
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  function update(key: keyof typeof form, value: string) {
    setForm((current) => ({
      ...current,
      [key]: value,
    }));
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setLoading(true);
    setError("");

    try {
      const invoice = await createInvoice({
        invoice_number: form.invoice_number,
        vendor_name: form.vendor_name,
        vendor_email: form.vendor_email || undefined,
        invoice_date: form.invoice_date || undefined,
        due_date: form.due_date || undefined,
        amount: Number(form.amount),
        currency: form.currency,
        description: form.description || undefined,
      });

      router.push(`/invoices/${invoice.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create invoice");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="app-shell">
      {/* <Sidebar /> */}

      <main className="main">
        {/* <Header
          title="Create invoice"
          description="Add a new invoice to the automation pipeline."
        /> */}

        <div className="page-content">
          <Link href="/invoices" className="back-link">
            <ArrowLeft size={16} />
            Back to invoices
          </Link>

          <form className="content-card create-form" onSubmit={handleSubmit}>
            <div className="section-header">
              <div>
                <h3>Invoice details</h3>
                <p>Enter the invoice information below.</p>
              </div>
            </div>

            {error && <div className="error-box">{error}</div>}

            <div className="form-grid">
              <label className="field">
                <span>Invoice number *</span>
                <input
                  required
                  value={form.invoice_number}
                  onChange={(e) => update("invoice_number", e.target.value)}
                  placeholder="INV-2026-010"
                />
              </label>

              <label className="field">
                <span>Vendor name *</span>
                <input
                  required
                  value={form.vendor_name}
                  onChange={(e) => update("vendor_name", e.target.value)}
                  placeholder="Acme Corporation"
                />
              </label>

              <label className="field">
                <span>Vendor email</span>
                <input
                  type="email"
                  value={form.vendor_email}
                  onChange={(e) => update("vendor_email", e.target.value)}
                  placeholder="finance@acme.com"
                />
              </label>

              <label className="field">
                <span>Currency</span>
                <select
                  value={form.currency}
                  onChange={(e) => update("currency", e.target.value)}
                >
                  <option value="INR">INR — Indian Rupee</option>
                  <option value="USD">USD — US Dollar</option>
                  <option value="EUR">EUR — Euro</option>
                  <option value="GBP">GBP — Pound</option>
                </select>
              </label>

              <label className="field">
                <span>Amount *</span>
                <input
                  required
                  type="number"
                  min="0.01"
                  step="0.01"
                  value={form.amount}
                  onChange={(e) => update("amount", e.target.value)}
                  placeholder="85000"
                />
              </label>

              <label className="field">
                <span>Invoice date</span>
                <input
                  type="date"
                  value={form.invoice_date}
                  onChange={(e) => update("invoice_date", e.target.value)}
                />
              </label>

              <label className="field">
                <span>Due date</span>
                <input
                  type="date"
                  value={form.due_date}
                  onChange={(e) => update("due_date", e.target.value)}
                />
              </label>

              <label className="field field-full">
                <span>Description</span>
                <textarea
                  rows={5}
                  value={form.description}
                  onChange={(e) => update("description", e.target.value)}
                  placeholder="Describe the goods or services..."
                />
              </label>
            </div>

            <div className="form-footer">
              <Link href="/invoices" className="button button-secondary">
                Cancel
              </Link>

              <button
                type="submit"
                className="button button-primary"
                disabled={loading}
              >
                {loading ? (
                  <Loader2 size={17} className="spin" />
                ) : (
                  <Save size={17} />
                )}

                {loading ? "Processing..." : "Create Invoice"}
              </button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}
