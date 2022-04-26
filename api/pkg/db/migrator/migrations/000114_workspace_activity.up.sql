CREATE TABLE workspace_activity
(
    id            TEXT PRIMARY KEY,
    user_id       TEXT                     NOT NULL,
    workspace_id  TEXT                     NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    activity_type TEXT                     NOT NULL,
    reference     TEXT                     NOT NULL
);

CREATE INDEX workspace_activity_workspace_id_created_at_idx
    ON workspace_activity (workspace_id, created_at);