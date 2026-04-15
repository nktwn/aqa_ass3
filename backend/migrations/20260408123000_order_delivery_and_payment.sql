-- +goose Up
ALTER TABLE orders
    ADD COLUMN delivery_address TEXT NOT NULL DEFAULT '',
    ADD COLUMN delivery_comment TEXT NOT NULL DEFAULT '',
    ADD COLUMN payment_id TEXT,
    ADD COLUMN payment_status VARCHAR(50) NOT NULL DEFAULT 'pending';

-- +goose Down
ALTER TABLE orders
    DROP COLUMN IF EXISTS payment_status,
    DROP COLUMN IF EXISTS payment_id,
    DROP COLUMN IF EXISTS delivery_comment,
    DROP COLUMN IF EXISTS delivery_address;
