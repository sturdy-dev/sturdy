CREATE TABLE IF NOT EXISTS conflicting_check
(
    id                 text      NOT NULL,
    repository_id      text      NOT NULL,
    base               text      NOT NULL,
    onto               text      NOT NULL,
    conflicting        boolean   NOT NULL,
    conflicting_commit text      NOT NULL,
    created_at         timestamp NOT NULL,
    created_by         text      NOT NULL,
    PRIMARY KEY (id)
);