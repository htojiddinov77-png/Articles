-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS paragraphs(
    id BIGSERIAL PRIMARY KEY,
    headline VARCHAR(255) NOT NULL,
    body TEXT,
    order_index INT NOT NULL,
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE paragraphs;
-- +goose StatementEnd