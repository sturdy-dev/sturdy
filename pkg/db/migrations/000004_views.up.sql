CREATE TABLE IF NOT EXISTS views
(
    id          text PRIMARY KEY,
    user_id     text NOT NULL,
    codebase_id text NOT NULL,
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id),
    CONSTRAINT fk_codebase_id
        FOREIGN KEY (codebase_id)
            REFERENCES codebases (id)
);
CREATE INDEX views_user_id_codebase_id_idx
    ON views (user_id, codebase_id);
