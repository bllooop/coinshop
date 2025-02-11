-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop
(
    id serial PRIMARY KEY,
    name varchar(150) NOT NULL,
    price int NOT NULL CHECK (price >= 0)
);
CREATE TABLE userlist
(
    id serial PRIMARY KEY,
    username varchar(255) NOT NULL unique,
    password varchar(255) NOT NULL,
    coins int NOT NULL CHECK (coins >= 0)
); 
CREATE TABLE transactions
(
    id serial PRIMARY KEY,
    source int NOT NULL,
    destination int NOT NULL,
    amount int NOT NULL CHECK (amount > 0),
    transaction_time TIMESTAMP DEFAULT now(), 
    FOREIGN KEY (source) REFERENCES userlist(id) ON DELETE CASCADE,
    FOREIGN KEY (destination) REFERENCES userlist(id) ON DELETE CASCADE
); 
CREATE TABLE purchases
(
    id serial PRIMARY KEY,
    user_id int NOT NULL,
    item_id int NOT NULL,
    price int NOT NULL,
    purchase_date TIMESTAMP DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES userlist(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES shop(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop;
DROP TABLE userlist;
DROP TABLE transactions;
DROP TABLE purchases;
-- +goose StatementEnd
