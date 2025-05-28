CREATE TABLE IF NOT EXISTS profiles (
  user_id uuid PRIMARY KEY UNIQUE NOT NULL,
  fullname varchar DEFAULT '',
  phone varchar null DEFAULT '',
  address varchar DEFAULT '',
  image varchar not NULL DEFAULT 'avatar_default.webp',
  created_at timestamptz DEFAULT now() NOT NULL,
  updated_at timestamptz,
  CONSTRAINT fk_profile_user FOREIGN KEY (user_id) REFERENCES users (id)
);