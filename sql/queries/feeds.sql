-- name: AddFeed :one
insert into feeds (name, url, user_id)
values (
        $1,
        $2,
        $3
       )
returning *;