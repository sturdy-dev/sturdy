ALTER TABLE users
	ADD COLUMN email_verified BOOLEAN;

UPDATE users
	SET email_verified = true;

ALTER TABLE users
	ALTER COLUMN email_verified SET NOT NULL;
