CREATE TABLE IF NOT EXISTS gh_integrations
(
    id      bigint NOT NULL PRIMARY KEY,
    user_id text   NOT NULL,
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
);
