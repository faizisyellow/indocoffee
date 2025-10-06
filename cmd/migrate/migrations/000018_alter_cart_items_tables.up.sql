ALTER TABLE cart_items ADD COLUMN status ENUM("open","ordered") DEFAULT "open";
