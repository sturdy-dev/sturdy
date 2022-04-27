ALTER TABLE github_pull_requests
    ADD COLUMN codebase_id TEXT;

UPDATE github_pull_requests
	SET (codebase_id) = (
		SELECT workspaces.codebase_id FROM workspaces WHERE id = github_pull_requests.workspace_id
	);

ALTER TABLE github_pull_requests
    ALTER COLUMN codebase_id SET NOT NULL;
