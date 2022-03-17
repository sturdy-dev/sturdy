CREATE TABLE remotes
(
    id             TEXT PRIMARY KEY,
    codebase_id    TEXT NOT NULL,
    name           TEXT,
    url            TEXT NOT NULL,
    basic_username TEXT,
    basic_password TEXT,
    tracked_branch TEXT NOT NULL
)