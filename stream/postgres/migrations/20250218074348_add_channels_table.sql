-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.channels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rf_user_id UUID not null unique,
    live boolean default false not null,
    rf_active_stream_id UUID default null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted boolean not null default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.channels
-- +goose StatementEnd
