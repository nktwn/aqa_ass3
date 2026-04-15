-- +goose Up
-- SQL в секции Up выполняется при применении миграции
------------------------------------------------------------------------
/* 1. Справочные таблицы */
CREATE TABLE delivery_conditions (
    condition_id           SERIAL PRIMARY KEY,
    minimum_free_delivery_amount NUMERIC,
    delivery_fee           NUMERIC
);

CREATE TABLE users (
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(50) NOT NULL,
    phone_number   VARCHAR(12) NOT NULL UNIQUE,
    hashed_password TEXT       NOT NULL,
    role           INTEGER     NOT NULL DEFAULT 0,
    created_at     TIMESTAMP   NOT NULL DEFAULT now(),
    updated_at     TIMESTAMP   NOT NULL DEFAULT now()
);

/* 2. Поставщики */
CREATE TABLE suppliers (
    user_id      INTEGER     PRIMARY KEY,
    condition_id INTEGER,
    order_amount INTEGER     NOT NULL,
    name         VARCHAR(100),
    FOREIGN KEY (user_id)      REFERENCES users(id),
    FOREIGN KEY (condition_id) REFERENCES delivery_conditions(condition_id)
);

/* 3. Товары и цены */
CREATE TABLE products (
    id                 SERIAL PRIMARY KEY,
    name               TEXT     NOT NULL,
    image_url          TEXT,
    gtin               BIGINT,
    lowest_price       INTEGER,
    created_at         TIMESTAMP NOT NULL DEFAULT now(),
    updated_at         TIMESTAMP NOT NULL DEFAULT now(),
    lowest_supplier_id INTEGER,
    FOREIGN KEY (lowest_supplier_id) REFERENCES suppliers(user_id)
);

CREATE TABLE products_supplier (
    product_id  INTEGER NOT NULL,
    supplier_id INTEGER NOT NULL,
    price       INTEGER NOT NULL,
    sell_amount INTEGER NOT NULL,
    PRIMARY KEY (product_id, supplier_id),
    FOREIGN KEY (product_id)  REFERENCES products(id),
    FOREIGN KEY (supplier_id) REFERENCES suppliers(user_id)
);

/* 4. Корзины */
CREATE TABLE carts (
    id          SERIAL PRIMARY KEY,
    customer_id INTEGER,
    total       INTEGER   NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (customer_id) REFERENCES users(id)
);

CREATE TABLE cart_items (
    cart_id     INTEGER NOT NULL,
    supplier_id INTEGER NOT NULL,
    product_id  INTEGER NOT NULL,
    quantity    INTEGER NOT NULL,
    price       INTEGER NOT NULL,
    PRIMARY KEY (cart_id, supplier_id, product_id),
    FOREIGN KEY (cart_id)     REFERENCES carts(id)          ON DELETE CASCADE,
    FOREIGN KEY (product_id)  REFERENCES products(id),
    FOREIGN KEY (supplier_id) REFERENCES suppliers(user_id)
);

/* 5. Заказы */
CREATE TABLE order_status (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Sample data for order statuses
INSERT INTO order_status (name) VALUES
('Pending'),
('In Progress'),
('Completed'),
('Cancelled');

CREATE TABLE orders (
    id          SERIAL PRIMARY KEY,
    order_date  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status_id   INTEGER     NOT NULL DEFAULT 1,
    customer_id INTEGER,
    supplier_id INTEGER,
    FOREIGN KEY (customer_id) REFERENCES users(id),
    FOREIGN KEY (supplier_id) REFERENCES suppliers(user_id),
    FOREIGN KEY (status_id) REFERENCES order_status(id)
);

CREATE INDEX idx_orders_supplier_id ON orders(supplier_id);

CREATE TABLE order_products (
    order_id   INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity   INTEGER NOT NULL,
    price      INTEGER NOT NULL,
    PRIMARY KEY (order_id, product_id),
    FOREIGN KEY (order_id)   REFERENCES orders(id)   ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE contracts (
    id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    supplier_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    supplier_sig TEXT,
    customer_sig TEXT,
    status SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    signed_at TIMESTAMP
);

CREATE INDEX idx_order_products_order_id  ON order_products(order_id);
CREATE INDEX idx_order_products_product_id ON order_products(product_id);

------------------------------------------------------------------------
-- +goose Down
-- SQL в секции Down выполняется при откате миграции
------------------------------------------------------------------------
-- Снимаем зависимости в порядке «от ребёнка к родителю»
DROP INDEX IF EXISTS idx_order_products_product_id;
DROP INDEX IF EXISTS idx_order_products_order_id;
DROP TABLE IF EXISTS order_products;

DROP INDEX IF EXISTS idx_orders_supplier_id;
DROP TABLE IF EXISTS orders;

DROP TABLE IF EXISTS contracts;

DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;

DROP TABLE IF EXISTS products_supplier;
DROP TABLE IF EXISTS products;

DROP TABLE IF EXISTS suppliers;
DROP TABLE IF EXISTS delivery_conditions;
DROP TABLE IF EXISTS users;

