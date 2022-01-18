ALTER TABLE github_app_installed_repositories ADD COLUMN default_branch TEXT;

-- Don't worry, this value will be updated on the next push webhook from the server
UPDATE github_app_installed_repositories SET default_branch = 'master';

ALTER TABLE github_app_installed_repositories ALTER COLUMN default_branch SET NOT NULL;

ALTER TABLE conflicting_check RENAME TO conflicting_check_log;

CREATE TABLE conflicting_check (
    id TEXT NOT NULL,
    repository_id  TEXT NOT NULL,
    base  TEXT NOT NULL,
    onto  TEXT NOT NULL,
    onto_name TEXT,
    conflicting BOOLEAN NOT NULL,
    is_conflict_in_working_directory BOOLEAN NOT NULL,
    conflicting_commit TEXT,
    commit_message TEXT,
    user_id TEXT NOT NULL,
    checked_at TIMESTAMP NOT NULL,
    PRIMARY KEY(id)
);

CREATE UNIQUE INDEX ON conflicting_check (base, onto, user_id);

ALTER TABLE conflicting_files drop constraint fk_conflicting_check_id;

TRUNCATE TABLE conflicting_files;

ALTER TABLE conflicting_files ADD CONSTRAINT fk_conflicting_check_id
        FOREIGN KEY (conflicting_check_id)
            REFERENCES conflicting_check (id);
