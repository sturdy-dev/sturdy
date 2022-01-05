CREATE TABLE acl_requested_access
(
    id         SERIAL PRIMARY KEY,
    email      TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL
);
