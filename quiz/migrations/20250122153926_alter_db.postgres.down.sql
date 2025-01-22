CREATE TABLE IF NOT EXISTS public.user
(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(255) not null,
    username varchar(255) not null unique,
    password_hash varchar(255) not null,
    deleted boolean not null default false
);

ALTER TABLE IF EXISTS public.quiz ADD CONSTRAINT quiz_rf_user_id_fkey FOREIGN KEY (rf_user_id) references public.user(id) on delete cascade;