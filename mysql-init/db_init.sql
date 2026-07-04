CREATE TABLE IF NOT EXISTS orders (
    id varchar(100) NOT NULL,
    price decimal(10,2),
    tax decimal(10,2),
    final_price decimal(10,2)
);