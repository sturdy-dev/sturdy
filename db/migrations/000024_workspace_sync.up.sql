CREATE TABLE IF NOT EXISTS workspace_sync
(
    id                  text NOT NULL PRIMARY KEY,
    user_id             text NOT NULL,
    codebase_id         text NOT NULL,
    workspace_id        text NOT NULL,
    view_id             text NOT NULL,
    base_commit         text NOT NULL,
    onto_commit         text NOT NULL,
    current_commit      text NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_codebase_id FOREIGN KEY (codebase_id) REFERENCES codebases (id)
);
