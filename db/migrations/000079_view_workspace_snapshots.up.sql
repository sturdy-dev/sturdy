CREATE TABLE view_workspace_snapshots
(
    id              TEXT PRIMARY KEY,
    view_id         TEXT NOT NULL,
    workspace_id    TEXT NOT NULL,
    snapshot_id     TEXT NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX view_workspace_snapshots_view_id_workspace_id_idx ON
    view_workspace_snapshots (view_id, workspace_id);
