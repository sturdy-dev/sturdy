ALTER TABLE gh_codebases
    ADD CONSTRAINT fk_gh_integration_id
        FOREIGN KEY (gh_integration_id)
            REFERENCES gh_integrations (id);

ALTER TABLE gh_repos
    ADD CONSTRAINT fk_gh_integration_id
        FOREIGN KEY (integration_id)
            REFERENCES gh_integrations (id);

ALTER TABLE gh_integrations
    DROP constraint gh_integrations_pkey;

ALTER TABLE gh_integrations
    ADD PRIMARY KEY (id);
