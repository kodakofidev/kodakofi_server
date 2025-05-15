CREATE TABLE IF NOT EXISTS status (
  id serial PRIMARY KEY,
  status varchar UNIQUE NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz
);