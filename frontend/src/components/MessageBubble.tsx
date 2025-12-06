import { Message } from '../services/chatService'
import './MessageBubble.css'

interface MessageBubbleProps {
  message: Message
  isOwn: boolean
}

export default function MessageBubble({ message, isOwn }: MessageBubbleProps) {
  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
  }

  return (
    <div className={`message-bubble ${isOwn ? 'own' : 'other'}`}>
      {message.type === 'image' && message.media_url ? (
        <div className="message-image">
          <img src={message.media_url} alt="Shared image" />
        </div>
      ) : null}
      {message.content && (
        <div className="message-content">{message.content}</div>
      )}
      <div className="message-footer">
        <span className="message-time">{formatTime(message.created_at)}</span>
        {isOwn && (
          <span className={`message-status ${message.status}`}>
            {message.status === 'sent' && '✓'}
            {message.status === 'delivered' && '✓✓'}
            {message.status === 'read' && '✓✓'}
          </span>
        )}
      </div>
    </div>
  )
}


