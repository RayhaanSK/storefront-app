CREATE TABLE product (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    price MONEY NOT NULL,
    description TEXT NULL,
    listed_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    sold BOOLEAN DEFAULT FALSE NOT NULL,
    seller_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    CONSTRAINT fk_seller
    FOREIGN KEY (seller_id)
    REFERENCES user_account(id),
    CONSTRAINT fk_category
    FOREIGN KEY (category_id)
    REFERENCES category(id)
);  