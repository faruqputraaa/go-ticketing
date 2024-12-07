BEGIN;

-- Mendefinisikan Tipe ENUM untuk role
-- CREATE TYPE user_role AS ENUM ('ADMIN', 'BUYER');

CREATE TABLE IF NOT EXISTS users (
    id_user BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    email varchar(255) NOT NULL UNIQUE ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
