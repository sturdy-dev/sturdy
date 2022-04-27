ALTER TABLE github_pull_requests
    ADD COLUMN state TEXT;

UPDATE github_pull_requests
    SET state = 'open'
    WHERE open = TRUE;

UPDATE github_pull_requests
    SET state = 'merged'
    WHERE merged = TRUE;

UPDATE github_pull_requests
    SET state = 'closed'
    WHERE merged = FALSE AND open = FALSE;
