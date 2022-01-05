CREATE TABLE IF NOT EXISTS oauth_user
(
    id           text      NOT NULL,
    username     text      NOT NULL,
    provider     text      NOT NULL,
    user_id      text      NOT NULL,
    access_token text      NOT NULL,
    created_at   timestamp NOT NULL,
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id),
    PRIMARY KEY (id),
    UNIQUE (username, provider)
);