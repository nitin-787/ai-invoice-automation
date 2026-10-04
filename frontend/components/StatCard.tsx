import type { LucideIcon } from "lucide-react";

interface StatCardProps {
  label: string;
  value: string | number;
  description: string;
  icon: LucideIcon;
}

export default function StatCard({
  label,
  value,
  description,
  icon: Icon,
}: StatCardProps) {
  return (
    <div className="stat-card">
      <div className="stat-top">
        <span>{label}</span>

        <div className="stat-icon">
          <Icon size={18} />
        </div>
      </div>

      <div className="stat-value">{value}</div>

      <div className="stat-description">{description}</div>
    </div>
  );
}
