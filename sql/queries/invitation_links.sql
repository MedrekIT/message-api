-- name: CreateInvitationLink :one
INSERT INTO invitation_links (token, created_at, updated_at, of_group_id, expires_at)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  $3
)
RETURNING *;

-- name: ClearInvitationLinks :exec
DELETE FROM invitation_links
WHERE expires_at < NOW();
