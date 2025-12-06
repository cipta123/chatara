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
  onSendMessage: (content: string, type?: 'text' | 'image', mediaUrl?: string) => void
  currentUser: User
}

export default function ChatWindow({
  conversation,
  messages,
  onSendMessage,
  currentUser,
}: ChatWindowProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const [isTyping, setIsTyping] = useState(false)
  
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
    
    // Return first participant or null if none
    return participants[0] || null
  }

  const otherParticipant = getOtherParticipant()
  
  // Don't render if no other participant
  if (!otherParticipant) {
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
        <div className="chat-window-user">
          <div className="chat-window-avatar">
            {otherParticipant?.username?.charAt(0).toUpperCase() || 'U'}
          </div>
          <div>
            <div className="chat-window-name">{otherParticipant?.username || 'Unknown'}</div>
            <div className="chat-window-status">Online</div>
          </div>
        </div>
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

      <MessageInput onSendMessage={onSendMessage} conversationId={conversation.id} />
    </div>
  )
}


