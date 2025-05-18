CREATE TABLE IF NOT EXISTS profiles (
  user_id uuid PRIMARY KEY NOT NULL,
  fullname varchar,
  phone varchar null,
  address varchar,
  image varchar,
  created_at timestamptz DEFAULT now() NOT NULL,
  updated_at timestamptz,
  CONSTRAINT fk_profile_user FOREIGN KEY (user_id) REFERENCES users (id)
);