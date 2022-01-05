CREATE TABLE notification_preferences (
    user_id TEXT    NOT NULL,
    type    TEXT    NOT NULL,
    channel TEXT    NOT NULL,
    enabled BOOLEAN NOT NULL
);

ALTER TABLE notification_preferences
    ADD CONSTRAINT user_id_type_channel_uq_ix UNIQUE (user_id, type, channel);
