CREATE TABLE IF NOT EXISTS leads (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL UNIQUE,
  linkedin_url TEXT,
  profile TEXT,
  goal TEXT,
  maturity TEXT,
  challenge TEXT,
  lang TEXT,
  utm TEXT,
  created_at TEXT NOT NULL
);
