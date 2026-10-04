"use client";

import { useEffect, useMemo, useState } from "react";
import { Search, SlidersHorizontal, RefreshCw, Plus } from "lucide-react";
import Link from "next/link";

import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import InvoiceTable from "@/components/InvoiceTable";
import { getInvoices, type Invoice, type InvoiceStatus } from "@/lib/api";

const filters: Array<"ALL" | InvoiceStatus> = [
  "ALL",
  "PENDING_APPROVAL",
  "APPROVED",
  "REJECTED",
  "RECEIVED",
];

export default function InvoicesPage() {
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState<"ALL" | InvoiceStatus>("ALL");

  async function loadInvoices() {
    setLoading(true);

    try {
      const data = await getInvoices();
      setInvoices(data.invoices);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadInvoices();
  }, []);

  const filtered = useMemo(() => {
    return invoices.filter((invoice) => {
      const matchesStatus = filter === "ALL" || invoice.status === filter;

      const query = search.toLowerCase();

      const matchesSearch =
        !query ||
        invoice.invoice_number.toLowerCase().includes(query) ||
        invoice.vendor_name.toLowerCase().includes(query) ||
        invoice.status.toLowerCase().includes(query);

      return matchesStatus && matchesSearch;
    });
  }, [invoices, search, filter]);

  return (
    <div className="app-shell">
      <Sidebar />

      <main className="main">
        <Header
          title="Invoices"
          description="Manage and review your invoice pipeline."
        />

        <div className="page-content">
          <div className="page-toolbar">
            <div className="search-box">
              <Search size={17} />
              <input
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Search invoices, vendors..."
              />
            </div>

            <button
              className="button button-secondary"
              onClick={loadInvoices}
              disabled={loading}
            >
              <RefreshCw size={16} className={loading ? "spin" : ""} />
              Refresh
            </button>

            <Link href="/create" className="button button-primary">
              <Plus size={17} />
              New Invoice
            </Link>
          </div>

          <div className="filter-bar">
            <div className="filter-title">
              <SlidersHorizontal size={15} />
              Filter
            </div>

            {filters.map((item) => (
              <button
                key={item}
                className={`filter-button ${filter === item ? "active" : ""}`}
                onClick={() => setFilter(item)}
              >
                {item === "ALL" ? "All" : item.replace("_", " ")}
              </button>
            ))}
          </div>

          <div className="content-card">
            <div className="section-header">
              <div>
                <h3>All invoices</h3>
                <p>
                  Showing {filtered.length} of {invoices.length} invoices.
                </p>
              </div>
            </div>

            {loading ? (
              <div className="loading-state">
                <RefreshCw className="spin" size={24} />
                <span>Loading invoices...</span>
              </div>
            ) : (
              <InvoiceTable invoices={filtered} />
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
