ALTER TABLE snapshots
    DROP COLUMN workspace_id;

DELETE INDEX snapshots_workspace_id_view_id_idx;