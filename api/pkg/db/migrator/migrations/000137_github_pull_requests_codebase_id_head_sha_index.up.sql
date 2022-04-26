CREATE INDEX github_pull_requests_codebase_id_head_sha_ix
    ON github_pull_requests (codebase_id, head_sha);
