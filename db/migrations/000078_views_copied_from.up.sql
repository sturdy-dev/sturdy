ALTER TABLE views
    ADD COLUMN copied_from_branch_name TEXT;

ALTER TABLE snapshots
    ADD COLUMN view_copied_from_branch_name TEXT;