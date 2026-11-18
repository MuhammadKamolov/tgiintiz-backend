-- +goose Up
ALTER TABLE request_services
ADD COLUMN created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL;


-- +goose Down
ALTER TABLE request_services
DROP COLUMN created_at;