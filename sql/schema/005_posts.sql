-- +goose Up
CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  title TEXT NOT NULL,
  url TEXT UNIQUE NOT NULL,
  description TEXT,
  published_at TIMESTAMP,
  feed_id UUID NOT NULL REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;
