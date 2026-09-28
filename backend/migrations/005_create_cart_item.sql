CREATE TABLE cart_item (
    customer_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    PRIMARY KEY (customer_id, product_id),
    CONSTRAINT fk_customer
    FOREIGN KEY (customer_id)
    REFERENCES user_account(id),
    CONSTRAINT fk_product
    FOREIGN KEY (product_id)
    REFERENCES product(id) 
);