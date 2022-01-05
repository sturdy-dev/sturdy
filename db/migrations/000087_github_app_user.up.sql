CREATE TABLE IF NOT EXISTS github_users
(
    id                             text      NOT NULL,
    username                       text      NOT NULL,
    user_id                        text      NOT NULL,
    access_token                   text      NOT NULL,
    created_at                     timestamp NOT NULL,
    access_token_last_validated_at timestamp NOT NULL,
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id),
    PRIMARY KEY (id),
    UNIQUE (username)
);
