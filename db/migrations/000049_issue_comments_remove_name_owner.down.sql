ALTER TABLE issue_comments ADD COLUMN repo_owner TEXT NOT NULL;
ALTER TABLE issue_comments ADD COLUMN repo_name TEXT NOT NULL;
ALTER TABLE issue_comments ADD COLUMN number_of_conflicts INT NOT NULL;
