CREATE TABLE customers (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  phone TEXT NOT NULL UNIQUE,
  timezone TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  price REAL NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_time_customers
  BEFORE UPDATE 
  ON customers
  FOR EACH ROW
BEGIN 
  UPDATE customers SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; 
END;

CREATE TRIGGER update_time_products
  BEFORE UPDATE 
  ON products
  FOR EACH ROW
BEGIN 
  UPDATE products SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; 
END;