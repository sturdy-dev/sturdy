CREATE TABLE IF NOT EXISTS pr_check_events (
    repository_id TEXT NOT NULL,
    pr_number_left INT NOT NULL,
    pr_number_right INT NOT NULL,
    processed_at TIMESTAMP,
    installation_id BIGINT NOT NULL,
    github_id BIGINT NOT NULL,
    has_conflict BOOLEAN NOT NULL,
    conflict_resolved BOOLEAN NOT NULL,
    PRIMARY KEY(repository_id, pr_number_left, pr_number_right)
);

CREATE INDEX pr_check_events_processed_at_idx
    ON pr_check_events (processed_at);
