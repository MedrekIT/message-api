-- name: CreateFriendship :one
INSERT INTO relations (id, created_at, updated_at, user_id, receiver_id, relationship)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  $3,
  'pending'
)
RETURNING *;

-- name: CreateBlock :one
INSERT INTO relations (id, created_at, updated_at, user_id, receiver_id, relationship)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  $3,
  'blocked'
)
RETURNING *;

-- name: AcceptFriendship :exec
UPDATE relations
SET relationship = 'friends', updated_at = NOW()
WHERE receiver_id = $1 AND user_id = $2 AND relationship = 'pending';

-- name: DeclineFriendship :exec
DELETE FROM relations
WHERE receiver_id = $1 AND user_id = $2 AND relationship = 'pending';

-- name: GetBlocks :many
SELECT *, users.login AS user_login
FROM relations INNER JOIN users
ON relations.receiver_id = users.id
WHERE (relations.user_id = $1) AND relations.relationship = 'blocked';

-- name: GetFriends :many
SELECT *, users.login AS user_login
FROM relations INNER JOIN users
ON (relations.user_id = users.id AND relations.receiver_id = $1)
OR (relations.receiver_id = users.id AND relations.user_id = $1)
WHERE relations.relationship = 'friends';

-- name: DeleteFriend :one
DELETE FROM relations
WHERE (user_id = $1 OR receiver_id = $1)
AND (user_id = $2 OR receiver_id = $2)
AND relationship = 'friends'
RETURNING *;
