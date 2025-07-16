CREATE TABLE IF NOT EXISTS courses (
  id UUID PRIMARY KEY,
  slug TEXT NOT NULL UNIQUE,
  title TEXT NOT NULL,
  description TEXT,
  author_id UUID REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW()
);
