ALTER TABLE products
ADD CONSTRAINT quantity_non_negative CHECK (quantity >= 0);
