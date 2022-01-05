CREATE TABLE IF NOT EXISTS conflicting_files
(
    filename TEXT,
    conflicting_check_id TEXT,
    CONSTRAINT fk_conflicting_check_id
        FOREIGN KEY (conflicting_check_id)
            REFERENCES conflicting_check (id),
    PRIMARY KEY(filename, conflicting_check_id)
);
CREATE INDEX conflicting_files_check_id_idx
    ON conflicting_files (conflicting_check_id);
