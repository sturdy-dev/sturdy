ALTER TABLE workspaces
    DROP COLUMN latest_snapshot_id;

DROP INDEX workspaces_view_id_idx;