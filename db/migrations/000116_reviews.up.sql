CREATE TABLE workspace_reviews
(
    id           TEXT PRIMARY KEY,
    codebase_id  TEXT                     NOT NULL,
    workspace_id TEXT                     NOT NULL,
    user_id      TEXT                     NOT NULL,
    grade        TEXT                     NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    dismissed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX workspace_reviews_workspace_id_idx
    ON workspace_reviews (workspace_id);