ALTER TABLE snapshots
    ADD COLUMN workspace_id TEXT;

CREATE INDEX snapshots_workspace_id_view_id_idx
    ON snapshots(workspace_id, view_id);