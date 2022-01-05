CREATE TABLE IF NOT EXISTS github_pull_requests (
    repository_id TEXT NOT NULL,
    pr_number INT NOT NULL,
    is_open BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    closed_at TIMESTAMP NOT NULL,
    PRIMARY KEY(repository_id, pr_number)
);
