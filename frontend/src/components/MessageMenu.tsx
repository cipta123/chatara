import './MessageMenu.css'

interface MessageMenuProps {
  isOpen: boolean
  onClose: () => void
  onCopy: () => void
  onReply: () => void
  onDelete: () => void
  position: { top: number; left: number }
  isOwn: boolean
}

export default function MessageMenu({
  isOpen,
  onClose,
  onCopy,
  onReply,
  onDelete,
  position,
  isOwn,
}: MessageMenuProps) {
  if (!isOpen) return null

  const handleCopy = () => {
    onCopy()
    onClose()
  }

  const handleReply = () => {
    onReply()
    onClose()
  }

  const handleDelete = () => {
    onDelete()
    onClose()
  }

  return (
    <>
      <div className="message-menu-overlay" onClick={onClose} />
      <div
        className="message-menu"
        style={{
          top: `${position.top}px`,
          left: `${position.left}px`,
        }}
      >
        <button className="message-menu-item" onClick={handleCopy}>
          <svg className="menu-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
          </svg>
          <span>Copy</span>
        </button>
        <button className="message-menu-item" onClick={handleReply}>
          <svg className="menu-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <polyline points="9 10 4 15 9 20"></polyline>
            <path d="M20 4v7a4 4 0 0 1-4 4H4"></path>
          </svg>
          <span>Reply</span>
        </button>
        {isOwn && (
          <button className="message-menu-item message-menu-item-danger" onClick={handleDelete}>
            <svg className="menu-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            </svg>
            <span>Delete</span>
          </button>
        )}
      </div>
    </>
  )
}

