CREATE TABLE IF NOT EXISTS transactions (
  transaction_code varchar PRIMARY KEY UNIQUE,
  order_id int NOT NULL,
  total int NOT NULL,
  delivery_fee int NOT NULL,
  tax int NOT null,
  total_amount int NOT null,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz,
  CONSTRAINT fk_transactions_order FOREIGN KEY (order_id) REFERENCES orders (id)
);