package database

import (
	"database/sql"
	"log"
)

// RunMigrations creates necessary database tables if they don't exist
func RunMigrations(db *sql.DB) error {
	// Create categories table
	createCategoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(createCategoriesTable)
	if err != nil {
		return err
	}
	log.Println("Categories table ready")

	// Create products table with foreign key to categories
	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price INTEGER NOT NULL DEFAULT 0,
		stock INTEGER NOT NULL DEFAULT 0,
		category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createProductsTable)
	if err != nil {
		return err
	}
	log.Println("Products table ready")

	// Create index on category_id for better JOIN performance
	createIndexQuery := `
	CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
	`

	_, err = db.Exec(createIndexQuery)
	if err != nil {
		return err
	}
	log.Println("Database indexes ready")

	return nil
}
