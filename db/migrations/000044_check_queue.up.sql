CREATE TABLE conflicting_check_queue (
    repository_id TEXT NOT NULL,
    user_id TEXT NOT NULL ,
    head TEXT NOT NULL,
    working_directory TEXT NOT NULL,
    working_directory_hash TEXT NOT NULL,
    processed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (repository_id, user_id)
);

CREATE INDEX check_conflicts_queue_processed_at_idx
    ON conflicting_check_queue (processed_at);