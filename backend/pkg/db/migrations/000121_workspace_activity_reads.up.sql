CREATE TABLE workspace_activity_reads
(
    id                   TEXT PRIMARY KEY,
    workspace_id         TEXT                     NOT NULL,
    user_id              TEXT                     NOT NULL,
    last_read_created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX workspace_activity_reads_workspace_id_user_id_idx
    ON workspace_activity_reads (workspace_id, user_id);