CREATE TABLE IF NOT EXISTS tax (
  id serial PRIMARY KEY NOT NULL UNIQUE,
  tax_value float NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz
);