ALTER TABLE orders
    MODIFY COLUMN items JSON NOT NULL COMMENT 'Array of OrderItem {id,image,bean_name,form_name,roasted,price,order_quantity}';
