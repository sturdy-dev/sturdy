ALTER TABLE views
    DROP COLUMN copied_from_branch_name;

ALTER TABLE snapshots
    ADD COLUMN view_copied_from_branch_name;