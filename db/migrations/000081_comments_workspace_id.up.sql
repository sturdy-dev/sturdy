ALTER TABLE comments
    ALTER COLUMN change_id DROP NOT NULL,
    ADD COLUMN workspace_id TEXT,
    ADD COLUMN context TEXT,
    ADD COLUMN context_starts_at_line INT;