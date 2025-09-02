CREATE TABLE
    forms (
        id INT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(16) NOT NULL,
        is_delete BOOLEAN NOT NULL DEFAULT 0
    )