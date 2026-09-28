CREATE TABLE user_account (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(30) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    user_type_id BIGINT NOT NULL,
    CONSTRAINT fk_user_type
    FOREIGN KEY (user_type_id)
    REFERENCES user_type(id)
);