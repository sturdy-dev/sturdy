CREATE TABLE codebases_garbage_collection_status
(
    codebase_id                     TEXT PRIMARY KEY         NOT NULL,
    last_collected_completed_at     TIMESTAMP WITH TIME ZONE NOT NULL,
    last_collection_duration_millis INT                      NOT NULL
);
