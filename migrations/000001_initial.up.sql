CREATE TABLE teams (team_name TEXT PRIMARY KEY);

CREATE TABLE users (
  user_id TEXT PRIMARY KEY,
  username TEXT NOT NULL,
  team_name TEXT NOT NULL REFERENCES teams (team_name) ON DELETE CASCADE,
  is_active BOOLEAN NOT NULL
);

CREATE TABLE pull_requests (
  pull_request_id TEXT PRIMARY KEY,
  pull_request_name TEXT NOT NULL,
  author_id TEXT NOT NULL REFERENCES users (user_id) ON DELETE RESTRICT,
  status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
  created_at TIMESTAMPTZ NULL DEFAULT NOW (),
  merged_at TIMESTAMPTZ NULL
);

CREATE TABLE pull_request_reviewers (
  pull_request_id TEXT NOT NULL REFERENCES pull_requests (pull_request_id) ON DELETE CASCADE,
  reviewer_id TEXT NOT NULL REFERENCES users (user_id) ON DELETE RESTRICT,
  PRIMARY KEY (pull_request_id, reviewer_id)
);
