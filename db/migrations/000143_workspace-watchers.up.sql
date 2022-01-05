CREATE TABLE workspace_watchers (
    workspace_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX workspace_watchers_workspace_id_idx ON workspace_watchers (workspace_id);
CREATE INDEX workspace_watchers_status_idx ON workspace_watchers(status);
