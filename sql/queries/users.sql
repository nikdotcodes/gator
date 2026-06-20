-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
           $1,
           $2,
           $3,
           $4
       )
RETURNING *;

-- name: GetUserByName :one
SELECT name
from users
where name = $1;

-- name: GetUserIDByName :one
select id
from users
where name = $1;

-- name: GetUsers :many
select name
from users;

-- name: DeleteAllUsers :exec
truncate table users;