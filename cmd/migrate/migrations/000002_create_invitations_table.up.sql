CREATE TABLE
    invitations (
        user_id INT NOT NULL,
        token VARBINARY(72) NOT NULL,
        expire_at DATETIME NOT NULL
    )