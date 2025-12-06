import { useState, useEffect } from 'react'
import { Conversation, User, chatService, userService } from '../services/chatService'
import './GroupInfoModal.css'

interface GroupInfoModalProps {
  isOpen: boolean
  onClose: () => void
  conversation: Conversation | null
  currentUserId: number
  onGroupUpdated?: () => void
}

export default function GroupInfoModal({
  isOpen,
  onClose,
  conversation,
  currentUserId,
  onGroupUpdated,
}: GroupInfoModalProps) {
  const [isEditing, setIsEditing] = useState(false)
  const [groupName, setGroupName] = useState('')
  const [description, setDescription] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [addingMember, setAddingMember] = useState(false)
  const [showAddMember, setShowAddMember] = useState(false)

  const isAdmin = conversation?.participants?.find(
    (p) => p.id === currentUserId && p.role === 'admin'
  )

  useEffect(() => {
    if (conversation) {
      setGroupName(conversation.name || '')
      setDescription(conversation.description || '')
    }
  }, [conversation])

  useEffect(() => {
    if (!isOpen) {
      setIsEditing(false)
      setShowAddMember(false)
      setSearchQuery('')
      setUsers([])
    }
  }, [isOpen])

  useEffect(() => {
    if (showAddMember && searchQuery.trim().length >= 1) {
      const timeoutId = setTimeout(() => {
        searchUsers()
      }, 300)
      return () => clearTimeout(timeoutId)
    } else {
      setUsers([])
    }
  }, [searchQuery, showAddMember])

  const searchUsers = async () => {
    const trimmedQuery = searchQuery.trim()
    if (trimmedQuery.length < 1) {
      setUsers([])
      return
    }

    setLoading(true)
    try {
      const results = await userService.searchUsers(trimmedQuery)
      // Filter out current user and existing participants
      const participantIds = conversation?.participants?.map((p) => p.id) || []
      const filtered = results.filter(
        (u) => u.id !== currentUserId && !participantIds.includes(u.id)
      )
      setUsers(filtered)
    } catch (error: any) {
      console.error('Error searching users:', error)
      setUsers([])
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async () => {
    if (!conversation || !groupName.trim()) {
      alert('Please enter a group name')
      return
    }

    setSaving(true)
    try {
      await chatService.updateGroupInfo(conversation.id, groupName.trim(), description.trim())
      setIsEditing(false)
      if (onGroupUpdated) {
        onGroupUpdated()
      }
    } catch (error: any) {
      console.error('Error updating group info:', error)
      alert(error.response?.data?.message || 'Failed to update group info')
    } finally {
      setSaving(false)
    }
  }

  const handleAddMember = async (userId: number) => {
    if (!conversation) return

    setAddingMember(true)
    try {
      await chatService.addGroupMember(conversation.id, userId)
      setSearchQuery('')
      setUsers([])
      if (onGroupUpdated) {
        onGroupUpdated()
      }
    } catch (error: any) {
      console.error('Error adding member:', error)
      alert(error.response?.data?.message || 'Failed to add member')
    } finally {
      setAddingMember(false)
    }
  }

  const handleRemoveMember = async (userId: number) => {
    if (!conversation) return

    if (!window.confirm('Are you sure you want to remove this member?')) {
      return
    }

    try {
      await chatService.removeGroupMember(conversation.id, userId)
      if (onGroupUpdated) {
        onGroupUpdated()
      }
    } catch (error: any) {
      console.error('Error removing member:', error)
      alert(error.response?.data?.message || 'Failed to remove member')
    }
  }

  const handleLeaveGroup = async () => {
    if (!conversation) return

    if (!window.confirm('Are you sure you want to leave this group?')) {
      return
    }

    try {
      await chatService.leaveGroup(conversation.id)
      onClose()
      if (onGroupUpdated) {
        onGroupUpdated()
      }
    } catch (error: any) {
      console.error('Error leaving group:', error)
      alert(error.response?.data?.message || 'Failed to leave group')
    }
  }

  if (!isOpen || !conversation) return null

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content group-info-modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>Group Info</h2>
          <button className="modal-close" onClick={onClose}>×</button>
        </div>
        <div className="modal-body">
          {!isEditing ? (
            <>
              <div className="group-info-section">
                <div className="group-name-display">{conversation.name || 'Unnamed Group'}</div>
                {conversation.description && (
                  <div className="group-description-display">{conversation.description}</div>
                )}
                {isAdmin && (
                  <button className="btn-edit" onClick={() => setIsEditing(true)}>
                    Edit Group Info
                  </button>
                )}
              </div>

              <div className="members-section">
                <div className="section-header">
                  <h3>Members ({conversation.participants?.length || 0})</h3>
                  {isAdmin && (
                    <button
                      className="btn-add-member"
                      onClick={() => setShowAddMember(!showAddMember)}
                    >
                      {showAddMember ? 'Cancel' : 'Add Member'}
                    </button>
                  )}
                </div>

                {showAddMember && isAdmin && (
                  <div className="add-member-section">
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
                    <div className="user-list">
                      {users.map((user) => (
                        <div
                          key={user.id}
                          className="user-item"
                          onClick={() => handleAddMember(user.id)}
                        >
                          <div className="user-avatar">
                            {user.username.charAt(0).toUpperCase()}
                          </div>
                          <div className="user-info">
                            <div className="user-name">{user.username}</div>
                            <div className="user-email">{user.email}</div>
                          </div>
                          {addingMember && (
                            <div className="loading-small">Adding...</div>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                <div className="members-list">
                  {conversation.participants?.map((participant) => (
                    <div key={participant.id} className="member-item">
                      <div className="user-avatar">
                        {participant.username.charAt(0).toUpperCase()}
                      </div>
                      <div className="member-info">
                        <div className="member-name">
                          {participant.username}
                          {participant.id === currentUserId && ' (You)'}
                        </div>
                        <div className="member-role">
                          {participant.role === 'admin' ? 'Admin' : 'Member'}
                        </div>
                      </div>
                      {isAdmin && participant.id !== currentUserId && (
                        <button
                          className="btn-remove-member"
                          onClick={() => handleRemoveMember(participant.id)}
                          title="Remove member"
                        >
                          ×
                        </button>
                      )}
                    </div>
                  ))}
                </div>
              </div>

              <div className="modal-footer">
                <button className="btn-danger" onClick={handleLeaveGroup}>
                  Leave Group
                </button>
              </div>
            </>
          ) : (
            <>
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

              <div className="modal-footer">
                <button
                  className="btn-secondary"
                  onClick={() => {
                    setIsEditing(false)
                    setGroupName(conversation.name || '')
                    setDescription(conversation.description || '')
                  }}
                  disabled={saving}
                >
                  Cancel
                </button>
                <button
                  className="btn-primary"
                  onClick={handleSave}
                  disabled={saving || !groupName.trim()}
                >
                  {saving ? 'Saving...' : 'Save'}
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  )
}

