CREATE TABLE IF NOT EXISTS product_discounts (
  product_id uuid NOT NULL,
  discount_id int NOT NULL,  
  CONSTRAINT product_discount_pk PRIMARY KEY (product_id, discount_id),
  CONSTRAINT fk_product_discount_product FOREIGN KEY (product_id) REFERENCES products (id),
  CONSTRAINT fk_product_discount_discount FOREIGN KEY (discount_id) REFERENCES discounts (id)
);