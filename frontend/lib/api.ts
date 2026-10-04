export const API_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export type InvoiceStatus =
  | "RECEIVED"
  | "APPROVED"
  | "PENDING_APPROVAL"
  | "REJECTED";

export interface Invoice {
  id: string;
  invoice_number: string;
  vendor_name: string;
  vendor_email?: string;
  invoice_date?: string;
  due_date?: string;
  amount: number;
  currency: string;
  description?: string;
  status: InvoiceStatus;
  source?: string;
  source_reference?: string;
  extracted_data?: Record<string, unknown>;
  validation_errors?: string[];
  created_at: string;
  updated_at: string;
}

export interface InvoiceResponse {
  invoices: Invoice[];
  count: number;
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    cache: "no-store",
  });

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`;

    try {
      const data = await response.json();

      message = data.details || data.error || data.message || message;
    } catch {
      // Ignore invalid JSON error response.
    }

    throw new Error(message);
  }

  return response.json();
}

export async function getInvoices(): Promise<InvoiceResponse> {
  return request<InvoiceResponse>("/api/v1/invoices");
}

export async function getInvoice(id: string): Promise<Invoice> {
  return request<Invoice>(`/api/v1/invoices/${id}`);
}

export async function createInvoice(data: {
  invoice_number: string;
  vendor_name: string;
  vendor_email?: string;
  invoice_date?: string;
  due_date?: string;
  amount: number;
  currency: string;
  description?: string;
}): Promise<Invoice> {
  const payload = {
    ...data,
    invoice_date: data.invoice_date
      ? `${data.invoice_date}T00:00:00Z`
      : undefined,
    due_date: data.due_date ? `${data.due_date}T00:00:00Z` : undefined,
  };

  return request<Invoice>("/api/v1/invoices", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function approveInvoice(id: string): Promise<void> {
  await request(`/api/v1/invoices/${id}/approve`, {
    method: "POST",
    headers: {
      "X-API-Key": "KnAQLqesC7BWCw4DVa7Y2kgqTnh1zGhz290FwHKhotU=",
    },
  });
}

export async function rejectInvoice(id: string): Promise<void> {
  await request(`/api/v1/invoices/${id}/reject`, {
    method: "POST",
    headers: {
      "X-API-Key": "KnAQLqesC7BWCw4DVa7Y2kgqTnh1zGhz290FwHKhotU=",
    },
  });
}
