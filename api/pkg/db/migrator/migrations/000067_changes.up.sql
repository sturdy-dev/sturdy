CREATE TABLE changes
(
    id                  TEXT PRIMARY KEY, -- A globally unique ID
    codebase_id         TEXT NOT NULL,
    commit_id           TEXT NOT NULL,
    updated_description TEXT
);

CREATE INDEX changes_codebase_id_commit_id_idx
    ON changes (codebase_id, commit_id);
