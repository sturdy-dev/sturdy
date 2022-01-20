ALTER TABLE github_repositories
    ADD COLUMN last_push_error_message TEXT,
    ADD COLUMN last_push_at TIMESTAMP WITH TIME ZONE;
