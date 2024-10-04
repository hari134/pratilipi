CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,            -- Unique identifier for each order
    user_id INT NOT NULL,                   -- Reference to the user who placed the order
    total_price DECIMAL(10, 2) NOT NULL,    -- Total price of the order
    status VARCHAR(50) NOT NULL,            -- Order status: placed, shipped, completed, etc.
    placed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the order was placed
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Timestamp for the last update
);


--bun:split

CREATE TABLE order_items (
    order_item_id SERIAL PRIMARY KEY,       -- Unique identifier for each order item
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE, -- Reference to the order
    product_id INT NOT NULL,                -- Reference to the product
    quantity INT NOT NULL,                  -- Quantity of the product ordered
    price_at_order DECIMAL(10, 2) NOT NULL  -- Product price at the time of the order
);


--bun:split

CREATE TABLE users (
    user_id INT PRIMARY KEY,                -- User ID (received from the User Service)
    email VARCHAR(255) NOT NULL,            -- Email of the user (received from the User Service)
    phone_no VARCHAR(20),                   -- Phone number of the user (optional, received from the User Service)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Timestamp when the user was added
);


--bun:split

CREATE TABLE products (
    product_id INT PRIMARY KEY,             -- Product ID (received from the Product Service)
    price DECIMAL(10, 2) NOT NULL,          -- Product price
    inventory_count INT NOT NULL            -- Inventory count for the product
);



