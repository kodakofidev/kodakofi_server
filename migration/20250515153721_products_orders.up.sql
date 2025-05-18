CREATE TABLE IF NOT EXISTS products_orders (
  order_id int NOT NULL,
  product_id uuid NOT NULL,
  base_price int not null,
  qty int NOT NULL,
  added_price int not null,
  sub_total int not null,
  CONSTRAINT products_orders_pk PRIMARY KEY (order_id, product_id, qty),
  CONSTRAINT fk_products_orders_order FOREIGN KEY (order_id) REFERENCES orders (id),
  CONSTRAINT fk_products_orders_product FOREIGN KEY (product_id) REFERENCES products (id)
);