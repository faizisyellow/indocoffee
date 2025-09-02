CREATE TABLE
    users (
        id INT PRIMARY KEY AUTO_INCREMENT,
        username VARCHAR(16) NOT NULL,
        email VARCHAR(32) NOT NULL UNIQUE,
        password VARBINARY(72) NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT 0,
        updated_at TIMESTAMP,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        profile VARCHAR(32),
        is_delete BOOLEAN NOT NULL DEFAULT 0
    )