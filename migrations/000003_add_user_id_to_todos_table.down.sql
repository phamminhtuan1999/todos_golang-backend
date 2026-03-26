alter table todos drop constraint if exists fk_todos_user_id;

alter table todos drop column if exists user_id;