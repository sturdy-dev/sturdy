ALTER TABLE snapshots
    ADD COLUMN new_files TEXT[],
    ADD COLUMN changed_files TEXT[],
    ADD COLUMN deleted_files TEXT[];