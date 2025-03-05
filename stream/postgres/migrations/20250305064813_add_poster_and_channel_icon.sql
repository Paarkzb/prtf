-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.streams
ADD poster varchar(255);
ALTER TABLE public.channels
ADD icon varchar(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.streams
DROP COLUMN poster;
ALTER TABLE public.channels
DROP COLUMN icon;
-- +goose StatementEnd
