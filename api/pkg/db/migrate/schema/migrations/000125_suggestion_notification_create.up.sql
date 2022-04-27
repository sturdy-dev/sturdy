CREATE TABLE suggestion_notifications (
    user_id      TEXT      NOT NULL,
    workspace_id TEXT      NOT NULL
);

CREATE UNIQUE INDEX suggestion_notifications_user_id_workspace_unique_ix ON
    suggestion_notifications(user_id, workspace_id);