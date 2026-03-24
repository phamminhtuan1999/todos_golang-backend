create table if not exists todos (
    id serial primary key,
    title varchar(255) not null,
    description text,
    completed boolean not null default false,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);