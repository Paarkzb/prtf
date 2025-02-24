-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.channels
ADD channel_name varchar(255) UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.channels
DROP COLUMN channel_name;
-- +goose StatementEnd
