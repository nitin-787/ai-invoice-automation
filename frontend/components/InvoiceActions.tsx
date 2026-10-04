"use client";

import { useState } from "react";
import { Check, X } from "lucide-react";
import { approveInvoice, rejectInvoice } from "@/lib/api";

export default function InvoiceActions({ id }: { id: string }) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [confirmAction, setConfirmAction] = useState<
    "approve" | "reject" | null
  >(null);

  async function handleAction() {
    if (!confirmAction) return;

    setLoading(true);
    setError("");

    try {
      if (confirmAction === "approve") {
        await approveInvoice(id);
      } else {
        await rejectInvoice(id);
      }

      window.location.reload();
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : `Failed to ${confirmAction} invoice`,
      );
      setLoading(false);
      setConfirmAction(null);
    }
  }

  return (
    <>
      <div className="invoice-actions">
        <button
          type="button"
          className="button button-danger"
          onClick={() => setConfirmAction("reject")}
          disabled={loading}
        >
          <X size={16} />
          Reject
        </button>

        <button
          type="button"
          className="button button-primary"
          onClick={() => setConfirmAction("approve")}
          disabled={loading}
        >
          <Check size={16} />
          Approve Invoice
        </button>

        {error && <div className="action-error">{error}</div>}
      </div>

      {confirmAction && (
        <div className="confirmation-overlay">
          <div className="confirmation-modal">
            <div className="confirmation-icon">
              {confirmAction === "approve" ? (
                <Check size={20} />
              ) : (
                <X size={20} />
              )}
            </div>

            <h3>
              {confirmAction === "approve"
                ? "Approve this invoice?"
                : "Reject this invoice?"}
            </h3>

            <p>
              {confirmAction === "approve"
                ? "This invoice will be marked as approved."
                : "This invoice will be marked as rejected."}
            </p>

            <div className="confirmation-actions">
              <button
                type="button"
                className="button button-secondary"
                onClick={() => setConfirmAction(null)}
                disabled={loading}
              >
                Cancel
              </button>

              <button
                type="button"
                className={
                  confirmAction === "approve"
                    ? "button button-primary"
                    : "button button-danger"
                }
                onClick={handleAction}
                disabled={loading}
              >
                {loading
                  ? "Processing..."
                  : confirmAction === "approve"
                    ? "Approve"
                    : "Reject"}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
