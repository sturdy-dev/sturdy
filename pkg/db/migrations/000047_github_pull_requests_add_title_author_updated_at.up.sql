ALTER TABLE github_pull_requests ADD COLUMN title TEXT;
ALTER TABLE github_pull_requests ADD COLUMN author TEXT;
ALTER TABLE github_pull_requests ADD COLUMN updated_at TIMESTAMP;
