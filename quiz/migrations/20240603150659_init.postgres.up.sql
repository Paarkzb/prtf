CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.user
(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(255) not null,
    username varchar(255) not null unique,
    password_hash varchar(255) not null,
    deleted boolean not null default false
);

CREATE TABLE IF NOT EXISTS public.quiz
(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    rf_user_id uuid references public.user(id) on delete cascade not null,
    name varchar(255) not null,
    description varchar(255),
    deleted boolean not null default false
);