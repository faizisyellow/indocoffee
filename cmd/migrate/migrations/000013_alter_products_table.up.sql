ALTER TABLE products
ADD CONSTRAINT unique_product UNIQUE (roasted, bean_id, form_id);
