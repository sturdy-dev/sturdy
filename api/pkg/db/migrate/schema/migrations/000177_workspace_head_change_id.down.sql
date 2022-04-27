ALTER TABLE workspaces
    ADD COLUMN head_commit_id TEXT,
    ADD COLUMN head_commit_show BOOLEAN NOT NULL DEFAULT FALSE,
    DROP COLUMN head_change_id,
    DROP COLUMN head_change_computed;