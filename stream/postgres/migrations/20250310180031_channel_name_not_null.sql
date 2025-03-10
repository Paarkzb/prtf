-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.channels
ALTER COLUMN channel_name SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.channels
ALTER COLUMN channel_name DROP NOT NULL;
-- +goose StatementEnd
