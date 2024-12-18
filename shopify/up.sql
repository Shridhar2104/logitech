-- Create the database if it doesn't exist
CREATE DATABASE IF NOT EXISTS shopify_user_dev;

-- Switch to the `shopify_user_dev` database
\c shopify_user_dev;

-- Create `tokens` table if it doesn't exist
CREATE TABLE IF NOT EXISTS tokens (
    shop_name VARCHAR(255) NOT NULL,
    account_id VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    CONSTRAINT tokens_unique UNIQUE (shop_name, account_id)
);

-- Create `orders` table if it doesn't exist
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,  -- Auto-incrementing order id
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    shop_name VARCHAR(255) NOT NULL,
    account_id VARCHAR(255) NOT NULL,
    order_id VARCHAR(255) NOT NULL,
    total_price FLOAT NOT NULL,
    CONSTRAINT orders_shop_account_fk FOREIGN KEY (shop_name, account_id)
        REFERENCES tokens(shop_name, account_id)
);
