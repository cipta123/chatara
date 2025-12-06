import { useState, useEffect, useRef } from 'react'
import { useAuth } from '../hooks/useAuth'
import { useWebSocket } from '../hooks/useWebSocket'
import { chatService, Message, Conversation, User } from '../services/chatService'
import LeftSidebar from '../components/LeftSidebar'
import Sidebar from '../components/Sidebar'
import ChatWindow from '../components/ChatWindow'
import NewChatModal from '../components/NewChatModal'
import CreateGroupModal from '../components/CreateGroupModal'
import GroupInfoModal from '../components/GroupInfoModal'
import './Chat.css'

// WebSocket URL - must be full URL for direct connection
// Auto-detect WebSocket URL based on current hostname
const getWebSocketUrl = () => {
  if (import.meta.env.VITE_WS_URL) {
    return import.meta.env.VITE_WS_URL
  }
  // Use current hostname and port 8080 for WebSocket (bypass proxy)
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const hostname = window.location.hostname
  return `${protocol}//${hostname}:8080/ws`
}
const WS_URL = getWebSocketUrl()

export default function Chat() {
  const { user, logout } = useAuth()
  const [conversations, setConversations] = useState<Conversation[]>([])
  const [selectedConversation, setSelectedConversation] = useState<Conversation | null>(null)
  const [messages, setMessages] = useState<Message[]>([])
  const [loading, setLoading] = useState(true)
  const [showNewChatModal, setShowNewChatModal] = useState(false)
  const [showCreateGroupModal, setShowCreateGroupModal] = useState(false)
  const [showGroupInfoModal, setShowGroupInfoModal] = useState(false)
  const [canCreateGroup, setCanCreateGroup] = useState(false)
  const [isMobile, setIsMobile] = useState(window.innerWidth < 768)
  const [showConversation, setShowConversation] = useState(false)
  const token = localStorage.getItem('token')

  // Detect mobile view
  useEffect(() => {
    const handleResize = () => {
      const mobile = window.innerWidth < 768
      setIsMobile(mobile)
      // On mobile, if screen gets bigger, show conversation if one is selected
      if (!mobile && selectedConversation) {
        setShowConversation(false)
      }
    }

    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [selectedConversation])

  // Load conversations on mount
  useEffect(() => {
    loadConversations()
    checkPermission()
  }, [])

  const checkPermission = async () => {
    // Check if current user has permission to create groups
    // This is a simple check - in production, you'd fetch this from the backend
    if (user?.email === 'ciptaanugrahh@gmail.com') {
      setCanCreateGroup(true)
    } else {
      // Try to get permission from user object if available
      setCanCreateGroup(user?.can_create_group || false)
    }
  }

  // WebSocket connection
  const { isConnected, sendMessage } = useWebSocket(
    WS_URL,
    token,
    (wsMessage) => {
      if (wsMessage.type === 'new_message' && wsMessage.message) {
        const newMessage: Message = wsMessage.message
        setMessages((prev) => {
          // Avoid duplicates
          if (prev.some((m) => m.id === newMessage.id)) {
            return prev
          }
          return [...prev, newMessage]
        })

        // Update conversation list with new message
        setConversations((prev) =>
          prev.map((conv) =>
            conv.id === newMessage.conversation_id
              ? { ...conv, last_message: newMessage, updated_at: newMessage.created_at }
              : conv
          )
        )
      }
    }
  )

  // Load messages when conversation is selected
  useEffect(() => {
    if (selectedConversation) {
      loadMessages(selectedConversation.id)
    }
  }, [selectedConversation])

  const loadConversations = async () => {
    try {
      const data = await chatService.getConversations()
      // Ensure data is an array
      const conversations = Array.isArray(data) ? data : []
      setConversations(conversations)
      if (conversations.length > 0 && !selectedConversation) {
        setSelectedConversation(conversations[0])
      }
    } catch (error) {
      console.error('Error loading conversations:', error)
      setConversations([]) // Set empty array on error
    } finally {
      setLoading(false)
    }
  }

  const loadMessages = async (conversationId: number) => {
    try {
      const data = await chatService.getMessages(conversationId)
      // Ensure data is always an array
      const messages = Array.isArray(data) ? data : []
      setMessages(messages)
      // Mark as read
      await chatService.markAsRead(conversationId)
    } catch (error) {
      console.error('Error loading messages:', error)
      setMessages([]) // Set empty array on error
    }
  }

  const handleSendMessage = async (content: string, type: 'text' | 'image' = 'text', mediaUrl?: string, replyToId?: number) => {
    if (!selectedConversation) return

    try {
      const newMessage = await chatService.sendMessage({
        conversation_id: selectedConversation.id,
        content,
        type,
        media_url: mediaUrl,
        reply_to_id: replyToId,
      })

      setMessages((prev) => [...prev, newMessage])
    } catch (error) {
      console.error('Error sending message:', error)
    }
  }

  const handleMessagesChange = (updatedMessages: Message[]) => {
    setMessages(updatedMessages)
  }

  const handleSelectConversation = (conversation: Conversation) => {
    setSelectedConversation(conversation)
    loadMessages(conversation.id)
    // On mobile, show conversation screen
    if (isMobile) {
      setShowConversation(true)
    }
  }

  const handleBackToList = () => {
    setShowConversation(false)
  }

  const handleStartChat = async (selectedUser: User) => {
    try {
      console.log('Creating conversation with user:', selectedUser)
      const conversation = await chatService.createDirectConversation(selectedUser.id)
      console.log('Conversation created:', conversation)
      
      // Reload conversations to get updated list
      await loadConversations()
      
      // Set the newly created conversation as selected
      setSelectedConversation(conversation)
      
      // Load messages for this conversation
      loadMessages(conversation.id)
    } catch (error) {
      console.error('Error starting chat:', error)
      alert('Failed to start chat. Please try again.')
    }
  }

  const handleCreateGroup = async (conversation: Conversation) => {
    // Reload conversations to get updated list
    await loadConversations()
    
    // Set the newly created group as selected
    setSelectedConversation(conversation)
    
    // Load messages for this conversation
    loadMessages(conversation.id)
  }

  const handleGroupUpdated = async () => {
    // Reload conversations and current conversation
    await loadConversations()
    if (selectedConversation) {
      const updated = await chatService.getConversation(selectedConversation.id)
      setSelectedConversation(updated)
    }
  }

  const [activeNavItem, setActiveNavItem] = useState<'chats' | 'status' | 'communities' | 'calls' | 'settings'>('chats')

  if (loading) {
    return <div className="chat-loading">Loading...</div>
  }

  return (
    <div className="chat-container">
      {/* Left sidebar - desktop only, bottom nav on mobile (handled by CSS) */}
      <LeftSidebar 
        activeItem={activeNavItem} 
        onItemClick={setActiveNavItem}
        hideOnMobile={isMobile && showConversation}
      />
      <div className="chat-content-wrapper">
        {activeNavItem === 'chats' && (
          <>
            {/* On mobile: show either chat list OR conversation, not both */}
            {isMobile ? (
              <>
                {!showConversation && (
                  <Sidebar
                    conversations={conversations}
                    selectedConversation={selectedConversation}
                    onSelectConversation={handleSelectConversation}
                    currentUserId={user?.id || 0}
                    currentUsername={user?.username || ''}
                    isConnected={isConnected}
                    onNewChat={() => setShowNewChatModal(true)}
                    onNewGroup={canCreateGroup ? () => setShowCreateGroupModal(true) : undefined}
                    canCreateGroup={canCreateGroup}
                    onLogout={logout}
                  />
                )}
                {showConversation && selectedConversation && (
                  <div className="chat-main mobile-fullscreen">
                    <ChatWindow
                        conversation={selectedConversation}
                        messages={messages}
                        onSendMessage={handleSendMessage}
                        currentUser={user!}
                        onBack={handleBackToList}
                        isMobile={isMobile}
                        onMessagesChange={handleMessagesChange}
                        onGroupInfoClick={selectedConversation.type === 'group' ? () => setShowGroupInfoModal(true) : undefined}
                      />
                  </div>
                )}
              </>
            ) : (
              <>
                <Sidebar
                  conversations={conversations}
                  selectedConversation={selectedConversation}
                  onSelectConversation={handleSelectConversation}
                  currentUserId={user?.id || 0}
                  currentUsername={user?.username || ''}
                  isConnected={isConnected}
                  onNewChat={() => setShowNewChatModal(true)}
                  onNewGroup={canCreateGroup ? () => setShowCreateGroupModal(true) : undefined}
                  canCreateGroup={canCreateGroup}
                  onLogout={logout}
                />
                <div className="chat-main">
                  {selectedConversation ? (
                    <ChatWindow
                        conversation={selectedConversation}
                        messages={messages}
                        onSendMessage={handleSendMessage}
                        currentUser={user!}
                        isMobile={isMobile}
                        onMessagesChange={handleMessagesChange}
                        onGroupInfoClick={selectedConversation.type === 'group' ? () => setShowGroupInfoModal(true) : undefined}
                      />
                  ) : (
                    <div className="chat-empty">
                      <div className="empty-illustration">
                        <div className="empty-icon">💬</div>
                        <h2>UMess</h2>
                        <p>Select a conversation to start chatting</p>
                        <p className="empty-subtitle">Your personal messages are end-to-end encrypted</p>
                      </div>
                    </div>
                  )}
                </div>
              </>
            )}
          </>
        )}
        {activeNavItem === 'status' && (
          <div className="chat-main full-width">
            <div className="chat-empty">
              <div className="empty-illustration">
                <div className="empty-icon">⚪</div>
                <h2>Status</h2>
                <p>Status updates coming soon</p>
              </div>
            </div>
          </div>
        )}
        {activeNavItem === 'communities' && (
          <div className="chat-main full-width">
            <div className="chat-empty">
              <div className="empty-illustration">
                <div className="empty-icon">👥</div>
                <h2>Communities</h2>
                <p>Communities feature coming soon</p>
              </div>
            </div>
          </div>
        )}
        {activeNavItem === 'calls' && (
          <div className="chat-main full-width">
            <div className="chat-empty">
              <div className="empty-illustration">
                <div className="empty-icon">📞</div>
                <h2>Calls</h2>
                <p>Voice and video calls coming soon</p>
              </div>
            </div>
          </div>
        )}
        {activeNavItem === 'settings' && (
          <div className="chat-main full-width">
            <div className="chat-empty">
              <div className="empty-illustration">
                <div className="empty-icon">⚙️</div>
                <h2>Settings</h2>
                <p>Settings page coming soon</p>
              </div>
            </div>
          </div>
        )}
      </div>
      <NewChatModal
        isOpen={showNewChatModal}
        onClose={() => setShowNewChatModal(false)}
        onSelectUser={handleStartChat}
        currentUserId={user?.id || 0}
      />
      {canCreateGroup && (
        <CreateGroupModal
          isOpen={showCreateGroupModal}
          onClose={() => setShowCreateGroupModal(false)}
          onCreateGroup={handleCreateGroup}
          currentUserId={user?.id || 0}
        />
      )}
      <GroupInfoModal
        isOpen={showGroupInfoModal}
        onClose={() => setShowGroupInfoModal(false)}
        conversation={selectedConversation}
        currentUserId={user?.id || 0}
        onGroupUpdated={handleGroupUpdated}
      />
    </div>
  )
}


