CREATE TABLE snapshots
(
    id            TEXT PRIMARY KEY,
    view_id       TEXT                     NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    new_files     TEXT[],
    changed_files TEXT[],
    deleted_files TEXT[]
);

CREATE
INDEX snapshots_view_id_idx
    ON snapshots (view_id);