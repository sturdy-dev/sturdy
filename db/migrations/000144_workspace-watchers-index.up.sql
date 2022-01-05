CREATE INDEX workspace_watchers_workspace_id_user_id_idx ON workspace_watchers (workspace_id, user_id);

DROP INDEX workspace_watchers_workspace_id_idx;
