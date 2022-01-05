CREATE TABLE workspace_sync
(
    id             text NOT NULL,
    user_id        text NOT NULL,
    codebase_id    text NOT NULL,
    workspace_id   text NOT NULL,
    view_id        text NOT NULL,
    base_commit    text NOT NULL,
    onto_commit    text NOT NULL,
    current_commit text NOT NULL,
    created_at     timestamp without time zone,
    completed_at   timestamp without time zone,
    aborted_at     timestamp without time zone,
    unsaved_commit text
);