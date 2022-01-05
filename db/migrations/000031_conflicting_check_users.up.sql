CREATE TABLE IF NOT EXISTS conflicting_check_users
(
    repository_id    text      NOT NULL,
    user_id          text      NOT NULL,
    first_pushed_at  timestamp NOT NULL,
    last_pushed_at   timestamp NOT NULL,
    first_checked_at timestamp NOT NULL,
    last_checked_at  timestamp NOT NULL,
    PRIMARY KEY (repository_id, user_id)
);