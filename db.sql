create table person
(
    id         serial
        primary key,
    content    jsonb,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp default now(),
    cache      boolean
)