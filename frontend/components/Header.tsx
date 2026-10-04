"use client";

import { Bell, Search, UserCircle } from "lucide-react";

interface HeaderProps {
  title: string;
  description?: string;
}

export default function Header({ title, description }: HeaderProps) {
  return (
    <header className="topbar">
      <div>
        <h1>{title}</h1>
        {description && <p>{description}</p>}
      </div>

      <div className="topbar-actions">
        <div className="header-search">
          <Search size={17} />
          <input placeholder="Search..." />
          <kbd>⌘ K</kbd>
        </div>

        <button className="icon-button">
          <Bell size={18} />
          <span className="notification-dot" />
        </button>

        <div className="user">
          <div className="avatar">A</div>
          <div className="user-info">
            <strong>Admin</strong>
            <span>Reviewer</span>
          </div>
        </div>
      </div>
    </header>
  );
}
