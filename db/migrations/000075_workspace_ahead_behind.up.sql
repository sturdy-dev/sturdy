ALTER TABLE workspaces
    ADD COLUMN behind_count INT,
    ADD COLUMN ahead_count INT,
    ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE;