-- +goose Up
INSERT INTO users (id, name, phone_number, hashed_password, role)
VALUES (
    999,
    'Admin',
    '+77000000999',
    '$2a$10$MqiMKXRYxt6R8boUra6E1.n3V4.JT61wRVW7sE8xZGPIYdF/lOOr6',
    2
)
ON CONFLICT (id) DO UPDATE
SET
    name = EXCLUDED.name,
    phone_number = EXCLUDED.phone_number,
    hashed_password = EXCLUDED.hashed_password,
    role = EXCLUDED.role;

-- +goose Down
DELETE FROM users WHERE id = 999;
