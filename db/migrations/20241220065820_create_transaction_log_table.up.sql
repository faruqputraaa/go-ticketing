BEGIN;

CREATE TABLE IF NOT EXISTS transaction_log(
    id_transaction_log BIGSERIAL PRIMARY KEY,
    id_transaction VARCHAR(255),
	status VARCHAR(255),
	message VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (id_transaction) REFERENCES transactions(id_transaction)
);

COMMIT;