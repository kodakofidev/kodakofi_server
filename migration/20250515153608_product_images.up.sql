CREATE TABLE IF NOT EXISTS product_images (
  product_id uuid NOT NULL,
  path varchar NOT NULL DEFAULT '',
  created_at timestamptz DEFAULT now() NOT NULL,
  updated_at timestamptz,
  CONSTRAINT product_images_pk PRIMARY KEY (product_id, path),
  CONSTRAINT fk_product_images_product FOREIGN KEY (product_id) REFERENCES products (id)
);