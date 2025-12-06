import { useState, useEffect } from 'react'
import { userService, User } from '../services/chatService'
import './NewChatModal.css'

interface NewChatModalProps {
  isOpen: boolean
  onClose: () => void
  onSelectUser: (user: User) => void
  currentUserId: number
}

export default function NewChatModal({ isOpen, onClose, onSelectUser, currentUserId }: NewChatModalProps) {
  const [searchQuery, setSearchQuery] = useState('')
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (isOpen && searchQuery.trim().length >= 1) {
      const timeoutId = setTimeout(() => {
        searchUsers()
      }, 300) // Debounce 300ms

      return () => clearTimeout(timeoutId)
    } else {
      setUsers([])
    }
  }, [searchQuery, isOpen])

  const searchUsers = async () => {
    const trimmedQuery = searchQuery.trim()
    if (trimmedQuery.length < 1) {
      setUsers([])
      return
    }

    setLoading(true)
    try {
      const results = await userService.searchUsers(trimmedQuery)
      console.log('Search results:', results) // Debug log
      const filtered = results.filter(u => u.id !== currentUserId)
      console.log('Filtered results:', filtered) // Debug log
      setUsers(filtered)
    } catch (error: any) {
      console.error('Error searching users:', error)
      console.error('Error details:', error.response?.data) // Debug log
      setUsers([])
    } finally {
      setLoading(false)
    }
  }

  const handleSelectUser = (user: User) => {
    onSelectUser(user)
    setSearchQuery('')
    setUsers([])
    onClose()
  }

  if (!isOpen) return null

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>New Chat</h2>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>
        <div className="modal-body">
          <input
            type="text"
            className="search-input"
            placeholder="Search by username or email..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            autoFocus
          />
          {loading && <div className="loading">Searching...</div>}
          {!loading && searchQuery.trim().length >= 1 && users.length === 0 && (
            <div className="no-results">No users found</div>
          )}
          {!loading && searchQuery.trim().length === 0 && (
            <div className="no-results">Type at least 1 character to search</div>
          )}
          <div className="user-list">
            {users.map((user) => (
              <div
                key={user.id}
                className="user-item"
                onClick={() => handleSelectUser(user)}
              >
                <div className="user-avatar">
                  {user.username.charAt(0).toUpperCase()}
                </div>
                <div className="user-info">
                  <div className="user-name">{user.username}</div>
                  <div className="user-email">{user.email}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
