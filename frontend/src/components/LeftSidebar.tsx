import { useState } from 'react'
import './LeftSidebar.css'

type NavItem = 'chats' | 'status' | 'communities' | 'calls' | 'settings'

interface LeftSidebarProps {
  activeItem?: NavItem
  onItemClick?: (item: NavItem) => void
  hideOnMobile?: boolean
}

export default function LeftSidebar({ activeItem = 'chats', onItemClick, hideOnMobile = false }: LeftSidebarProps) {
  const [statusBadge, setStatusBadge] = useState(1) // Example: 1 unread status

  const navItems: { id: NavItem; icon: string; label: string; badge?: number }[] = [
    { id: 'status', icon: '⚪', label: 'Status', badge: statusBadge },
    { id: 'chats', icon: '💬', label: 'Chats' },
    { id: 'communities', icon: '👥', label: 'Communities' },
    { id: 'calls', icon: '📞', label: 'Calls' },
    { id: 'settings', icon: '⚙️', label: 'Settings' },
  ]

  return (
    <div className={`left-sidebar ${hideOnMobile ? 'hide-on-mobile' : ''}`}>
      {navItems.map((item) => (
        <button
          key={item.id}
          className={`left-sidebar-item ${activeItem === item.id ? 'active' : ''}`}
          onClick={() => onItemClick?.(item.id)}
          title={item.label}
        >
          <div className="left-sidebar-icon">
            {item.id === 'status' ? (
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <circle
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="2"
                  fill="none"
                />
                <circle cx="12" cy="12" r="6" fill="currentColor" />
              </svg>
            ) : item.id === 'chats' ? (
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path
                  d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            ) : item.id === 'communities' ? (
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path
                  d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
                <circle cx="9" cy="7" r="4" stroke="currentColor" strokeWidth="2" />
                <path
                  d="M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            ) : item.id === 'calls' ? (
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path
                  d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            ) : (
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="3" stroke="currentColor" strokeWidth="2" />
                <path
                  d="M12 1v6m0 6v6M5.64 5.64l4.24 4.24m4.24 4.24l4.24 4.24M1 12h6m6 0h6M5.64 18.36l4.24-4.24m4.24-4.24l4.24-4.24"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                />
              </svg>
            )}
            {item.badge && item.badge > 0 && (
              <span className="left-sidebar-badge">{item.badge}</span>
            )}
          </div>
        </button>
      ))}
    </div>
  )
}

