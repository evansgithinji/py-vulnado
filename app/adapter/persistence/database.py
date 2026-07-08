import sqlite3
import os


def init_database(db_path: str) -> sqlite3.Connection:
    os.makedirs(os.path.dirname(db_path), exist_ok=True)

    conn = sqlite3.connect(db_path, check_same_thread=False)
    conn.row_factory = sqlite3.Row
    conn.execute("PRAGMA journal_mode=WAL")

    _create_tables(conn)
    _seed_data(conn)

    return conn


def _create_tables(conn: sqlite3.Connection):
    conn.executescript("""
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password TEXT,
            email TEXT,
            role TEXT DEFAULT 'user',
            is_admin INTEGER DEFAULT 0
        );

        CREATE TABLE IF NOT EXISTS products (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            price REAL NOT NULL,
            stock INTEGER DEFAULT 0,
            category TEXT
        );

        CREATE TABLE IF NOT EXISTS orders (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            product_id INTEGER NOT NULL,
            quantity INTEGER DEFAULT 1,
            total REAL NOT NULL,
            status TEXT DEFAULT 'pending',
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (product_id) REFERENCES products(id)
        );

        CREATE TABLE IF NOT EXISTS reviews (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            product_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            rating INTEGER NOT NULL,
            comment TEXT,
            FOREIGN KEY (product_id) REFERENCES products(id),
            FOREIGN KEY (user_id) REFERENCES users(id)
        );

        CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            content TEXT NOT NULL,
            author TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );
    """)


def _seed_data(conn: sqlite3.Connection):
    cursor = conn.cursor()

    # Users
    cursor.execute("SELECT COUNT(*) FROM users")
    if cursor.fetchone()[0] == 0:
        cursor.executemany(
            "INSERT INTO users (id, username, password, email, role, is_admin) VALUES (?, ?, ?, ?, ?, ?)",
            [
                (1, "admin", "admin123", "admin@vulnerable.app", "admin", 1),
                (2, "john", "password", "john@vulnerable.app", "user", 0),
                (3, "jane", "password123", "jane@vulnerable.app", "user", 0),
                (4, "guest", "guest", "guest@vulnerable.app", "guest", 0),
            ],
        )

    # Products
    cursor.execute("SELECT COUNT(*) FROM products")
    if cursor.fetchone()[0] == 0:
        cursor.executemany(
            "INSERT INTO products (id, name, price, stock, category) VALUES (?, ?, ?, ?, ?)",
            [
                (1, "Laptop Pro", 1299.99, 50, "Electronics"),
                (2, "Wireless Mouse", 29.99, 200, "Electronics"),
                (3, "USB-C Cable", 9.99, 500, "Accessories"),
                (4, "Mechanical Keyboard", 149.99, 75, "Electronics"),
                (5, "Monitor Stand", 49.99, 100, "Accessories"),
                (6, "Webcam HD", 79.99, 150, "Electronics"),
            ],
        )

    # Orders
    cursor.execute("SELECT COUNT(*) FROM orders")
    if cursor.fetchone()[0] == 0:
        cursor.executemany(
            "INSERT INTO orders (id, user_id, product_id, quantity, total, status) VALUES (?, ?, ?, ?, ?, ?)",
            [
                (1, 2, 1, 1, 1299.99, "completed"),
                (2, 2, 2, 2, 59.98, "pending"),
                (3, 3, 3, 5, 49.95, "shipped"),
                (4, 3, 4, 1, 149.99, "completed"),
                (5, 1, 6, 2, 159.98, "pending"),
            ],
        )

    # Reviews
    cursor.execute("SELECT COUNT(*) FROM reviews")
    if cursor.fetchone()[0] == 0:
        cursor.executemany(
            "INSERT INTO reviews (id, product_id, user_id, rating, comment) VALUES (?, ?, ?, ?, ?)",
            [
                (1, 1, 2, 5, "Excellent laptop, very fast!"),
                (2, 2, 3, 4, "Good mouse, comfortable grip"),
                (3, 4, 2, 5, "Best keyboard I have ever used"),
            ],
        )

    # Messages
    cursor.execute("SELECT COUNT(*) FROM messages")
    if cursor.fetchone()[0] == 0:
        cursor.executemany(
            "INSERT INTO messages (id, content, author) VALUES (?, ?, ?)",
            [
                (1, "Welcome to the message board!", "admin"),
                (2, "Hello everyone!", "john"),
            ],
        )

    conn.commit()
