CREATE TABLE IF NOT EXISTS payment_methods (
  id serial PRIMARY KEY,
  name varchar UNIQUE NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz
);