ALTER TABLE snapshots
    DROP COLUMN new_files,
    DROP COLUMN changed_files,
    DROP COLUMN deleted_files;