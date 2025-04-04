-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stream_id VARCHAR(255),
    username VARCHAR(255),
    text TEXT,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted boolean not null default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.messages;
-- +goose StatementEnd
