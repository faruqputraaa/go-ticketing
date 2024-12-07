BEGIN;

CREATE TABLE IF NOT EXISTS offers(
    id_offer BIGSERIAL PRIMARY KEY,
    id_user INT NOT NULL,
    email VARCHAR(255) NOT NULL,
    name_event VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,

    FOREIGN KEY (id_user) REFERENCES users(id_user)
);

COMMIT;