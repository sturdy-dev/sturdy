TRUNCATE remotes;

ALTER TABLE remotes
    ADD COLUMN browser_link_repo TEXT NOT NULL,
    ADD COLUMN browser_link_branch TEXT NOT NULL;