-- Add name field to users table (idempotent)
ALTER TABLE users ADD COLUMN IF NOT EXISTS name VARCHAR(255);
