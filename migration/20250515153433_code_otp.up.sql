CREATE TABLE IF NOT EXISTS code_otp (
  id serial PRIMARY KEY,
  user_id uuid NOT NULL,
  code varchar NOT NULL, -- Changed from int to varchar
  type_id int NOT NULL,
  expired_at timestamptz NOT NULL,
  created_at timestamptz DEFAULT now() NOT NULL,
  updated_at timestamptz,
  CONSTRAINT fk_code_otp_user FOREIGN KEY (user_id) REFERENCES users (id),
  CONSTRAINT fk_code_otp_type FOREIGN KEY (type_id) REFERENCES otp_type (id)
);