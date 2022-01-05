CREATE TABLE IF NOT EXISTS gh_codebases
(
    gh_integration_id bigint NOT NULL,
    codebase_id       text   NOT NULL,
    gh_owner_name     text   NOT NULL,
    gh_repo_name      text   NOT NULL,
    gh_export_branch  text   NOT NULL,
    CONSTRAINT fk_gh_integration_id
        FOREIGN KEY (gh_integration_id)
            REFERENCES gh_integrations (id),
    CONSTRAINT fk_codebase_id
        FOREIGN KEY (codebase_id)
            REFERENCES codebases (id),
    PRIMARY KEY(gh_integration_id, codebase_id)
);

CREATE INDEX gh_codebases_codebase_id_idx
    ON gh_codebases (codebase_id);
