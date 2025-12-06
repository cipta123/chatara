import { useState, useRef, useEffect } from 'react'
import { Message } from '../services/chatService'
import MessageMenu from './MessageMenu'
import './MessageBubble.css'

interface MessageBubbleProps {
  message: Message
  isOwn: boolean
  onCopy?: () => void
  onReply?: () => void
  onDelete?: () => void
}

export default function MessageBubble({ message, isOwn, onCopy, onReply, onDelete }: MessageBubbleProps) {
  const [showMenu, setShowMenu] = useState(false)
  const [menuPosition, setMenuPosition] = useState({ top: 0, left: 0 })
  const messageRef = useRef<HTMLDivElement>(null)
  const buttonRef = useRef<HTMLButtonElement>(null)

  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
  }

  // Resolve media URL - if it starts with /api, prepend the base URL
  const getImageUrl = (mediaUrl: string | undefined): string => {
    if (!mediaUrl) return ''
    if (mediaUrl.startsWith('http://') || mediaUrl.startsWith('https://')) {
      return mediaUrl
    }
    // If it's a relative path starting with /api, prepend the API base URL
    if (mediaUrl.startsWith('/api/')) {
      const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'
      // Remove /api from the base if present to avoid double /api
      const base = apiBase.replace(/\/api$/, '')
      return `${base}${mediaUrl}`
    }
    return mediaUrl
  }

  const imageUrl = message.type === 'image' && message.media_url ? getImageUrl(message.media_url) : null

  // Handle click outside to close menu
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (messageRef.current && !messageRef.current.contains(event.target as Node)) {
        setShowMenu(false)
      }
    }

    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setShowMenu(false)
      }
    }

    if (showMenu) {
      document.addEventListener('mousedown', handleClickOutside)
      document.addEventListener('keydown', handleEscape)
      return () => {
        document.removeEventListener('mousedown', handleClickOutside)
        document.removeEventListener('keydown', handleEscape)
      }
    }
  }, [showMenu])

  const handleMenuToggle = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation()
    if (!showMenu && buttonRef.current && messageRef.current) {
      const buttonRect = buttonRef.current.getBoundingClientRect()
      const messageRect = messageRef.current.getBoundingClientRect()
      
      // Position menu close to the message bubble
      const menuWidth = 180
      const menuItemHeight = 44
      const menuHeight = isOwn ? (menuItemHeight * 3) : (menuItemHeight * 2)
      
      // Position menu near the top of the message, close to the button
      let top = messageRect.top + 4
      let left: number
      
      if (isOwn) {
        // Own message: menu appears to the left, very close
        left = messageRect.left - menuWidth - 4
      } else {
        // Other message: menu appears to the right, very close
        left = messageRect.right + 4
      }
      
      // Adjust if menu would go off screen vertically
      if (top + menuHeight > window.innerHeight - 8) {
        top = Math.max(8, window.innerHeight - menuHeight - 8)
      }
      if (top < 8) {
        top = 8
      }
      
      // Adjust if menu would go off screen horizontally - flip to other side
      if (left < 8) {
        // If off left edge, show on right side
        left = messageRect.right + 4
      }
      if (left + menuWidth > window.innerWidth - 8) {
        // If off right edge, show on left side
        left = messageRect.left - menuWidth - 4
        if (left < 8) {
          // If still off, align to edge
          left = 8
        }
      }
      
      setMenuPosition({ top, left })
    }
    setShowMenu(!showMenu)
  }

  const handleCopy = () => {
    if (onCopy) {
      onCopy()
    }
  }

  const handleReply = () => {
    if (onReply) {
      onReply()
    }
  }

  const handleDelete = () => {
    if (onDelete) {
      onDelete()
    }
  }

  return (
    <div className={`message-bubble-wrapper ${isOwn ? 'own' : 'other'}`} ref={messageRef}>
      <div className={`message-bubble ${isOwn ? 'own' : 'other'}`}>
      {message.reply_to && (
        <div className="message-reply">
          <div className="message-reply-line" />
          <div className="message-reply-content">
            <div className="message-reply-name">{message.reply_to.sender?.username || 'Unknown'}</div>
            <div className="message-reply-text">
              {message.reply_to.type === 'image' ? 'Image' : message.reply_to.content || ''}
            </div>
          </div>
        </div>
      )}
      {imageUrl && (
        <div className="message-image">
          <img 
            src={imageUrl} 
            alt={message.content || 'Shared image'} 
            loading="lazy"
            onError={(e) => {
              console.error('Failed to load image:', imageUrl)
              // Optionally show error state
              const target = e.target as HTMLImageElement
              target.style.display = 'none'
            }}
          />
        </div>
      )}
      {message.content && (
        <div className="message-content">
          <span className="message-text">{message.content}</span>
          <div className="message-footer">
            <span className="message-time">{formatTime(message.created_at)}</span>
            {isOwn && (
              <span className={`message-status ${message.status}`}>
                {message.status === 'sent' && '✓'}
                {message.status === 'delivered' && '✓✓'}
                {message.status === 'read' && '✓✓'}
              </span>
            )}
            <button
              ref={buttonRef}
              className="message-dropdown-btn"
              onClick={handleMenuToggle}
              aria-label="Message options"
              type="button"
            >
              <svg width="10" height="10" viewBox="0 0 10 10" fill="none">
                <path
                  d="M5 7.5L2 2.5h6L5 7.5z"
                  fill="currentColor"
                />
              </svg>
            </button>
          </div>
        </div>
      )}
      {!message.content && imageUrl && (
        <div className="message-footer message-footer-standalone">
          <span className="message-time">{formatTime(message.created_at)}</span>
          {isOwn && (
            <span className={`message-status ${message.status}`}>
              {message.status === 'sent' && '✓'}
              {message.status === 'delivered' && '✓✓'}
              {message.status === 'read' && '✓✓'}
            </span>
          )}
          <button
            ref={buttonRef}
            className="message-dropdown-btn"
            onClick={handleMenuToggle}
            aria-label="Message options"
            type="button"
          >
            <svg width="10" height="10" viewBox="0 0 10 10" fill="none">
              <path
                d="M5 7.5L2 2.5h6L5 7.5z"
                fill="currentColor"
              />
            </svg>
          </button>
        </div>
      )}
      </div>
      <MessageMenu
        isOpen={showMenu}
        onClose={() => setShowMenu(false)}
        onCopy={handleCopy}
        onReply={handleReply}
        onDelete={handleDelete}
        position={menuPosition}
        isOwn={isOwn}
      />
    </div>
  )
}


