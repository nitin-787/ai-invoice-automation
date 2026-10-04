"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import {
  ArrowRight,
  CheckCircle2,
  Clock3,
  FileText,
  Plus,
  Receipt,
  XCircle,
} from "lucide-react";

import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import StatCard from "@/components/StatCard";
import StatusBadge from "@/components/StatusBadge";
import { getInvoices, type Invoice } from "@/lib/api";

export default function DashboardPage() {
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadInvoices() {
      try {
        const response = await getInvoices();
        setInvoices(response.invoices);
      } catch (error) {
        console.error("Failed to load invoices:", error);
      } finally {
        setLoading(false);
      }
    }

    loadInvoices();
  }, []);

  const total = invoices.length;

  const approved = invoices.filter(
    (invoice) => invoice.status === "APPROVED",
  ).length;

  const pending = invoices.filter(
    (invoice) => invoice.status === "PENDING_APPROVAL",
  ).length;

  const rejected = invoices.filter(
    (invoice) => invoice.status === "REJECTED",
  ).length;

  const totalAmount = invoices.reduce(
    (sum, invoice) => sum + Number(invoice.amount || 0),
    0,
  );

  const recentInvoices = [...invoices]
    .sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    )
    .slice(0, 5);

  function formatCurrency(amount: number, currency = "INR") {
    return new Intl.NumberFormat("en-IN", {
      style: "currency",
      currency,
      maximumFractionDigits: 0,
    }).format(amount);
  }

  function formatDate(date: string) {
    return new Date(date).toLocaleDateString("en-IN", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    });
  }

  return (
    <div className="app-shell">
      {/* <Sidebar /> */}

      <main className="main">
        {/* <Header
          title="Dashboard"
          description="Monitor invoices, approvals, and automation activity."
        /> */}

        <div className="page-content">
          <div className="dashboard-hero">
            <div>
              <span className="detail-eyebrow">INVOICE AUTOMATION</span>
              <h2>Everything under control.</h2>
              <p>
                Track your invoices and let automation handle the repetitive
                approval work.
              </p>
            </div>

            <Link href="/create" className="button button-primary">
              <Plus size={17} />
              Create invoice
            </Link>
          </div>

          <div className="stats-grid">
            <StatCard title="Total invoices" value={total} icon={Receipt} />

            <StatCard title="Approved" value={approved} icon={CheckCircle2} />

            <StatCard title="Pending approval" value={pending} icon={Clock3} />

            <StatCard title="Rejected" value={rejected} icon={XCircle} />
          </div>

          <div className="dashboard-grid">
            <section className="content-card dashboard-main-card">
              <div className="section-header">
                <div>
                  <h3>Recent invoices</h3>
                  <p>Your latest invoice activity.</p>
                </div>

                <Link href="/invoices" className="text-link">
                  View all
                  <ArrowRight size={15} />
                </Link>
              </div>

              {loading ? (
                <div className="empty-state">
                  <p>Loading invoices...</p>
                </div>
              ) : recentInvoices.length === 0 ? (
                <div className="empty-state">
                  <FileText size={30} />
                  <h3>No invoices yet</h3>
                  <p>Create your first invoice to get started.</p>

                  <Link href="/create" className="button button-primary">
                    <Plus size={16} />
                    Create invoice
                  </Link>
                </div>
              ) : (
                <div className="recent-invoice-list">
                  {recentInvoices.map((invoice) => (
                    <Link
                      key={invoice.id}
                      href={`/invoices/${invoice.id}`}
                      className="recent-invoice-row"
                    >
                      <div className="recent-invoice-icon">
                        <FileText size={18} />
                      </div>

                      <div className="recent-invoice-info">
                        <strong>{invoice.invoice_number}</strong>
                        <span>{invoice.vendor_name}</span>
                      </div>

                      <div className="recent-invoice-amount">
                        <strong>
                          {formatCurrency(invoice.amount, invoice.currency)}
                        </strong>
                        <span>{formatDate(invoice.created_at)}</span>
                      </div>

                      <StatusBadge status={invoice.status} />

                      <ArrowRight size={16} className="recent-invoice-arrow" />
                    </Link>
                  ))}
                </div>
              )}
            </section>

            <section className="content-card dashboard-summary-card">
              <div className="section-header">
                <div>
                  <h3>Overview</h3>
                  <p>Current invoice volume.</p>
                </div>
              </div>

              <div className="overview-amount">
                <span>Total invoice value</span>
                <strong>{formatCurrency(totalAmount)}</strong>
              </div>

              <div className="overview-list">
                <div>
                  <span>
                    <i className="overview-dot approved-dot" />
                    Approved
                  </span>
                  <strong>{approved}</strong>
                </div>

                <div>
                  <span>
                    <i className="overview-dot pending-dot" />
                    Pending
                  </span>
                  <strong>{pending}</strong>
                </div>

                <div>
                  <span>
                    <i className="overview-dot rejected-dot" />
                    Rejected
                  </span>
                  <strong>{rejected}</strong>
                </div>
              </div>

              <Link
                href="/invoices"
                className="button button-secondary full-width"
              >
                Manage invoices
                <ArrowRight size={16} />
              </Link>
            </section>
          </div>
        </div>
      </main>
    </div>
  );
}
