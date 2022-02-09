ALTER TABLE installation_statistics 
    ALTER COLUMN first_user_email DROP NOT NULL,
    ALTER COLUMN first_user_email DROP DEFAULT;
