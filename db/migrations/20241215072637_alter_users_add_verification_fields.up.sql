BEGIN;

ALTER TABLE users
    ADD COLUMN reset_password_token VARCHAR(255),
    ADD COLUMN verify_email_token VARCHAR(255),
    ADD COLUMN is_verified INTEGER NOT NULL DEFAULT 0;

COMMIT;
