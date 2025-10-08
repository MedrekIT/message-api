-- +goose Up
CREATE TYPE member_type AS ENUM('member', 'moderator', 'admin', 'blocked');
CREATE TABLE users_groups (
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id UUID NOT NULL,
  of_group_id UUID NOT NULL,
  member_type MEMBER_TYPE NOT NULL DEFAULT 'member',
  CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_groups FOREIGN KEY (of_group_id) REFERENCES groups(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, of_group_id)
);

-- +goose Down
DROP TABLE users_groups;
DROP TYPE member_type;
