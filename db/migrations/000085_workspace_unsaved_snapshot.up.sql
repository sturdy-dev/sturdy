ALTER TABLE workspaces
    ADD COLUMN latest_snapshot_id TEXT;

CREATE UNIQUE INDEX workspaces_view_id_idx
    ON workspaces (view_id);

-- Migration
WITH subquery AS (
    SELECT id as view_id, workspace_id
    FROM views
)
UPDATE workspaces
SET view_id = subquery.view_id
FROM subquery
WHERE workspaces.id = subquery.workspace_id;