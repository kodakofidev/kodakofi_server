CREATE TABLE IF NOT EXISTS products_orders (
  order_id int NOT NULL,
  product_id uuid NOT NULL,  
  qty int NOT NULL,
  CONSTRAINT products_orders_pk PRIMARY KEY (order_id, product_id),
  CONSTRAINT fk_products_orders_order FOREIGN KEY (order_id) REFERENCES orders (id),
  CONSTRAINT fk_products_orders_product FOREIGN KEY (product_id) REFERENCES products (id)
);