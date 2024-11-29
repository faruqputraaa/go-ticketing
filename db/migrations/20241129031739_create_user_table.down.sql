CREATE TABLE IF NOT EXISTS transactions(
  id_transaction INT PRIMARY KEY,
  id_user INT NOT NULL,
  quantity_ticket INT NOT NULL,
  id_event INT NOT NULL,
  id_ticket INT NOT NULL,
  total_price FLOAT,
  date_order TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,)