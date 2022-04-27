CREATE TABLE IF NOT EXISTS codebase_users
(
    id          text NOT NULL PRIMARY KEY,
    user_id     text NOT NULL,
    codebase_id text NOT NULL,
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id),
    CONSTRAINT fk_codebase_id
        FOREIGN KEY (codebase_id)
            REFERENCES codebases (id)
);

CREATE INDEX codebase_users_user_id_idx
    ON views (user_id);

CREATE INDEX codebase_users_codebase_id_idx
    ON views (codebase_id);
