-- +goose Up
-- +goose StatementBegin

CREATE TABLE address (
    id SERIAL PRIMARY KEY,
    street TEXT NOT NULL,
    description TEXT,
    user_id INTEGER,
    CONSTRAINT fk_address_users FOREIGN KEY (user_id) REFERENCES users(id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS address;

-- +goose StatementEnd
