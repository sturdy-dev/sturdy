CREATE TABLE suggestions
(
    id              TEXT PRIMARY KEY,
    workspace_id    TEXT NOT NULL,
    snapshot_id     TEXT NOT NULL,
    dismissed_at    TIMESTAMP WITH TIME ZONE,
    dismissed_files TEXT[]
);

CREATE UNIQUE INDEX suggestions_workspace_id_snapshot_id_idx ON
    suggestions (workspace_id, snapshot_id);