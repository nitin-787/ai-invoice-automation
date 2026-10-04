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
    } catch (error) {
      console.error("Failed to load invoices:", error);
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

      const query = search.trim().toLowerCase();

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
      {/* <Sidebar /> */}

      <main className="main">
        {/* <Header
          title="Invoices"
          description="Manage and review your invoice pipeline."
        /> */}

        <div className="invoice-page">
          <div className="page-content">
            {/* Page intro */}
            <div className="invoices-page-intro">
              <div>
                <span className="detail-eyebrow">INVOICE MANAGEMENT</span>
                <h2>Your invoice pipeline</h2>
                <p>
                  Search, filter, review, and manage all your invoices from one
                  place.
                </p>
              </div>

              <Link href="/create" className="button button-primary">
                <Plus size={17} />
                New invoice
              </Link>
            </div>

            {/* Toolbar */}
            <div className="invoice-toolbar">
              <div className="search-box invoice-search">
                <Search size={17} />

                <input
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  placeholder="Search by invoice number or vendor..."
                />

                {search && (
                  <button
                    type="button"
                    className="search-clear"
                    onClick={() => setSearch("")}
                    aria-label="Clear search"
                  >
                    ×
                  </button>
                )}
              </div>

              <button
                type="button"
                className="button button-secondary"
                onClick={loadInvoices}
                disabled={loading}
              >
                <RefreshCw size={16} className={loading ? "spin" : ""} />
                Refresh
              </button>
            </div>

            {/* Filters */}
            <div className="invoice-filters">
              <div className="filter-title">
                <SlidersHorizontal size={15} />
                Status
              </div>

              <div className="filter-buttons">
                {filters.map((item) => (
                  <button
                    key={item}
                    type="button"
                    className={`filter-button ${filter === item ? "active" : ""}`}
                    onClick={() => setFilter(item)}
                  >
                    {item === "ALL"
                      ? "All"
                      : item
                          .replaceAll("_", " ")
                          .toLowerCase()
                          .replace(/\b\w/g, (char) => char.toUpperCase())}
                  </button>
                ))}
              </div>
            </div>

            {/* Table */}
            <div className="content-card invoices-table-card">
              <div className="section-header invoices-section-header">
                <div>
                  <h3>All invoices</h3>
                  <p>
                    {search || filter !== "ALL"
                      ? `Showing ${filtered.length} matching ${
                          filtered.length === 1 ? "invoice" : "invoices"
                        }.`
                      : `${invoices.length} ${
                          invoices.length === 1 ? "invoice" : "invoices"
                        } in your pipeline.`}
                  </p>
                </div>

                <div className="invoice-count">{filtered.length}</div>
              </div>

              {loading ? (
                <div className="loading-state">
                  <RefreshCw className="spin" size={24} />
                  <span>Loading invoices...</span>
                </div>
              ) : filtered.length === 0 ? (
                <div className="empty-state invoice-empty-state">
                  <div className="empty-state-icon">
                    <Search size={22} />
                  </div>

                  <h3>
                    {search || filter !== "ALL"
                      ? "No invoices found"
                      : "No invoices yet"}
                  </h3>

                  <p>
                    {search || filter !== "ALL"
                      ? "Try changing your search or filter."
                      : "Create your first invoice to start building your pipeline."}
                  </p>

                  {search || filter !== "ALL" ? (
                    <button
                      type="button"
                      className="button button-secondary"
                      onClick={() => {
                        setSearch("");
                        setFilter("ALL");
                      }}
                    >
                      Clear filters
                    </button>
                  ) : (
                    <Link href="/create" className="button button-primary">
                      <Plus size={16} />
                      Create invoice
                    </Link>
                  )}
                </div>
              ) : (
                <InvoiceTable invoices={filtered} />
              )}
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
