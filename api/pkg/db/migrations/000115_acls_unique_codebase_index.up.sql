DROP INDEX acls_codebase_id_idx;

CREATE UNIQUE INDEX acls_unique_codebase_id_idx ON
    acls (codebase_id);
