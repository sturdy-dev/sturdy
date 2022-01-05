CREATE TABLE IF NOT EXISTS users
(
    id         text PRIMARY KEY,
    first_name text NOT NULL,
    last_name  text NOT NULL,
    email      text NOT NULL
);