-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.messages
ADD rf_stream_channel_id uuid not null;

ALTER TABLE public.messages
RENAME COLUMN stream_id TO rf_stream_id;

ALTER TABLE public.messages
ALTER COLUMN username TYPE uuid using username::uuid,
ALTER COLUMN username SET NOT NULL;


ALTER TABLE public.messages
RENAME COLUMN username TO rf_channel_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.messages
DROP COLUMN rf_stream_channel_id;

ALTER TABLE public.messages
RENAME COLUMN rf_stream_id TO stream_id;

ALTER TABLE public.messages
RENAME COLUMN rf_channel_id TO username;

ALTER TABLE public.messages
ALTER COLUMN username TYPE varchar(255),
ALTER COLUMN username DROP NOT NULL;


-- +goose StatementEnd
