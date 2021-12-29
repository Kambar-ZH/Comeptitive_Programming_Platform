CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    handle VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    country VARCHAR(255) DEFAULT '',
    city VARCHAR(255) DEFAULT '',
    rating INTEGER DEFAULT 0,
    max_rating INTEGER DEFAULT 0,
    avatar VARCHAR(255) DEFAULT '',
    encrypted_password VARCHAR(255) NOT NULL
);