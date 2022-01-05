CREATE TABLE comments
(
    id          TEXT PRIMARY KEY,
    codebase_id TEXT NOT NULL,
    change_id   TEXT NOT NULL,
    user_id     TEXT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at  TIMESTAMP WITH TIME ZONE, -- Nullable
    message     TEXT NOT NULL,
    path        TEXT NOT NULL,
    line_start  INT NOT NULL,
    line_end    INT NOT NULL,
    line_is_new BOOLEAN NOT NULL
);

CREATE INDEX comments_codebase_id_change_id_idx
    ON comments(codebase_id, change_id);