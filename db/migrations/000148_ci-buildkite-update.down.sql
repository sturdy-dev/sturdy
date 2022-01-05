ALTER TABLE ci_configurations_buildkite
    ADD  COLUMN org_slug TEXT NOT NULL,
    ADD  COLUMN token    TEXT NOT NULL,
    DROP COLUMN organization_slug,
    DROP COLUMN api_token,
    DROP COLUMN webhook_secret;
