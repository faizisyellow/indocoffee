CREATE TRIGGER trg_cart_open_unique_insert
BEFORE INSERT ON cart_items
FOR EACH ROW
BEGIN
  IF NEW.status = 'open' AND
     EXISTS (
       SELECT 1
       FROM cart_items
       WHERE user_id = NEW.user_id
         AND product_id = NEW.product_id
         AND status = 'open'
     ) THEN
    SIGNAL SQLSTATE '45000'
      SET MESSAGE_TEXT = 'Duplicate open cart for this user/product (INSERT)';
  END IF;
END;

CREATE TRIGGER trg_cart_open_unique_update
BEFORE UPDATE ON cart_items
FOR EACH ROW
BEGIN
  IF NEW.status = 'open' AND
     EXISTS (
       SELECT 1
       FROM cart_items
       WHERE user_id = NEW.user_id
         AND product_id = NEW.product_id
         AND status = 'open'
         AND id <> NEW.id
     ) THEN
    SIGNAL SQLSTATE '45000'
      SET MESSAGE_TEXT = 'Duplicate open cart for this user/product (UPDATE)';
  END IF;
END;
