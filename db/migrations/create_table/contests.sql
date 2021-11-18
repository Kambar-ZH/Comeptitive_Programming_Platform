CREATE TABLE contests (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) DEFAULT '',
    start_date DATE DEFAULT NOW(),
    description VARCHAR(255) DEFAULT '',
    author_handle VARCHAR(255) DEFAULT ''
);