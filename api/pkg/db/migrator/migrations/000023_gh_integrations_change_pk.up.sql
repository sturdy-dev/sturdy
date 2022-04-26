ALTER TABLE gh_codebases
    DROP CONSTRAINT fk_gh_integration_id;

ALTER TABLE gh_repos
    DROP CONSTRAINT fk_gh_integration_id;

ALTER TABLE gh_integrations
    DROP constraint gh_integrations_pkey;

ALTER TABLE gh_integrations
    ADD PRIMARY KEY (id, user_id);

