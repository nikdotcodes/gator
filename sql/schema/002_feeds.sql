-- +goose Up
create table feeds (
    id uuid primary key ,
    name varchar(255) not null ,
    url varchar(4000) not null unique ,
    user_id uuid not null references users(id) on delete cascade ,
    created_at timestamp not null ,
    updated_at timestamp not null
);

-- +goose Down
drop table feeds;