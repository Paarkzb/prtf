-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users_sessions
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rf_users_id     UUID REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    refresh_token   VARCHAR NOT NULL,
    expires_at      TIMESTAMP NOT NULL,
    deleted         BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.users_sessions;
-- +goose StatementEnd
