-- +goose Up
CREATE TABLE feeds (
  id integer PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT,
  url TEXT UNIQUE NOT NULL,
  user_id  UUID,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
