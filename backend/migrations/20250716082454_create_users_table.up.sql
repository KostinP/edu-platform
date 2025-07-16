CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  telegram_id TEXT NOT NULL UNIQUE,
  first_name TEXT,
  last_name TEXT,
  username TEXT,
  photo_url TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);
