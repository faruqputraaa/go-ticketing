BEGIN;

CREATE TABLE IF NOT EXISTS transactions(
    id_transaction VARCHAR(255) PRIMARY KEY,
    id_user INT NOT NULL,
    quantity_ticket INT NOT NULL,
    id_ticket INT NOT NULL,
    status VARCHAR(255) NOT NULL,
    total_price FLOAT,
    date_order TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (id_user) REFERENCES users(id_user),
    FOREIGN KEY (id_ticket) REFERENCES tickets(id_ticket)
);

COMMIT;