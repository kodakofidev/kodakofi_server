CREATE TABLE IF NOT EXISTS orders (
  id serial PRIMARY KEY UNIQUE,
  user_id uuid NOT NULL,
  fullname varchar NOT NULL,
  address varchar NOT NULL,
  delivery_method_id int NOT NULL,
  payment_method_id int NOT NULL,
  status_id int NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz,
  CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users (id),
  CONSTRAINT fk_orders_delivery FOREIGN KEY (delivery_method_id) REFERENCES delivery_methods (id),
  CONSTRAINT fk_orders_payment FOREIGN KEY (payment_method_id) REFERENCES payment_methods (id),
  CONSTRAINT fk_orders_status FOREIGN KEY (status_id) REFERENCES status (id)
);