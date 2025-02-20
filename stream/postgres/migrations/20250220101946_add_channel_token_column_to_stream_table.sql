-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.channels
ADD channel_token varchar(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.channels
DROP COLUMN channel_token;
-- +goose StatementEnd
