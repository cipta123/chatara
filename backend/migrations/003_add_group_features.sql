-- Add group features to conversations table
ALTER TABLE conversations
ADD COLUMN IF NOT EXISTS name VARCHAR(255),
ADD COLUMN IF NOT EXISTS description TEXT;

-- Add role to participants table (admin/member)
ALTER TABLE participants
ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'member';

-- Add permission to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS can_create_group BOOLEAN DEFAULT FALSE;

-- Set permission for specific user (ciptaanugrahh@gmail.com)
UPDATE users
SET can_create_group = TRUE
WHERE email = 'ciptaanugrahh@gmail.com';

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_conversations_name ON conversations(name);
CREATE INDEX IF NOT EXISTS idx_participants_role ON participants(role);
CREATE INDEX IF NOT EXISTS idx_users_can_create_group ON users(can_create_group);

-- Set default role for existing participants
UPDATE participants
SET role = 'member'
WHERE role IS NULL;

