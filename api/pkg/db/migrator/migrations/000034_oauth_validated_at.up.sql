alter table oauth_user add column access_token_last_validated_at timestamp;
update oauth_user set access_token_last_validated_at = created_at;
alter table oauth_user alter column access_token_last_validated_at set not null;