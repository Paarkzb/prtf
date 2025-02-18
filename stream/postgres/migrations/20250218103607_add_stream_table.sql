-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.streams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rf_channel_id UUID references public.channels(id) on delete cascade not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted boolean not null default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.streams
-- +goose StatementEnd
