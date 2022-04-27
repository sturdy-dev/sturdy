CREATE TABLE suggestions_v2 (
    id               TEXT                NOT NULL PRIMARY KEY,
    codebase_id      TEXT                NOT NULL,
    workspace_id     TEXT                NOT NULL,
    for_workspace_id TEXT                NOT NULL,
    for_snapshot_id  TEXT                NOT NULL,
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL,
    applied_hunks    TEXT[],
    dismissed_hunks  TEXT[],
    user_id          TEXT                NOT NULL
);

CREATE INDEX suggestions_v2_workspace_id_idx ON suggestions_v2 (workspace_id);
CREATE INDEX suggestions_v2_for_workspace_id_idx ON suggestions_v2 (for_workspace_id);
