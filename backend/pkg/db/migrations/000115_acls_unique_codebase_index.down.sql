CREATE INDEX acls_codebase_id_idx ON
    acls (codebase_id);

DROP INDEX acls_unique_codebase_id_idx;
