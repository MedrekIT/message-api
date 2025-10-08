-- +goose Up
CREATE TABLE invitation_links (
  token TEXT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  of_group_id UUID NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  CONSTRAINT fk_groups FOREIGN KEY (of_group_id) REFERENCES groups(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE invitation_links;
