"use client";

import { useState } from "react";
import { Check, X } from "lucide-react";
import { approveInvoice, rejectInvoice } from "@/lib/api";

export default function InvoiceActions({ id }: { id: string }) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleApprove() {
    setLoading(true);
    setError("");

    try {
      await approveInvoice(id);
      window.location.reload();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to approve invoice",
      );
      setLoading(false);
    }
  }

  async function handleReject() {
    setLoading(true);
    setError("");

    try {
      await rejectInvoice(id);
      window.location.reload();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to reject invoice");
      setLoading(false);
    }
  }

  return (
    <div className="invoice-actions">
      <button
        type="button"
        className="button button-danger"
        onClick={handleReject}
        disabled={loading}
      >
        <X size={16} />
        Reject
      </button>

      <button
        type="button"
        className="button button-primary"
        onClick={handleApprove}
        disabled={loading}
      >
        <Check size={16} />
        {loading ? "Processing..." : "Approve Invoice"}
      </button>

      {error && <div className="action-error">{error}</div>}
    </div>
  );
}
