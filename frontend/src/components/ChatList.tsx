import { Conversation } from '../services/chatService'
import './ChatList.css'

interface ChatListProps {
  conversations: Conversation[]
  selectedConversation: Conversation | null
  onSelectConversation: (conversation: Conversation) => void
  currentUserId: number
}

export default function ChatList({
  conversations,
  selectedConversation,
  onSelectConversation,
  currentUserId,
}: ChatListProps) {
  // Ensure conversations is always an array
  const safeConversations = Array.isArray(conversations) ? conversations : []

  const getOtherParticipant = (conversation: Conversation) => {
    // Ensure participants is an array
    const participants = Array.isArray(conversation.participants) ? conversation.participants : []
    
    // For direct conversations, get the participant that is NOT the current user
    if (conversation.type === 'direct' && participants.length >= 1) {
      // Filter out the current user to get the other participant
      const otherParticipant = participants.find((p) => p.id !== currentUserId)
      return otherParticipant || null
    }
    
    // Return first participant or null if none
    return participants[0] || null
  }

  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))

    if (days === 0) {
      return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
    } else if (days === 1) {
      return 'Yesterday'
    } else if (days < 7) {
      return date.toLocaleDateString('en-US', { weekday: 'short' })
    } else {
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
    }
  }

  return (
    <div className="chat-list">
      {safeConversations.length === 0 ? (
        <div className="chat-list-empty">No conversations yet</div>
      ) : (
        safeConversations
          .filter((conversation) => {
            // Only show conversations with participants
            const otherParticipant = getOtherParticipant(conversation)
            return otherParticipant !== null
          })
          .map((conversation) => {
            const otherParticipant = getOtherParticipant(conversation)
            const isSelected = selectedConversation?.id === conversation.id

            return (
              <div
                key={conversation.id}
                className={`chat-item ${isSelected ? 'selected' : ''}`}
                onClick={() => onSelectConversation(conversation)}
              >
              <div className="chat-item-avatar">
                {otherParticipant?.username?.charAt(0).toUpperCase() || 'U'}
              </div>
              <div className="chat-item-content">
                <div className="chat-item-header">
                  <div className="chat-item-name">{otherParticipant?.username || 'Unknown'}</div>
                  {conversation.last_message && (
                    <div className="chat-item-time">
                      {formatTime(conversation.last_message.created_at)}
                    </div>
                  )}
                </div>
                <div className="chat-item-preview">
                  {conversation.last_message ? (
                    <>
                      {conversation.last_message.type === 'image' ? (
                        <span>📷 Image</span>
                      ) : (
                        <span>{conversation.last_message.content}</span>
                      )}
                    </>
                  ) : (
                    <span className="chat-item-empty">Start a conversation</span>
                  )}
                </div>
              </div>
            </div>
          )
        })
      )}
    </div>
  )
}


