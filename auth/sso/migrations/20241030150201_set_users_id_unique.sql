-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.users_sessions ADD UNIQUE (rf_users_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.users_sessions DROP CONSTRAINT users_sessions_rf_users_id_key;
-- +goose StatementEnd
