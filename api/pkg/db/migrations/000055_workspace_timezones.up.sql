ALTER TABLE workspaces
    DROP COLUMN last_landed_at,
    ADD COLUMN last_landed_at TIMESTAMP WITH TIME ZONE,
    DROP COLUMN created_at,
    ADD COLUMN created_at TIMESTAMP WITH TIME ZONE,
    DROP COLUMN deleted_at,
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
