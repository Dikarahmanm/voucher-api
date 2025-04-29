-- +goose Up
CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE vouchers (
    id SERIAL PRIMARY KEY,
    brand_id INTEGER REFERENCES brands (id),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    cost_in_points INTEGER NOT NULL,
    expiry_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    points_balance INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers (id),
    total_points INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transaction_vouchers (
    transaction_id INTEGER REFERENCES transactions (id),
    voucher_id INTEGER REFERENCES vouchers (id),
    quantity INTEGER NOT NULL,
    PRIMARY KEY (transaction_id, voucher_id)
);

-- +goose Down
DROP TABLE transaction_vouchers;

DROP TABLE transactions;

DROP TABLE customers;

DROP TABLE vouchers;

DROP TABLE brands;