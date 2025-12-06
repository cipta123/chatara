import { useState, useEffect } from 'react'
import { userService, User, chatService } from '../services/chatService'
import './CreateGroupModal.css'

interface CreateGroupModalProps {
  isOpen: boolean
  onClose: () => void
  onCreateGroup: (conversation: any) => void
  currentUserId: number
}

export default function CreateGroupModal({
  isOpen,
  onClose,
  onCreateGroup,
  currentUserId,
}: CreateGroupModalProps) {
  const [groupName, setGroupName] = useState('')
  const [description, setDescription] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const [users, setUsers] = useState<User[]>([])
  const [selectedUsers, setSelectedUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    if (isOpen && searchQuery.trim().length >= 1) {
      const timeoutId = setTimeout(() => {
        searchUsers()
      }, 300)

      return () => clearTimeout(timeoutId)
    } else {
      setUsers([])
    }
  }, [searchQuery, isOpen])

  useEffect(() => {
    if (!isOpen) {
      // Reset form when modal closes
      setGroupName('')
      setDescription('')
      setSearchQuery('')
      setUsers([])
      setSelectedUsers([])
    }
  }, [isOpen])

  const searchUsers = async () => {
    const trimmedQuery = searchQuery.trim()
    if (trimmedQuery.length < 1) {
      setUsers([])
      return
    }

    setLoading(true)
    try {
      const results = await userService.searchUsers(trimmedQuery)
      // Filter out current user and already selected users
      const filtered = results.filter(
        (u) => u.id !== currentUserId && !selectedUsers.find((su) => su.id === u.id)
      )
      setUsers(filtered)
    } catch (error: any) {
      console.error('Error searching users:', error)
      setUsers([])
    } finally {
      setLoading(false)
    }
  }

  const handleSelectUser = (user: User) => {
    if (!selectedUsers.find((u) => u.id === user.id)) {
      setSelectedUsers([...selectedUsers, user])
      setSearchQuery('')
      setUsers([])
    }
  }

  const handleRemoveUser = (userId: number) => {
    setSelectedUsers(selectedUsers.filter((u) => u.id !== userId))
  }

  const handleCreateGroup = async () => {
    if (!groupName.trim()) {
      alert('Please enter a group name')
      return
    }

    if (selectedUsers.length < 1) {
      alert('Please select at least one member')
      return
    }

    setCreating(true)
    try {
      const participantIds = selectedUsers.map((u) => u.id)
      const conversation = await chatService.createGroup(
        groupName.trim(),
        description.trim(),
        participantIds
      )
      onCreateGroup(conversation)
      onClose()
    } catch (error: any) {
      console.error('Error creating group:', error)
      alert(error.response?.data?.message || 'Failed to create group')
    } finally {
      setCreating(false)
    }
  }

  if (!isOpen) return null

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content create-group-modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>Create Group</h2>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>
        <div className="modal-body">
          <div className="form-group">
            <label>Group Name *</label>
            <input
              type="text"
              className="form-input"
              placeholder="Enter group name..."
              value={groupName}
              onChange={(e) => setGroupName(e.target.value)}
              maxLength={255}
            />
          </div>

          <div className="form-group">
            <label>Description (Optional)</label>
            <textarea
              className="form-textarea"
              placeholder="Enter group description..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
            />
          </div>

          <div className="form-group">
            <label>Add Members *</label>
            <input
              type="text"
              className="search-input"
              placeholder="Search by username or email..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
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

          {selectedUsers.length > 0 && (
            <div className="selected-users">
              <label>Selected Members ({selectedUsers.length})</label>
              <div className="selected-users-list">
                {selectedUsers.map((user) => (
                  <div key={user.id} className="selected-user-item">
                    <div className="user-avatar small">
                      {user.username.charAt(0).toUpperCase()}
                    </div>
                    <span className="selected-user-name">{user.username}</span>
                    <button
                      className="remove-user-btn"
                      onClick={() => handleRemoveUser(user.id)}
                      type="button"
                    >
                      ×
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}

          <div className="modal-footer">
            <button className="btn-secondary" onClick={onClose} disabled={creating}>
              Cancel
            </button>
            <button
              className="btn-primary"
              onClick={handleCreateGroup}
              disabled={creating || !groupName.trim() || selectedUsers.length < 1}
            >
              {creating ? 'Creating...' : 'Create Group'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

