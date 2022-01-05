CREATE TABLE IF NOT EXISTS issue_comments (
    repository_id TEXT NOT NULL,
    pr_number INT NOT NULL,
    repo_owner TEXT NOT NULL,
    repo_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    number_of_conflicts INT NOT NULL,
    comment TEXT NOT NULL,
    github_comment_id BIGINT NOT NULL,
    PRIMARY KEY(repository_id, pr_number)
);