-- /shrug
TRUNCATE TABLE snapshots;

ALTER TABLE snapshots
    ADD COLUMN commit_id TEXT NOT NULL;