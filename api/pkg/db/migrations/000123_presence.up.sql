CREATE TABLE presence (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    workspace_id TEXT NOT NULL,
    last_active_at TIMESTAMP WITH TIME ZONE NOT NULL,
    state TEXT
);

CREATE UNIQUE INDEX presence_user_id_workspace_id_idx
    ON presence(user_id, workspace_id);