DROP TABLE IF EXISTS products;
CREATE TABLE IF NOT EXISTS products
(
    product_id bigserial not null constraint products_pkey primary key,
    name    text                        not null,
    price   numeric(10, 2) default 0.00 not null,
    created timestamp      default CURRENT_TIMESTAMP not null
);

alter table products owner to marius;
INSERT INTO products (name, price, created) SELECT CONCAT('Product-', MD5(random()::text)), random()*100, CURRENT_TIMESTAMP from generate_series(1, 1000000);

DROP TABLE IF EXISTS sales;
CREATE TABLE IF NOT EXISTS sales
(
    sale_id  bigserial not null constraint sales_pk primary key,
    product_id bigint,
    quantity   integer,
    created    timestamp default CURRENT_TIMESTAMP not null
);

alter table sales owner to marius;
INSERT INTO sales (product_id, quantity, created) SELECT random()*100,random()*100, CURRENT_TIMESTAMP from generate_series(1, 100);

