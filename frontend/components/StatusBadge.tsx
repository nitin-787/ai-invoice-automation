import type { InvoiceStatus } from "@/lib/api";

interface StatusBadgeProps {
  status: InvoiceStatus;
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  const labels: Record<InvoiceStatus, string> = {
    RECEIVED: "Received",
    APPROVED: "Approved",
    PENDING_APPROVAL: "Pending Review",
    REJECTED: "Rejected",
  };

  return (
    <span className={`status-badge status-${status.toLowerCase()}`}>
      <span className="status-dot" />
      {labels[status]}
    </span>
  );
}
