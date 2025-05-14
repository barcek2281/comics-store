CREATE TABLE comics (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  author TEXT NOT NULL,
  description TEXT,
  release_date TEXT,
  price REAL,
  quantity INTEGER
);

CREATE TABLE order_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT, -- Internal item ID
  order_id TEXT NOT NULL, -- Foreign key to the order
  product_id TEXT NOT NULL, -- Product identifier
  quantity INTEGER NOT NULL, -- Quantity of the product
  FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE FOREIGN KEY (product_id) REFERENCES comics (id) ON DELETE CASCADE
);

CREATE TABLE orders (
  id TEXT PRIMARY KEY, -- Unique order ID (UUID)
  user_id TEXT NOT NULL, -- ID of the user who created the order
  total_price REAL NOT NULL, -- Total price of the order
  status TEXT NOT NULL, -- Order status (e.g., created, closed)
  created_at TEXT NOT NULL -- Timestamp when the order was created
);