ALTER TABLE installation_statistics 
    ADD COLUMN users_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN codebases_count INTEGER NOT NULL DEFAULT 0;
