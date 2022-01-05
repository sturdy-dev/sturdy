ALTER TABLE github_pull_requests ADD COLUMN head_sha TEXT;
UPDATE github_pull_requests SET head_sha = 'from-a-migration';
ALTER TABLE github_pull_requests ALTER COLUMN head_sha SET NOT NULL;

ALTER TABLE conflicting_check
    ADD COLUMN onto_reference_type TEXT,
    ADD COLUMN onto_reference_id TEXT;
