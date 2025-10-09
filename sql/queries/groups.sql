-- name: CreateGroup :one
INSERT INTO groups (id, created_at, updated_at, name, creator_id, group_type)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  $3,
  $4
)
RETURNING *;

-- name: GetPublicGroups :many
SELECT *, COUNT(users_groups.user_id) AS users_count
FROM groups INNER JOIN users_groups
ON groups.id = users_groups.group_id
WHERE groups.group_type = 'public' AND groups.name LIKE '%' || $1 || '%'
GROUP BY groups.id
ORDER BY users_count DESC;

-- name: GetGroupByID :one
SELECT * FROM groups
WHERE id = $1;

-- name: GetGroupByName :one
SELECT * FROM groups
WHERE name = $1;

-- name: RenameGroup :exec
UPDATE groups
SET name = $2, updated_at = NOW()
WHERE id = $1;

-- name: ChangeType :exec
UPDATE groups
SET group_type = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1;
