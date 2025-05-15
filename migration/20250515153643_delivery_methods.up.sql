CREATE TABLE IF NOT EXISTS delivery_methods (
  id serial PRIMARY KEY,
  name varchar UNIQUE NOT NULL,
  fee int NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz
);