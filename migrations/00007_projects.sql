-- +goose Up
CREATE TABLE projects(
    id BIGSERIAL PRIMARY KEY,
    lang VARCHAR(5) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    location VARCHAR(255),
    completed_year SMALLINT,
    short_desc VARCHAR(255) NOT NULL,
    full_desc TEXT,
    image_url VARCHAR(255),
    images TEXT[] NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    client VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    UNIQUE(slug, lang)
);

CREATE TABLE project_services(
    project_id BIGINT NOT NULL,
    service_id BIGINT NOT NULL,
    PRIMARY KEY(project_id, service_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY(service_id) REFERENCES services(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE INDEX idx_project_services_service_id ON project_services(service_id);

-- +goose Down
DROP TABLE IF NOT EXISTS project_services;
DROP TABLE IF NOT EXISTS projects;