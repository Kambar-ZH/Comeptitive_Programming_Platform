CREATE TABLE contests (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) DEFAULT '',
    start_date DATE DEFAULT NOW(),
    end_date DATE DEFAULT NOW(),
    description VARCHAR(255) DEFAULT ''
);