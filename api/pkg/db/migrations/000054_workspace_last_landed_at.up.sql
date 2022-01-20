ALTER TABLE workspaces
    ADD COLUMN last_landed_at TIMESTAMP,
    ADD COLUMN deleted_at TIMESTAMP;
