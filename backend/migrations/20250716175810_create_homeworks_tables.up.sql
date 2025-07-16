CREATE TABLE homeworks (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  type TEXT NOT NULL, -- "file", "text", "link", "code"
  course_id UUID REFERENCES courses(id),
  module_id UUID REFERENCES modules(id),
  lesson_id UUID REFERENCES lessons(id),
  group_id UUID, -- TODO: сделать таблицу групп в будущем
  user_id UUID REFERENCES users(id), -- индивидуальная выдача
  author_id UUID REFERENCES users(id),
  due_at TIMESTAMP,
  is_required BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);

CREATE TABLE homework_submissions (
  id UUID PRIMARY KEY,
  homework_id UUID REFERENCES homeworks(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id),
  status TEXT DEFAULT 'submitted', -- "submitted", "reviewed", "returned"
  answer TEXT,
  file_url TEXT,
  review TEXT,
  score FLOAT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
