import { useState } from 'react'
import { Conversation } from '../services/chatService'
import ChatList from './ChatList'
import './Sidebar.css'

interface SidebarProps {
  conversations: Conversation[]
  selectedConversation: Conversation | null
  onSelectConversation: (conversation: Conversation) => void
  currentUserId: number
  currentUsername: string
  isConnected: boolean
  onNewChat: () => void
  onLogout: () => void
}

type FilterType = 'all' | 'unread' | 'favorites' | 'groups'

export default function Sidebar({
  conversations,
  selectedConversation,
  onSelectConversation,
  currentUserId,
  currentUsername,
  isConnected,
  onNewChat,
  onLogout,
}: SidebarProps) {
  const [searchQuery, setSearchQuery] = useState('')
  const [activeFilter, setActiveFilter] = useState<FilterType>('all')
  const [showMenu, setShowMenu] = useState(false)

  // Filter conversations based on active filter
  const filteredConversations = conversations.filter((conv) => {
    if (searchQuery.trim()) {
      const participants = Array.isArray(conv.participants) ? conv.participants : []
      const otherParticipant = participants.find((p) => p.id !== currentUserId)
      const searchLower = searchQuery.toLowerCase()
      const matchesSearch =
        otherParticipant?.username.toLowerCase().includes(searchLower) ||
        otherParticipant?.email.toLowerCase().includes(searchLower) ||
        conv.last_message?.content?.toLowerCase().includes(searchLower)

      if (!matchesSearch) return false
    }

    switch (activeFilter) {
      case 'unread':
        return (conv.unread_count || 0) > 0
      case 'groups':
        return conv.type === 'group'
      case 'favorites':
        // TODO: Implement favorites
        return false
      default:
        return true
    }
  })

  return (
    <div className="whatsapp-sidebar">
      {/* Header */}
      <div className="sidebar-header">
        <div className="sidebar-header-left">
          <button className="menu-button" onClick={() => setShowMenu(!showMenu)}>
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path
                d="M3 12H21M3 6H21M3 18H21"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </button>
            <div className="app-logo">
              <div className="logo-icon">💬</div>
              <span className="logo-text">UMess</span>
            </div>
        </div>
        <div className="sidebar-header-right">
          <button className="icon-button" onClick={onNewChat} title="New Chat">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path
                d="M12 5V19M5 12H19"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </button>
          <div className="dropdown-menu">
            <button className="icon-button" onClick={() => setShowMenu(!showMenu)} title="Menu">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path
                  d="M12 12H12.01M12 6H12.01M12 18H12.01"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            </button>
            {showMenu && (
              <div className="dropdown-content">
                <button onClick={() => {}}>New Group</button>
                <button onClick={() => {}}>Settings</button>
                <button onClick={onLogout}>Logout</button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Search Bar */}
      <div className="sidebar-search">
        <div className="search-container">
          <svg className="search-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
            <circle cx="11" cy="11" r="8" stroke="currentColor" strokeWidth="2" />
            <path d="m21 21-4.35-4.35" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
          </svg>
          <input
            type="text"
            className="search-input"
            placeholder="Search or start new chat"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          {searchQuery && (
            <button className="search-clear" onClick={() => setSearchQuery('')}>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="2" />
                <path d="m15 9-6 6m0-6 6 6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
              </svg>
            </button>
          )}
        </div>
      </div>

      {/* Filter Buttons */}
      <div className="sidebar-filters">
        <button
          className={`filter-button ${activeFilter === 'all' ? 'active' : ''}`}
          onClick={() => setActiveFilter('all')}
        >
          All
        </button>
        <button
          className={`filter-button ${activeFilter === 'unread' ? 'active' : ''}`}
          onClick={() => setActiveFilter('unread')}
        >
          Unread
        </button>
        <button
          className={`filter-button ${activeFilter === 'favorites' ? 'active' : ''}`}
          onClick={() => setActiveFilter('favorites')}
        >
          Favorites
        </button>
        <button
          className={`filter-button ${activeFilter === 'groups' ? 'active' : ''}`}
          onClick={() => setActiveFilter('groups')}
        >
          Groups
        </button>
      </div>

      {/* Chat List */}
      <div className="sidebar-chat-list">
        <ChatList
          conversations={filteredConversations}
          selectedConversation={selectedConversation}
          onSelectConversation={onSelectConversation}
          currentUserId={currentUserId}
        />
      </div>

      {/* Bottom Navigation */}
      <div className="sidebar-bottom">
        <button className="bottom-icon" title="Settings">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="3" stroke="currentColor" strokeWidth="2" />
            <path
              d="M12 1v6m0 6v6M5.64 5.64l4.24 4.24m4.24 4.24l4.24 4.24M1 12h6m6 0h6M5.64 18.36l4.24-4.24m4.24-4.24l4.24-4.24"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
            />
          </svg>
        </button>
        <button className="bottom-icon" title="Profile">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
            <path
              d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
            <circle cx="12" cy="7" r="4" stroke="currentColor" strokeWidth="2" />
          </svg>
        </button>
        <button className="bottom-icon status-icon" title="Status">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="2" />
            <circle cx="12" cy="12" r="6" fill="currentColor" />
          </svg>
          <span className="status-badge">1</span>
        </button>
      </div>
    </div>
  )
}

