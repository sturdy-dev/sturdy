ALTER TABLE workspaces
    DROP COLUMN unarchived_at;

ALTER TABLE workspaces RENAME COLUMN archived_at TO deleted_at;
