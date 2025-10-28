-- +goose Up
-- +goose StatementBegin

CREATE TABLE articles IF NOT EXISTS {
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image VARCHAR(255),
    author_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
}
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE articles;
-- +goose StatementEnd