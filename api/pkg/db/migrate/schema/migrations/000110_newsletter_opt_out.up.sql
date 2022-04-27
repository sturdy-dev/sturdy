CREATE TABLE notification_settings
(
    user_id            TEXT PRIMARY KEY,
    receive_newsletter BOOLEAN DEFAULT FALSE NOT NULL
);