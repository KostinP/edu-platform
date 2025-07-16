CREATE TABLE tests (
    id UUID PRIMARY KEY,
    author_id UUID REFERENCES users(id),
    title TEXT NOT NULL,
    time_limit INT DEFAULT 0,
    shuffle BOOLEAN DEFAULT FALSE,
    attempts INT DEFAULT 1,
    show_score BOOLEAN DEFAULT TRUE,
    show_answer BOOLEAN DEFAULT FALSE,
    access_from TIMESTAMPTZ,
    access_to TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE questions (
    id UUID PRIMARY KEY,
    test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
    author_id UUID REFERENCES users(id),
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    image_url TEXT,
    data JSONB NOT NULL,
    feedback TEXT,
    score FLOAT DEFAULT 1.0,
    ordinal INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE test_sessions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
  started_at TIMESTAMP NOT NULL DEFAULT NOW(),
  finished_at TIMESTAMP,
  score FLOAT,
  attempts INTEGER DEFAULT 1,
  status TEXT DEFAULT 'in_progress',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_answers (
  id UUID PRIMARY KEY,
  test_session_id UUID NOT NULL REFERENCES test_sessions(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  answer JSONB NOT NULL,
  duration_seconds INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW()
);
