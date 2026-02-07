package repositories

import (
	"category-management-api/models"
	"database/sql"
	"fmt"
)

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error)
	GetDailySalesReport() (*models.SalesReport, error)
	GetSalesReportByDateRange(startDate, endDate string) (*models.SalesReport, error)
}

// transactionRepository implements TransactionRepository interface
type transactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository instance
func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// CreateTransaction processes a checkout: validates products, deducts stock,
// creates transaction record and detail rows inside a single DB transaction.
func (repo *transactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0, len(items))

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow(
			"SELECT name, price, stock FROM products WHERE id = $1",
			item.ProductID,
		).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product '%s' (available: %d, requested: %d)",
				productName, stock, item.Quantity)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec(
			"UPDATE products SET stock = stock - $1 WHERE id = $2",
			item.Quantity, item.ProductID,
		)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Insert transaction header
	var transactionID int
	err = tx.QueryRow(
		"INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id",
		totalAmount,
	).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// Insert transaction details — use RETURNING id to capture the generated ID
	for i := range details {
		details[i].TransactionID = transactionID

		var detailID int
		err = tx.QueryRow(
			"INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal,
		).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		details[i].ID = detailID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

// GetDailySalesReport returns the sales summary for today
func (repo *transactionRepository) GetDailySalesReport() (*models.SalesReport, error) {
	report := &models.SalesReport{}

	// Get total revenue and transaction count for today
	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE created_at::date = CURRENT_DATE
	`).Scan(&report.TotalRevenue, &report.TotalTransactions)
	if err != nil {
		return nil, err
	}

	// Get best selling product for today
	var best models.BestSellingProduct
	err = repo.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) AS qty_sold
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at::date = CURRENT_DATE
		GROUP BY p.id, p.name
		ORDER BY qty_sold DESC
		LIMIT 1
	`).Scan(&best.Name, &best.QtySold)
	if err == sql.ErrNoRows {
		report.BestSellingProduct = nil
	} else if err != nil {
		return nil, err
	} else {
		report.BestSellingProduct = &best
	}

	return report, nil
}

// GetSalesReportByDateRange returns the sales summary for a given date range
func (repo *transactionRepository) GetSalesReportByDateRange(startDate, endDate string) (*models.SalesReport, error) {
	report := &models.SalesReport{}

	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE created_at::date >= $1::date AND created_at::date <= $2::date
	`, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransactions)
	if err != nil {
		return nil, err
	}

	var best models.BestSellingProduct
	err = repo.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) AS qty_sold
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at::date >= $1::date AND t.created_at::date <= $2::date
		GROUP BY p.id, p.name
		ORDER BY qty_sold DESC
		LIMIT 1
	`, startDate, endDate).Scan(&best.Name, &best.QtySold)
	if err == sql.ErrNoRows {
		report.BestSellingProduct = nil
	} else if err != nil {
		return nil, err
	} else {
		report.BestSellingProduct = &best
	}

	return report, nil
}
