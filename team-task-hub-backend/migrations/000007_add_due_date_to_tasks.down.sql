-- Rollback: Remove due_date column from tasks table
ALTER TABLE tasks
DROP COLUMN IF EXISTS due_date;
