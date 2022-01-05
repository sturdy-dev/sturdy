CREATE UNIQUE INDEX ON conflicting_check (repository_id, base, onto, user_id);
DROP INDEX conflicting_check_base_onto_user_id_idx;
