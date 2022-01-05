CREATE UNIQUE INDEX codebase_users_user_id_codebase_id_uq_idx ON
    codebase_users (user_id, codebase_id);
