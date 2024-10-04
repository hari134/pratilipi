CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,         -- Unique identifier for the user
    name VARCHAR(100) NOT NULL,         -- User's name
    phone_no VARCHAR(100) NOT NULL,         -- User's name
    email VARCHAR(150) UNIQUE NOT NULL, -- Unique email address
    password_hash VARCHAR(255) NOT NULL,-- Hashed password for authentication
    role VARCHAR(50) DEFAULT 'user',    -- Role: 'admin' or 'user', defaults to 'user'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- User registration timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp for last update
    CHECK (role IN ('admin', 'user'))   -- Ensure only 'admin' or 'user' roles are allowed
);
