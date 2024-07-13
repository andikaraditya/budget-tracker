CREATE TABLE IF NOT EXISTS "transfer" (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "to_id" VARCHAR(40) NOT NULL,
  "from_id" VARCHAR(40) NOT NULL,
  amount NUMERIC(10, 2) NOT NULL,
  user_id VARCHAR(40) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY ("to_id") REFERENCES "source" (id),
  FOREIGN KEY ("from_id") REFERENCES "source" (id),
  FOREIGN KEY (user_id) REFERENCES "user" (id)
);