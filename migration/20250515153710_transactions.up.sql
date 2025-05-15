CREATE TABLE IF NOT EXISTS transactions (
  transaction_code varchar PRIMARY KEY,
  order_id int NOT NULL,
  grand_total int NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz,
  CONSTRAINT fk_transactions_order FOREIGN KEY (order_id) REFERENCES orders (id)
);