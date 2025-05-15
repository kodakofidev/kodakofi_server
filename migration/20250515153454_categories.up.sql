CREATE TABLE IF NOT EXISTS categories (
  id serial PRIMARY KEY,
  name varchar NOT NULL UNIQUE,
  created_at timestamptz DEFAULT now() NOT NULL,
  updated_at timestamptz
);