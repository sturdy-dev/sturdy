
CREATE TABLE github_app_user_access (
    user_id TEXT NOT NULL,
    repository_id TEXT NOT NULL,
    has_access BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_checked_at TIMESTAMP NOT NULL,
    PRIMARY KEY(user_id, repository_id)
);
