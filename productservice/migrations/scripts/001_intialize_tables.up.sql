CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,   -- Unique identifier for each product
    name VARCHAR(200) NOT NULL,      -- Product name
    description TEXT,                -- Product description
    price DECIMAL(10, 2) NOT NULL,   -- Price of the product
    inventory_count INT NOT NULL,    -- Available inventory for the product
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Product creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Timestamp for last update
);
