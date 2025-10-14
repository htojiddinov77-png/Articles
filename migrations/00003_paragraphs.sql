-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS paraghraps (
    id BIGSERIAL PRIMARY KEY,
    headline VARCHAR(255) NOT NULL,
    body TEXT,
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    order_index INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE paraghraps;
-- +goose StatementEnd