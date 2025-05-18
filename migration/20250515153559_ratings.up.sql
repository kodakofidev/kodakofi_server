CREATE TABLE IF NOT EXISTS ratings (
  user_id uuid NOT NULL,
  product_id uuid NOT NULL,
  rating int NOT NULL,
  created_at timestamptz DEFAULT now() NOT NULL,
  updated_at timestamptz,
  CONSTRAINT ratings_pk PRIMARY KEY (user_id, product_id),
  CONSTRAINT fk_rating_product FOREIGN KEY (product_id) REFERENCES products (id),
  CONSTRAINT fk_rating_user FOREIGN KEY (user_id) REFERENCES users (id)
);

ALTER TABLE ratings
ALTER COLUMN rating TYPE BOOLEAN;

ALTER TABLE ratings 
DROP CONSTRAINT fk_rating_user;
