CREATE TABLE address (id SERIAL PRIMARY KEY, street TEXT NOT NULL, user_id INTEGER, CONSTRAINT fk_address_users FOREIGN KEY (user_id) REFERENCES users(id));
