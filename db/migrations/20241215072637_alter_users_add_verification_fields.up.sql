DO $$
    BEGIN
        -- Tambahkan kolom reset_password_token jika belum ada
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_name = 'users' AND column_name = 'reset_password_token'
        ) THEN
            ALTER TABLE users
                ADD COLUMN reset_password_token VARCHAR(255);
        END IF;

        -- Tambahkan kolom verify_email_token jika belum ada
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_name = 'users' AND column_name = 'verify_email_token'
        ) THEN
            ALTER TABLE users
                ADD COLUMN verify_email_token VARCHAR(255);
        END IF;

        -- Tambahkan kolom is_verified jika belum ada
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_name = 'users' AND column_name = 'is_verified'
        ) THEN
            ALTER TABLE users
                ADD COLUMN is_verified INTEGER NOT NULL DEFAULT 0;
        END IF;
    END $$;
