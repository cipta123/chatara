import { useState, useEffect, useRef } from 'react'
import { Conversation, Message } from '../services/chatService'
import { User } from '../services/authService'
import { chatService } from '../services/chatService'
import MessageBubble from './MessageBubble'
import MessageInput from './MessageInput'
import './ChatWindow.css'

interface ChatWindowProps {
  conversation: Conversation
  messages: Message[]
  onSendMessage: (content: string, type?: 'text' | 'image', mediaUrl?: string, replyToId?: number) => void
  currentUser: User
  onBack?: () => void
  isMobile?: boolean
  onMessagesChange?: (messages: Message[]) => void
  onGroupInfoClick?: () => void
}

export default function ChatWindow({
  conversation,
  messages,
  onSendMessage,
  currentUser,
  onBack,
  isMobile = false,
  onMessagesChange,
  onGroupInfoClick,
}: ChatWindowProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const [isTyping, setIsTyping] = useState(false)
  const [replyToMessage, setReplyToMessage] = useState<Message | null>(null)
  
  // Ensure messages is always an array
  const safeMessages = Array.isArray(messages) ? messages : []

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [safeMessages])

  const getOtherParticipant = () => {
    // Ensure participants is an array
    const participants = Array.isArray(conversation.participants) ? conversation.participants : []
    
    // For direct conversations, get the participant that is NOT the current user
    if (conversation.type === 'direct' && participants.length >= 1) {
      // Filter out the current user to get the other participant
      const otherParticipant = participants.find((p) => p.id !== currentUser.id)
      return otherParticipant || null
    }
    
    // For groups, return null (we'll use group name instead)
    if (conversation.type === 'group') {
      return null
    }
    
    // Return first participant or null if none
    return participants[0] || null
  }

  const otherParticipant = getOtherParticipant()
  const isGroup = conversation.type === 'group'
  const displayName = isGroup ? (conversation.name || 'Group') : (otherParticipant?.username || 'Unknown')
  const memberCount = isGroup ? (conversation.participants?.length || 0) : null

  const handleCopyMessage = (message: Message) => {
    const textToCopy = message.content || (message.type === 'image' ? 'Image' : '')
    if (textToCopy) {
      navigator.clipboard.writeText(textToCopy).catch((err) => {
        console.error('Failed to copy:', err)
        // Fallback for older browsers
        const textArea = document.createElement('textarea')
        textArea.value = textToCopy
        document.body.appendChild(textArea)
        textArea.select()
        document.execCommand('copy')
        document.body.removeChild(textArea)
      })
    }
  }

  const handleReplyToMessage = (message: Message) => {
    setReplyToMessage(message)
  }

  const handleDeleteMessage = async (message: Message) => {
    if (!confirm('Are you sure you want to delete this message?')) {
      return
    }

    try {
      await chatService.deleteMessage(conversation.id, message.id)
      // Remove message from local state
      if (onMessagesChange) {
        const updatedMessages = safeMessages.filter((m) => m.id !== message.id)
        onMessagesChange(updatedMessages)
      }
    } catch (error) {
      console.error('Failed to delete message:', error)
      alert('Failed to delete message. Please try again.')
    }
  }

  const handleSendWithReply = (content: string, type: 'text' | 'image' = 'text', mediaUrl?: string) => {
    const replyToId = replyToMessage?.id
    onSendMessage(content, type, mediaUrl, replyToId)
    setReplyToMessage(null) // Clear reply after sending
  }

  const handleCancelReply = () => {
    setReplyToMessage(null)
  }
  
  // Don't render if no other participant (for direct chats) and not a group
  if (!isGroup && !otherParticipant) {
    return (
      <div className="chat-window">
        <div className="chat-window-header">
          <div className="chat-window-header-info">
            <div className="chat-window-header-name">Loading...</div>
          </div>
        </div>
        <div className="chat-window-messages">
          <div className="chat-empty">Loading conversation...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="chat-window">
      <div className="chat-window-header">
        {isMobile && onBack && (
          <button className="chat-back-button" onClick={onBack} aria-label="Back to chat list">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <path
                d="M19 12H5M5 12L12 19M5 12L12 5"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </button>
        )}
        <div className="chat-window-user" onClick={isGroup && onGroupInfoClick ? onGroupInfoClick : undefined} style={{ cursor: isGroup && onGroupInfoClick ? 'pointer' : 'default' }}>
          <div className="chat-window-avatar">
            {isGroup ? '👥' : (otherParticipant?.username?.charAt(0).toUpperCase() || 'U')}
          </div>
          <div>
            <div className="chat-window-name">{displayName}</div>
            <div className="chat-window-status">
              {isGroup ? `${memberCount} ${memberCount === 1 ? 'member' : 'members'}` : 'Online'}
            </div>
          </div>
        </div>
        {isGroup && onGroupInfoClick && (
          <button className="chat-info-button" onClick={onGroupInfoClick} title="Group Info">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="2" />
              <path d="M12 16v-4M12 8h.01" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
            </svg>
          </button>
        )}
      </div>

      <div className="chat-window-messages">
        {safeMessages.length === 0 ? (
          <div className="chat-window-empty">
            No messages yet. Start the conversation!
          </div>
        ) : (
          safeMessages.map((message) => (
            <MessageBubble
              key={message.id}
              message={message}
              isOwn={message.sender_id === currentUser.id}
              onCopy={() => handleCopyMessage(message)}
              onReply={() => handleReplyToMessage(message)}
              onDelete={() => handleDeleteMessage(message)}
            />
          ))
        )}
        <div ref={messagesEndRef} />
      </div>

      {isTyping && (
        <div className="typing-indicator">
          <span>{otherParticipant?.username} is typing...</span>
        </div>
      )}

      <MessageInput 
        onSendMessage={handleSendWithReply} 
        conversationId={conversation.id}
        replyTo={replyToMessage}
        onCancelReply={handleCancelReply}
      />
    </div>
  )
}


