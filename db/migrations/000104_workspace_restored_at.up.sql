ALTER TABLE workspaces ADD COLUMN unarchived_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE workspaces RENAME COLUMN deleted_at TO archived_at;
