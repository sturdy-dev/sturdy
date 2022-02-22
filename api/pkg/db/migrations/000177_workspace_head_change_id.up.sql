ALTER TABLE workspaces
    ADD COLUMN head_change_id TEXT,
    ADD COLUMN head_change_computed BOOLEAN NOT NULL DEFAULT FALSE,
    DROP COLUMN head_commit_id,
    DROP COLUMN head_commit_show;