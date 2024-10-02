CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,                   -- Unique identifier for the order
    user_id VARCHAR(255) REFERENCES users(user_id) ON DELETE SET NULL,  -- Foreign key to local users table
    total_price DECIMAL(10, 2) NOT NULL,           -- Total price of the order
    status VARCHAR(50) NOT NULL,                   -- Order status: placed, shipped, completed, etc.
    placed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the order was placed
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Timestamp for last update
);


CREATE TABLE order_items (
    order_item_id SERIAL PRIMARY KEY,              -- Unique identifier for each item in an order
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE, -- Reference to the order, cascade delete
    product_id INT NOT NULL,                       -- Reference to the product (no foreign key to Product Service)
    quantity INT NOT NULL,                         -- Quantity of the product ordered
    price_at_order DECIMAL(10, 2) NOT NULL,        -- Product price at the time of the order
    stock_at_order INT NOT NULL,                   -- Inventory level when the order was placed (for reference)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the order item was added
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Timestamp for last update
);


CREATE TABLE users (
    user_id VARCHAR(255) PRIMARY KEY,  -- Unique identifier for the user (received from User Service)
    email VARCHAR(255) NOT NULL,       -- Email of the user (received from User Service)
    phone_no VARCHAR(20),              -- Phone number of the user (optional, received from User Service)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the user was added to the Order Service
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Timestamp for last update
);


CREATE TABLE products (
    product_id VARCHAR(255) PRIMARY KEY,  -- Unique identifier for the product (received from Product Service)
    name VARCHAR(255) NOT NULL,           -- Product name
    price DECIMAL(10, 2) NOT NULL,        -- Product price at the time of creation
    stock INT NOT NULL DEFAULT 0,         -- Current stock level of the product
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the product was added to the Order Service
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Timestamp for last update
);



