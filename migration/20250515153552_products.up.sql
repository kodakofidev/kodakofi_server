CREATE TABLE IF NOT EXISTS products (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
  name varchar UNIQUE NOT NULL,
  category_id int NOT NULL,
  is_deleted bool DEFAULT false NOT NULL,
  stock int NOT NULL,
  created_at timestamptz DEFAULT now() NOT NULL,
  updated_at timestamptz,
  CONSTRAINT fk_product_category FOREIGN KEY (category_id) REFERENCES categories (id)
);