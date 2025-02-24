-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.channels
ADD recording_path varchar(255),
ADD duration interval;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.channels
DROP COLUMN recording_path,
DROP COLUMN duration;
-- +goose StatementEnd
