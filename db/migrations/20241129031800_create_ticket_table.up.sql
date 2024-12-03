BEGIN;

CREATE TABLE IF NOT EXISTS tickets(
    id_ticket BIGSERIAL PRIMARY KEY,
    price FLOAT NOT NULL,
    category VARCHAR(255) NOT NULL,
    id_event INT NOT NULL,

    FOREIGN KEY (id_event) REFERENCES events(id_event)
);

COMMIT;