CREATE TABLE category (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL CHECK (length(trim(title)) > 0)
);