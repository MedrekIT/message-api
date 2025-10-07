-- +goose Up
CREATE TYPE relationship AS ENUM('pending', 'friends', 'blocked');
CREATE TABLE relations (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id UUID NOT NULL,
  receiver_id UUID NOT NULL,
  relationship RELATIONSHIP NOT NULL,
  CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_receiver FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE (user_id, receiver_id)
);

-- +goose Down
DROP TABLE relations;
DROP TYPE relationship;
