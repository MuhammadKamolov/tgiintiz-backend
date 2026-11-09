-- +goose Up
CREATE TABLE request_attachments (
    id BIGSERIAL PRIMARY KEY,
    request_id BIGINT NOT NULL REFERENCES survey_requests(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_request_attachments_request_id ON request_attachments(request_id);

-- +goose Down
DROP TABLE IF EXISTS request_attachments;