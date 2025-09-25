ALTER TABLE users
ADD COLUMN role_id INT NOT NULL,
ADD CONSTRAINT fk_users_roles
FOREIGN KEY (role_id) REFERENCES roles(id);
