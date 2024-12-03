-- +goose Up
-- +goose StatementBegin
INSERT INTO apps (id, name, secret)  
VALUES ('36c604ca-5f22-447c-a2a7-f220d2c1193b', 'test', 'test-secret') 
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM apps WHERE id 
-- +goose StatementEnd
