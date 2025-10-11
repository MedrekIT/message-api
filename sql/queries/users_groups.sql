-- name: AddMember :one
INSERT INTO users_groups (created_at, updated_at, user_id, of_group_id)
VALUES (
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;

-- name: ChangePermissions :exec
UPDATE users_groups
SET member_type = $3, updated_at = NOW()
WHERE user_id = $1 AND of_group_id = $2;

-- name: RemoveMember :exec
DELETE FROM users_groups
WHERE user_id = $1 AND of_group_id = $2;

-- name: GetBans :many
SELECT *, users.login AS user_login
FROM users_groups INNER JOIN users
ON users_groups.user_id = users.id
WHERE users_groups.of_group_id = $1 AND member_type = 'blocked';

-- name: GetMembers :many
SELECT *, users.login AS user_login
FROM users_groups INNER JOIN users
ON users_groups.user_id = users.id
WHERE users_groups.of_group_id = $1 AND member_type != 'blocked';

-- name: GetMember :one
SELECT * FROM users_groups
WHERE user_id = $1 AND of_group_id = $2 AND member_type != 'blocked';
