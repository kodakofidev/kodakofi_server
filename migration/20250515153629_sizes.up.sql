CREATE TABLE IF NOT EXISTS sizes (
  id serial PRIMARY KEY,
  size varchar UNIQUE NOT NULL,
  added_price float NOT null,
  created_at timestamptz DEFAULT now() NOT NULL,
  updated_at timestamptz
);