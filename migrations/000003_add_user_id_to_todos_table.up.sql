alter table todos add column user_id uuid not null;

alter table todos add constraint fk_todos_user_id foreign key (user_id) references users(id) on delete cascade;