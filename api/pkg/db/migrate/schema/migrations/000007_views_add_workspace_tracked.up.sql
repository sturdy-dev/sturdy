ALTER TABLE views
    ADD COLUMN workspace_id text;
ALTER TABLE views
    ADD CONSTRAINT fk_workspace_id
        FOREIGN KEY (workspace_id)
            REFERENCES workspaces (id);
