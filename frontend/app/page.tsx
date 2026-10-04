import Link from "next/link";
import {
  ArrowRight,
  CheckCircle2,
  Clock3,
  Receipt,
  TrendingUp,
} from "lucide-react";

import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import StatCard from "@/components/StatCard";
import InvoiceTable from "@/components/InvoiceTable";
import { getInvoices } from "@/lib/api";

export default async function Dashboard() {
  let invoices = [];

  try {
    const response = await getInvoices();
    invoices = response.invoices;
  } catch {
    invoices = [];
  }

  const total = invoices.length;

  const approved = invoices.filter(
    (invoice) => invoice.status === "APPROVED",
  ).length;

  const pending = invoices.filter(
    (invoice) => invoice.status === "PENDING_APPROVAL",
  ).length;

  const totalValue = invoices.reduce(
    (sum, invoice) => sum + Number(invoice.amount || 0),
    0,
  );

  const recentInvoices = invoices.slice(0, 6);

  return (
    <div className="app-shell">
      <Sidebar />

      <main className="main">
        <Header
          title="Dashboard"
          description="Monitor invoice processing and approvals."
        />

        <div className="page-content">
          <section className="hero">
            <div>
              <div className="eyebrow">
                <span className="live-dot" />
                AUTOMATION ACTIVE
              </div>

              <h2>
                Invoice operations,
                <br />
                <span>simplified.</span>
              </h2>

              <p>
                AI-powered extraction, validation and human approval in one
                workflow.
              </p>
            </div>

            <Link href="/create" className="button button-primary">
              <Receipt size={17} />
              Create Invoice
            </Link>
          </section>

          <section className="stats-grid">
            <StatCard
              label="Total Invoices"
              value={total}
              description="All processed invoices"
              icon={Receipt}
            />

            <StatCard
              label="Approved"
              value={approved}
              description="Successfully approved"
              icon={CheckCircle2}
            />

            <StatCard
              label="Pending Review"
              value={pending}
              description="Require human approval"
              icon={Clock3}
            />

            <StatCard
              label="Invoice Value"
              value={`₹${totalValue.toLocaleString("en-IN")}`}
              description="Total processed value"
              icon={TrendingUp}
            />
          </section>

          <section className="content-card">
            <div className="section-header">
              <div>
                <h3>Recent invoices</h3>
                <p>Latest activity across your invoice pipeline.</p>
              </div>

              <Link href="/invoices" className="text-link">
                View all
                <ArrowRight size={15} />
              </Link>
            </div>

            <InvoiceTable invoices={recentInvoices} />
          </section>
        </div>
      </main>
    </div>
  );
}
