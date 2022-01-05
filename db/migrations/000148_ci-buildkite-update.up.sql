ALTER TABLE ci_configurations_buildkite
    DROP COLUMN org_slug,
    DROP COLUMN token,
    ADD  COLUMN organization_slug TEXT NOT NULL,
    ADD  COLUMN api_token         TEXT NOT NULL,
    ADD  COLUMN webhook_secret    TEXT NOT NULL;
