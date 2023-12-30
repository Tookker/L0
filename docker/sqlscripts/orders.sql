BEGIN;

CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR(100) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR(100) NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR(100) NOT NULL,
    status INTEGER NOT NULL
);

INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES (134621, 'WBILMTESTTRACK', 453, 'ab4219087a764ae0btest', 'Mascaras', 30, '0', 317, 2389212, 'Vivienne Sabo', 202);

INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES (54543543, 'WBILMTESTTRACK', 10000, 'ab4219vf34a764ae0btest', 'MOBILE PHONE', 0, '10', 10000, 5189212, 'ReadMi', 202);

INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES (15443523, 'WBILMTESTTRACK', 4230, 'ab4221387a764ae0btest', 'Booster', 0, '50', 4230, 1689212, 'Titan', 202);

CREATE TABLE IF NOT EXISTS deliverys (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(100) NOT NULL,
    zip VARCHAR(10) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address VARCHAR(100) NOT NULL,
    region VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL
);

INSERT INTO deliverys (name, phone, zip, city, address, region, email)
VALUES ('Test Testov', '+9720000000', '2639809', 'Kiryat Mozkin', 'Ploshad Mira 15', 'Kraiot', 'test@gmail.com');

CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    transaction VARCHAR(100) NOT NULL,
    request_id VARCHAR(100),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(100) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL
);

INSERT INTO payments (transaction, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ('b563feb7b2b84b6test', 'USD', 'wbpay', 1817, 1637907727, 'alpha', 1500, 317, 0);

INSERT INTO payments (transaction, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ('b56fgveb7b2b84b6test', 'Rubles', 'sber', 1991, 16379075427, 'sber', 15120, 317, 0);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_uid VARCHAR(100) NOT NULL,
    track_number VARCHAR(100) NOT NULL,
    entry VARCHAR(100) NOT NULL,
    id_delivery BIGINT NOT NULL,
    id_payment BIGINT NOT NULL,
    locale VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(100),
    customer_id VARCHAR(100) NOT NULL,
    delivery_service VARCHAR(100) NOT NULL,
    shardkey VARCHAR(100) NOT NULL,
    sm_id INTEGER,
    date_created DATE NOT NULL,
    oof_shard VARCHAR(100) NOT NULL
);

ALTER TABLE IF EXISTS orders
    ADD FOREIGN KEY (id_delivery)
    REFERENCES deliverys (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE IF EXISTS orders
    ADD FOREIGN KEY (id_payment)
    REFERENCES payments (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

INSERT INTO orders (order_uid, track_number, entry, id_delivery, id_payment, 
                    locale,  customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ('b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', 1, 1, 'en', 'test', 'meest', '9', 99, '2021-11-26T06:22:19Z', '1');

INSERT INTO orders (order_uid, track_number, entry, id_delivery, id_payment, 
                    locale,  customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ('gfad63feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', 1, 2, 'en', 'root', 'meest', '9', 99, '2021-12-26T06:22:19Z', '2');

CREATE TABLE IF NOT EXISTS orders_items(
    id BIGSERIAL PRIMARY KEY,
    id_order BIGINT NOT NULL,
    id_item BIGINT NOT NULL
);

ALTER TABLE IF EXISTS orders_items
    ADD FOREIGN KEY (id_order)
    REFERENCES orders (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

ALTER TABLE IF EXISTS orders_items
    ADD FOREIGN KEY (id_item)
    REFERENCES items (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

INSERT INTO orders_items (id_order, id_item)
VALUES (1, 1);

INSERT INTO orders_items (id_order, id_item)
VALUES (2, 2);

INSERT INTO orders_items (id_order, id_item)
VALUES (2, 3);


END;