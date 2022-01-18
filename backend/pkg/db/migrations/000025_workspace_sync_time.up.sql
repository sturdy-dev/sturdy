ALTER TABLE workspace_sync
    ADD COLUMN created_at timestamp,
    ADD COLUMN completed_at timestamp,
    ADD COLUMN aborted_at timestamp;