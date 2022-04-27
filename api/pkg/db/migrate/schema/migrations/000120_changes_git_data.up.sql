ALTER TABLE changes
    ADD COLUMN created_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN git_created_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN git_creator_name TEXT,
    ADD COLUMN git_creator_email TEXT;
