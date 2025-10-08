-- +goose Up
CREATE TYPE group_type AS ENUM('public', 'invite_only', 'private');
CREATE TABLE groups (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT UNIQUE NOT NULL,
  creator_id UUID UNIQUE NOT NULL,
  group_type GROUP_TYPE NOT NULL,
  CONSTRAINT fk_users FOREIGN KEY (creator_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE users;
DROP TYPE group_type;
