CREATE TABLE IF NOT EXISTS conflict_repo
(
    id         text      NOT NULL,
    full_name  text      NOT NULL,
    created_at timestamp NOT NULL,
    PRIMARY KEY (id)
);