import { useState, useEffect, useRef } from 'react'
import { useAuth } from '../hooks/useAuth'
import { useWebSocket } from '../hooks/useWebSocket'
import { chatService, Message, Conversation, User } from '../services/chatService'
import ChatList from '../components/ChatList'
import ChatWindow from '../components/ChatWindow'
import NewChatModal from '../components/NewChatModal'
import './Chat.css'

// WebSocket URL - must be full URL for direct connection
const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'

export default function Chat() {
  const { user, logout } = useAuth()
  const [conversations, setConversations] = useState<Conversation[]>([])
  const [selectedConversation, setSelectedConversation] = useState<Conversation | null>(null)
  const [messages, setMessages] = useState<Message[]>([])
  const [loading, setLoading] = useState(true)
  const [showNewChatModal, setShowNewChatModal] = useState(false)
  const token = localStorage.getItem('token')

  // Load conversations on mount
  useEffect(() => {
    loadConversations()
  }, [])

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

  const handleSendMessage = async (content: string, type: 'text' | 'image' = 'text', mediaUrl?: string) => {
    if (!selectedConversation) return

    try {
      const newMessage = await chatService.sendMessage({
        conversation_id: selectedConversation.id,
        content,
        type,
        media_url: mediaUrl,
      })

      setMessages((prev) => [...prev, newMessage])
    } catch (error) {
      console.error('Error sending message:', error)
    }
  }

  const handleSelectConversation = (conversation: Conversation) => {
    setSelectedConversation(conversation)
    loadMessages(conversation.id)
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

  if (loading) {
    return <div className="chat-loading">Loading...</div>
  }

  return (
    <div className="chat-container">
      <div className="chat-sidebar">
        <div className="chat-header">
          <div className="user-info">
            <div className="user-avatar">
              {user?.username?.charAt(0).toUpperCase() || 'U'}
            </div>
            <div>
              <div className="user-name">{user?.username}</div>
              <div className="user-status">{isConnected ? 'Online' : 'Offline'}</div>
            </div>
          </div>
          <div className="header-actions">
            <button onClick={() => setShowNewChatModal(true)} className="new-chat-button" title="New Chat">
              +
            </button>
            <button onClick={logout} className="logout-button">
              Logout
            </button>
          </div>
        </div>
        <ChatList
          conversations={conversations}
          selectedConversation={selectedConversation}
          onSelectConversation={handleSelectConversation}
          currentUserId={user?.id || 0}
        />
      </div>
      <div className="chat-main">
        {selectedConversation ? (
          <ChatWindow
            conversation={selectedConversation}
            messages={messages}
            onSendMessage={handleSendMessage}
            currentUser={user!}
          />
        ) : (
          <div className="chat-empty">Select a conversation to start chatting</div>
        )}
      </div>
      <NewChatModal
        isOpen={showNewChatModal}
        onClose={() => setShowNewChatModal(false)}
        onSelectUser={handleStartChat}
        currentUserId={user?.id || 0}
      />
    </div>
  )
}


