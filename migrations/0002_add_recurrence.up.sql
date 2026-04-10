-- Add scheduled_at to tasks so generated instances know their target date.
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS scheduled_at DATE NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_scheduled_at ON tasks (scheduled_at);
