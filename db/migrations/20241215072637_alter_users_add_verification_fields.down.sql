BEGIN;

ALTER TABLE users
    DROP COLUMN IF EXISTS reset_password_token,
    DROP COLUMN IF EXISTS verify_email_token,
    DROP COLUMN IF EXISTS is_verified;

COMMIT;
