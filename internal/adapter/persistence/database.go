package persistence

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// InitDatabase initializes the SQLite database
func InitDatabase(dbPath string) (*sql.DB, error) {
	// Ensure directory exists
	if err := os.MkdirAll("/app/data", 0755); err != nil {
		// Try local path if /app/data doesn't work
		os.MkdirAll("./data", 0755)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	// Seed initial data
	if err := seedData(db); err != nil {
		return nil, fmt.Errorf("failed to seed data: %v", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			email TEXT,
			role TEXT DEFAULT 'user',
			is_admin INTEGER DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			price REAL NOT NULL,
			stock INTEGER DEFAULT 0,
			category TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER DEFAULT 1,
			total REAL NOT NULL,
			status TEXT DEFAULT 'pending',
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`,
		`CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			product_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			rating INTEGER NOT NULL,
			comment TEXT,
			FOREIGN KEY (product_id) REFERENCES products(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			author TEXT,
			created_at TEXT
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func seedData(db *sql.DB) error {
	// Seed users (with weak passwords for demo)
	users := []string{
		`INSERT OR IGNORE INTO users (id, username, password, email, role, is_admin) VALUES (1, 'admin', 'admin123', 'admin@vulnerable.app', 'admin', 1)`,
		`INSERT OR IGNORE INTO users (id, username, password, email, role, is_admin) VALUES (2, 'john', 'password', 'john@vulnerable.app', 'user', 0)`,
		`INSERT OR IGNORE INTO users (id, username, password, email, role, is_admin) VALUES (3, 'jane', 'password123', 'jane@vulnerable.app', 'user', 0)`,
		`INSERT OR IGNORE INTO users (id, username, password, email, role, is_admin) VALUES (4, 'guest', 'guest', 'guest@vulnerable.app', 'guest', 0)`,
	}

	// Seed products
	products := []string{
		`INSERT OR IGNORE INTO products (id, name, price, stock, category) VALUES (1, 'Laptop Pro', 1299.99, 50, 'Electronics')`,
		`INSERT OR IGNORE INTO products (id, name, price, stock, category) VALUES (2, 'Wireless Mouse', 29.99, 200, 'Electronics')`,
		`INSERT OR IGNORE INTO products (id, name, price, stock, category) VALUES (3, 'USB-C Cable', 9.99, 500, 'Accessories')`,
		`INSERT OR IGNORE INTO products (id, name, price, stock, category) VALUES (4, 'Mechanical Keyboard', 149.99, 75, 'Electronics')`,
		`INSERT OR IGNORE INTO products (id, name, price, stock, category) VALUES (5, 'Monitor Stand', 49.99, 100, 'Accessories')`,
		`INSERT OR IGNORE INTO products (id, name, price, stock, category) VALUES (6, 'Webcam HD', 79.99, 150, 'Electronics')`,
	}

	// Seed orders
	orders := []string{
		`INSERT OR IGNORE INTO orders (id, user_id, product_id, quantity, total, status) VALUES (1, 2, 1, 1, 1299.99, 'completed')`,
		`INSERT OR IGNORE INTO orders (id, user_id, product_id, quantity, total, status) VALUES (2, 2, 2, 2, 59.98, 'pending')`,
		`INSERT OR IGNORE INTO orders (id, user_id, product_id, quantity, total, status) VALUES (3, 3, 4, 1, 149.99, 'shipped')`,
		`INSERT OR IGNORE INTO orders (id, user_id, product_id, quantity, total, status) VALUES (4, 3, 3, 5, 49.95, 'completed')`,
		`INSERT OR IGNORE INTO orders (id, user_id, product_id, quantity, total, status) VALUES (5, 1, 6, 3, 239.97, 'pending')`,
	}

	// Seed reviews
	reviews := []string{
		`INSERT OR IGNORE INTO reviews (id, product_id, user_id, rating, comment) VALUES (1, 1, 2, 5, 'Excellent laptop! Fast and reliable.')`,
		`INSERT OR IGNORE INTO reviews (id, product_id, user_id, rating, comment) VALUES (2, 2, 3, 4, 'Good mouse, comfortable to use.')`,
		`INSERT OR IGNORE INTO reviews (id, product_id, user_id, rating, comment) VALUES (3, 4, 2, 5, 'Best keyboard I have ever used!')`,
	}

	// Seed messages
	messages := []string{
		`INSERT OR IGNORE INTO messages (id, content, author, created_at) VALUES (1, 'Welcome to our message board!', 'admin', '2024-01-01 10:00:00')`,
		`INSERT OR IGNORE INTO messages (id, content, author, created_at) VALUES (2, 'Hello everyone!', 'john', '2024-01-02 14:30:00')`,
	}

	allSeeds := append(users, products...)
	allSeeds = append(allSeeds, orders...)
	allSeeds = append(allSeeds, reviews...)
	allSeeds = append(allSeeds, messages...)

	for _, query := range allSeeds {
		if _, err := db.Exec(query); err != nil {
			// Ignore duplicate key errors
			continue
		}
	}

	return nil
}
