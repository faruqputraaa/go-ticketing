BEGIN;

-- Mendefinisikan Tipe ENUM untuk role
-- CREATE TYPE user_role AS ENUM ('ADMIN', 'BUYER');

CREATE TABLE IF NOT EXISTS users (
    id_user BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
