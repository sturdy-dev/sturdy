ALTER TABLE github_prs
    RENAME COLUMN open TO opened;
ALTER TABLE github_prs
    RENAME COLUMN created_by TO github_author_user_login;
