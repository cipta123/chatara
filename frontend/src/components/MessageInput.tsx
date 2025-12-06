import { useState, useRef } from 'react'
import { chatService, Message } from '../services/chatService'
import './MessageInput.css'

interface MessageInputProps {
  onSendMessage: (content: string, type?: 'text' | 'image', mediaUrl?: string, replyToId?: number) => void
  conversationId: number
  replyTo?: Message | null
  onCancelReply?: () => void
}

export default function MessageInput({ onSendMessage, conversationId, replyTo, onCancelReply }: MessageInputProps) {
  const [message, setMessage] = useState('')
  const [uploading, setUploading] = useState(false)
  const [imagePreview, setImagePreview] = useState<string | null>(null)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    // If there's an image preview, send the image first
    if (imagePreview && selectedFile) {
      await handleSendImage()
      return
    }
    
    if (!message.trim()) return

    const content = message.trim()
    const replyToId = replyTo?.id
    setMessage('')
    onSendMessage(content, 'text', undefined, replyToId)
    if (onCancelReply) {
      onCancelReply()
    }
  }

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    // Validate file type
    if (!file.type.startsWith('image/')) {
      alert('Please select an image file')
      return
    }

    // Validate file size (max 10MB)
    if (file.size > 10 * 1024 * 1024) {
      alert('Image size should be less than 10MB')
      return
    }

    // Create preview
    const reader = new FileReader()
    reader.onloadend = () => {
      setImagePreview(reader.result as string)
      setSelectedFile(file)
    }
    reader.readAsDataURL(file)
  }

  const handleSendImage = async () => {
    if (!selectedFile) return

    setUploading(true)

    try {
      const result = await chatService.uploadImage(selectedFile)
      const caption = message.trim() || undefined
      const replyToId = replyTo?.id
      onSendMessage(caption || '', 'image', result.url, replyToId)
      
      // Clear preview and file
      setImagePreview(null)
      setSelectedFile(null)
      setMessage('')
      if (onCancelReply) {
        onCancelReply()
      }
    } catch (error) {
      console.error('Error uploading image:', error)
      alert('Failed to upload image. Please try again.')
    } finally {
      setUploading(false)
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    }
  }

  const handleCancelPreview = () => {
    setImagePreview(null)
    setSelectedFile(null)
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  return (
    <div className="message-input-container">
      {replyTo && (
        <div className="reply-preview">
          <div className="reply-preview-content">
            <div className="reply-preview-line" />
            <div className="reply-preview-info">
              <div className="reply-preview-name">{replyTo.sender?.username || 'Unknown'}</div>
              <div className="reply-preview-text">
                {replyTo.type === 'image' ? 'Image' : replyTo.content || ''}
              </div>
            </div>
          </div>
          {onCancelReply && (
            <button
              type="button"
              className="reply-cancel-button"
              onClick={onCancelReply}
              aria-label="Cancel reply"
            >
              ×
            </button>
          )}
        </div>
      )}
      {imagePreview && (
        <div className="image-preview">
          <div className="image-preview-wrapper">
            <img src={imagePreview} alt="Preview" />
            <button
              type="button"
              className="cancel-preview-button"
              onClick={handleCancelPreview}
              aria-label="Cancel image"
            >
              ×
            </button>
          </div>
          {message && (
            <div className="image-caption">
              <input
                type="text"
                className="caption-input"
                placeholder="Add a caption (optional)..."
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                disabled={uploading}
              />
            </div>
          )}
        </div>
      )}
      <form onSubmit={handleSubmit} className="message-input-form">
        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          onChange={handleFileSelect}
          style={{ display: 'none' }}
        />
        <div className="message-input-wrapper">
          <button
            type="button"
            className="attach-button"
            onClick={() => fileInputRef.current?.click()}
            disabled={uploading}
            title="Attach image"
          >
            {uploading ? (
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="10" strokeDasharray="31.416" strokeDashoffset="31.416">
                  <animate attributeName="stroke-dasharray" dur="2s" values="0 31.416;15.708 15.708;0 31.416;0 31.416" repeatCount="indefinite"/>
                  <animate attributeName="stroke-dashoffset" dur="2s" values="0;-31.416;-15.708;-31.416;0" repeatCount="indefinite"/>
                </circle>
              </svg>
            ) : (
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"/>
              </svg>
            )}
          </button>
          <input
            type="text"
            className="message-input"
            placeholder={imagePreview ? "Add a caption (optional)..." : "Type a message..."}
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            disabled={uploading}
          />
        </div>
        <button 
          type="submit" 
          className="send-button" 
          disabled={uploading || (!message.trim() && !imagePreview)}
          title="Send message"
        >
          {uploading ? (
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="10" strokeDasharray="31.416" strokeDashoffset="31.416">
                <animate attributeName="stroke-dasharray" dur="2s" values="0 31.416;15.708 15.708;0 31.416;0 31.416" repeatCount="indefinite"/>
                <animate attributeName="stroke-dashoffset" dur="2s" values="0;-31.416;-15.708;-31.416;0" repeatCount="indefinite"/>
              </circle>
            </svg>
          ) : (
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <line x1="22" y1="2" x2="11" y2="13"></line>
              <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
            </svg>
          )}
        </button>
      </form>
    </div>
  )
}


