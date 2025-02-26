-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.channels
ALTER COLUMN rf_active_stream_id SET DEFAULT uuid_nil();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.channels
ALTER COLUMN rf_active_stream_id SET DEFAULT null;
-- +goose StatementEnd
