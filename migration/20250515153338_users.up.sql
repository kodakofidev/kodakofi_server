CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
  email varchar UNIQUE NOT NULL,
  password varchar NOT NULL,
  role varchar DEFAULT 'user' NOT NULL,
  is_verified bool DEFAULT false NOT NULL,
  created_at timestamptz DEFAULT now() NOT NULL,
  updated_at timestamptz
);