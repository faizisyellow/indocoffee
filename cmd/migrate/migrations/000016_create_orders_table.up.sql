CREATE TABLE orders(
id VARCHAR(255),
idempotency_key VARCHAR(255) UNIQUE NOT NULL,
customer_id INT NOT NULL,
customer_email VARCHAR(32) NOT NULL,
customer_name VARCHAR(24) NOT NULL,
status ENUM("confirm","roasting","shipped","complete","cancelled") DEFAULT "confirm",
items JSON NOT NULL,
total_price FLOAT NOT NULL,
phone_number INT NOT NULL,
alternative_phone_number VARCHAR(18),
street VARCHAR(16) NOT NULL,
city VARCHAR(16) NOT NULL,
PRIMARY KEY (id),
FOREIGN KEY (customer_id) REFERENCES users(id),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
