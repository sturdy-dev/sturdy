ALTER TABLE github_repositories
    ALTER COLUMN installation_access_token_expires_at TYPE TIMESTAMP WITH TIME ZONE USING installation_access_token_expires_at AT TIME ZONE 'UTC',
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN uninstalled_at TYPE TIMESTAMP WITH TIME ZONE USING uninstalled_at AT TIME ZONE 'UTC',
    ALTER COLUMN synced_at TYPE TIMESTAMP WITH TIME ZONE USING synced_at AT TIME ZONE 'UTC';

ALTER TABLE codebases
    ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN archived_at TYPE TIMESTAMP WITH TIME ZONE USING archived_at AT TIME ZONE 'UTC';
