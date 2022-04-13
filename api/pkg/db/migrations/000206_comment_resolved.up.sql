ALTER TABLE comments
    ADD COLUMN resolved_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN resolved_by TEXT;