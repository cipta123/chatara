-- Add reply_to_id column to messages table for reply functionality
ALTER TABLE messages ADD COLUMN IF NOT EXISTS reply_to_id INTEGER REFERENCES messages(id) ON DELETE SET NULL;

-- Add deleted_at column for soft delete
ALTER TABLE messages ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- Create index for reply_to_id for better query performance
CREATE INDEX IF NOT EXISTS idx_messages_reply_to_id ON messages(reply_to_id);

-- Create index for deleted_at for filtering deleted messages
CREATE INDEX IF NOT EXISTS idx_messages_deleted_at ON messages(deleted_at);

