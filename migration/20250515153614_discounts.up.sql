CREATE TABLE IF NOT EXISTS discounts (
  id serial PRIMARY key not null,
  name varchar UNIQUE NOT NULL,
  discount float NOT NULL,
  expired timestamptz,
  created_at timestamptz DEFAULT now() NOT NULL,
  updated_at timestamptz
);