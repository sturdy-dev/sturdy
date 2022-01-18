ALTER TABLE snapshots
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX snapshots_codebase_id_deleted_at
    ON snapshots (codebase_id, deleted_at);