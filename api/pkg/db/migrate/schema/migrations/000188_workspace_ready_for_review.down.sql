ALTER TABLE workspaces
    ADD COLUMN ready_for_review_change TEXT,
    ADD COLUMN approved_change TEXT;