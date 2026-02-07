package models

import "time"

// Transaction represents a completed transaction
// @Description Transaction information with details of purchased items
type Transaction struct {
	ID          int                 `json:"id" example:"1"`
	TotalAmount int                 `json:"total_amount" example:"45000"`
	CreatedAt   time.Time           `json:"created_at" example:"2026-02-08T12:00:00Z"`
	Details     []TransactionDetail `json:"details"`
}

// TransactionDetail represents a single item in a transaction
// @Description Detail of a single item within a transaction
type TransactionDetail struct {
	ID            int    `json:"id" example:"1"`
	TransactionID int    `json:"transaction_id" example:"1"`
	ProductID     int    `json:"product_id" example:"3"`
	ProductName   string `json:"product_name,omitempty" example:"Indomie Goreng"`
	Quantity      int    `json:"quantity" example:"5"`
	Subtotal      int    `json:"subtotal" example:"15000"`
}

// CheckoutItem represents a single item in a checkout request
// @Description Single item to be checked out
type CheckoutItem struct {
	ProductID int `json:"product_id" example:"3"`
	Quantity  int `json:"quantity" example:"5"`
}

// CheckoutRequest represents the request body for checkout
// @Description Request body for processing a checkout
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

// SalesReport represents the sales summary response
// @Description Sales summary report with revenue, transaction count, and best seller
type SalesReport struct {
	TotalRevenue       int                 `json:"total_revenue" example:"45000"`
	TotalTransactions  int                 `json:"total_transactions" example:"5"`
	BestSellingProduct *BestSellingProduct `json:"best_selling_product"`
}

// BestSellingProduct represents the best selling product in a report
// @Description Best selling product information
type BestSellingProduct struct {
	Name    string `json:"name" example:"Indomie Goreng"`
	QtySold int    `json:"qty_sold" example:"12"`
}
