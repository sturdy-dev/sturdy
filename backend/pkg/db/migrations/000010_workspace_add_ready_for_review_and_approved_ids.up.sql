ALTER TABLE workspaces
    ADD COLUMN ready_for_review_change text;
ALTER TABLE workspaces
    ADD COLUMN approved_change text;
