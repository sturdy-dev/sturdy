CREATE TABLE waitinglist
(
    id         SERIAL PRIMARY KEY,
    email      TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL
);