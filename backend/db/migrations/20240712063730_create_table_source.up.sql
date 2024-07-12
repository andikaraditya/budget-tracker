CREATE TABLE IF NOT EXISTS "source" (
    id VARCHAR(40) PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    description TEXT,
    initial NUMERIC(10,2),
    user_id VARCHAR(40) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES "user" (id)
)