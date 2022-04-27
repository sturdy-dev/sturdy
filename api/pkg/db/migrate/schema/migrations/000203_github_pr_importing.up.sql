ALTER TABLE github_pull_requests
    ADD COLUMN importing BOOLEAN;

UPDATE
    github_pull_requests
SET
    importing = FALSE;

ALTER TABLE github_pull_requests
    ALTER COLUMN importing SET NOT NULL;
