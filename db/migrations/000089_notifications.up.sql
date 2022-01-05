CREATE TABLE notifications
(
    id           TEXT PRIMARY KEY,
    codebase_id  TEXT NOT NULL,
    user_id      TEXT NOT NULL,
    type         TEXT NOT NULL,
    reference_id TEXT NOT NULL,
    archived_at  TIMESTAMP WITH TIME ZONE
);

CREATE INDEX notifications_user_id_codebase_id_archived_at_idx
    ON notifications (user_id, codebase_id, archived_at);

CREATE INDEX notifications_user_id_archived_at_idx
    ON notifications (user_id, archived_at);