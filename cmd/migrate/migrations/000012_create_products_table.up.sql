CREATE TABLE products(
id INT AUTO_INCREMENT,
roasted ENUM("light","medium","dark") DEFAULT "light",
price DECIMAL(10,2),
quantity INT,
bean_id INT NOT NULL,
form_id INT NOT NULL,
PRIMARY KEY(id),
FOREIGN KEY(bean_id) REFERENCES beans(id),
FOREIGN KEY(form_id) REFERENCES forms(id),
image VARCHAR(255)
);
