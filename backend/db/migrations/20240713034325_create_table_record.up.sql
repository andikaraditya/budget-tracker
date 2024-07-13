CREATE TABLE IF NOT EXISTS "record" (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  amount NUMERIC(10, 2) NOT NULL,
  description TEXT,
  type VARCHAR(30) NOT NULL,
  category_id uuid NOT NULL,
  source_id VARCHAR(40) NOT NULL,
  user_id VARCHAR(255) NOT NULL,
  date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY ("category_id") REFERENCES "category" (id),
  FOREIGN KEY ("source_id") REFERENCES "source" (id),
  FOREIGN KEY (user_id) REFERENCES "user" (id)
)