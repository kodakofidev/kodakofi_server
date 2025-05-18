CREATE TABLE IF NOT EXISTS code_otp (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL,
  code VARCHAR(6) NOT NULL,
  type_id INT NOT NULL,
  expired_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ,
  CONSTRAINT fk_code_otp_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_code_otp_type FOREIGN KEY (type_id) REFERENCES otp_type(id) ON DELETE CASCADE,
  CONSTRAINT unique_user_type UNIQUE (user_id, type_id)
);