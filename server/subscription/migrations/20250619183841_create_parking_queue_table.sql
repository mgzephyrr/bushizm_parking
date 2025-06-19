-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS parking_queue (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS parking_queue;
-- +goose StatementEnd
