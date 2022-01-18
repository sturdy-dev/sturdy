-- /shrug
TRUNCATE TABLE snapshots;

ALTER TABLE snapshots
    ADD COLUMN codebase_id TEXT NOT NULL;