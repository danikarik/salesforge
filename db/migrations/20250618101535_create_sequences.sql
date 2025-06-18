-- +goose Up
-- +goose StatementBegin
CREATE TABLE sequences (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    open_tracking_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    click_tracking_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sequences;
-- +goose StatementEnd
