
CREATE TABLE IF NOT EXISTS customers
(
    CustomerID INTEGER PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Surname VARCHAR(255) NOT NULL,
    Username VARCHAR(255) NOT NULL

);

create index customer_idx on customers (CustomerID);

CREATE TABLE IF NOT EXISTS products
(
    ProductID INTEGER PRIMARY KEY,
    Name VARCHAR(255) NOT NULL
);

create index product_idx on products (ProductID);

CREATE TABLE IF NOT EXISTS orders
(
    OrderID INTEGER PRIMARY KEY,
    CustomerID INTEGER ,
    ProductID INTEGER,
    quantity INTEGER,
    CONSTRAINT fk_customer Foreign KEY(CustomerID) References customers (CustomerID),
    CONSTRAINT fk_product Foreign KEY(ProductID) References products (ProductID)
);


