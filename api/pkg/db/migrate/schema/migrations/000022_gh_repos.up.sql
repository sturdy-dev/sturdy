CREATE TABLE IF NOT EXISTS gh_repos
(
    integration_id bigint NOT NULL,
    owner_name     text   NOT NULL,
    repo_name      text   NOT NULL,
    public_key     text   NOT NULL,
    private_key    text   NOT NULL,
    user_id        text   NOT NULL,
    CONSTRAINT fk_gh_integration_id
        FOREIGN KEY (integration_id)
            REFERENCES gh_integrations (id),
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id),
    PRIMARY KEY (integration_id, owner_name, repo_name)
);
