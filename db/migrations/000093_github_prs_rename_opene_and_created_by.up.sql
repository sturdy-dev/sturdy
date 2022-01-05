ALTER TABLE github_prs
    RENAME COLUMN opened TO open;
ALTER TABLE github_prs
    RENAME COLUMN github_author_user_login TO created_by;
