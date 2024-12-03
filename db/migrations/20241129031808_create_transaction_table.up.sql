BEGIN;

CREATE TABLE IF NOT EXISTS transactions(
    id_transaction INT PRIMARY KEY AUTO_INCREMENT,
    id_user INT NOT NULL,
    quantity_ticket INT NOT NULL,
    id_event INT NOT NULL,
    id_ticket INT NOT NULL,
    total_price FLOAT,
    date_order TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (id_user) REFERENCES users(id_user),
    FOREIGN KEY (id_event) REFERENCES events(id_events),
    FOREIGN KEY (id_ticket) REFERENCES tickets(id_ticket)
);

COMMIT;