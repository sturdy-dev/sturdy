CREATE TABLE view_status
(
    id                      TEXT PRIMARY KEY,
    state                   TEXT NOT NULL,
    staging_status_path     TEXT,
    staging_status_received INT,
    staging_status_total    INT,
    sturdy_version          TEXT NOT NULL,
    last_error              TEXT
)