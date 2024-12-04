BEGIN;

CREATE TABLE IF NOT EXISTS events(
    id_event BIGSERIAL PRIMARY KEY ,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    time DATE NOT NULL,
    description VARCHAR(255) NOT NULL   
);

COMMIT;