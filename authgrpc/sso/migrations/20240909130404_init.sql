-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.users
(
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username     VARCHAR(255) NOT NULL UNIQUE,
    email        VARCHAR(255) NOT NULL UNIQUE,
    pass_hash    VARCHAR(255)    NOT NULL,
    deleted      BOOLEAN NOT NULL DEFAULT false,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps
(
    id     UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name   TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS apps;
-- +goose StatementEnd
