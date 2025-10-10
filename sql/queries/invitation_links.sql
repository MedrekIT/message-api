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

-- name: GetInvitation :one
SELECT * FROM invitation_links
WHERE token = $1 AND expires_at > NOW();

-- name: ClearInvitationLinks :exec
DELETE FROM invitation_links
WHERE expires_at < NOW();
