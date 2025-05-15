CREATE TABLE IF NOT EXISTS size_products (
  product_id uuid NOT NULL,
  size_id int NOT NULL,
  created_at timestamptz DEFAULT now() NOT NULL,
  updated_at timestamptz,
  CONSTRAINT size_product_pk PRIMARY KEY (product_id, size_id),
  CONSTRAINT fk_size_products_product FOREIGN KEY (product_id) REFERENCES products (id),
  CONSTRAINT fk_size_products_size FOREIGN KEY (size_id) REFERENCES sizes (id)
);