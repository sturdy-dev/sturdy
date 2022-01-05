ALTER TABLE snapshots
    ADD COLUMN action TEXT NOT NULL DEFAULT 'view_sync';