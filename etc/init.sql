CREATE TABLE IF NOT EXISTS products
(
    id bigserial constraint products_pk primary key,
    name varchar(250) not null,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    created timestamp,
    modified timestamp
);

CREATE INDEX IF NOT EXISTS products_name_index ON products (name);

INSERT INTO public.products (id, name, price, created, modified) VALUES (1, 'Test product', 11.01, '2020-12-01 01:16:26.000000', '2020-12-01 01:16:31.000000');