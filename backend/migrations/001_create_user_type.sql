CREATE TABLE user_type (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(10) CHECK (length(trim(title)) > 0)
);