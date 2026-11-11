-- +goose Up
CREATE TABLE services(
    id BIGSERIAL PRIMARY KEY,
    lang VARCHAR(5) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    short_desc VARCHAR(255) NOT NULL,
    full_desc TEXT,
    image_url VARCHAR(255),
    standards JSONB NOT NULL DEFAULT '[]'::jsonb,
    equipment JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    UNIQUE(slug, lang)
);

CREATE TABLE request_services(
    request_id BIGINT NOT NULL, 
    service_id BIGINT NOT NULL, 
    PRIMARY KEY (request_id, service_id),
    FOREIGN KEY(request_id) REFERENCES survey_requests(id) ON DELETE CASCADE,
    FOREIGN KEY(service_id) REFERENCES services(id) ON DELETE RESTRICT
);

CREATE INDEX idx_request_services_service_id ON request_services(service_id);

-- +goose Down
DROP TABLE IF NOT EXISTS request_services;
DROP TABLE IF NOT EXISTS services;