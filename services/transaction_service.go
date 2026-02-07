package services

import (
	"category-management-api/models"
	"category-management-api/repositories"
	"errors"
)

// TransactionService defines the interface for transaction business logic
type TransactionService interface {
	Checkout(items []models.CheckoutItem) (*models.Transaction, error)
	GetDailySalesReport() (*models.SalesReport, error)
	GetSalesReportByDateRange(startDate, endDate string) (*models.SalesReport, error)
}

// transactionService implements TransactionService interface
type transactionService struct {
	repo repositories.TransactionRepository
}

// NewTransactionService creates a new transaction service instance
func NewTransactionService(repo repositories.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

// Checkout validates items and delegates to the repository
func (s *transactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	if len(items) == 0 {
		return nil, errors.New("checkout items cannot be empty")
	}

	for _, item := range items {
		if item.ProductID <= 0 {
			return nil, errors.New("invalid product ID")
		}
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0")
		}
	}

	return s.repo.CreateTransaction(items)
}

// GetDailySalesReport returns the sales summary for today
func (s *transactionService) GetDailySalesReport() (*models.SalesReport, error) {
	return s.repo.GetDailySalesReport()
}

// GetSalesReportByDateRange returns the sales summary for a given date range
func (s *transactionService) GetSalesReportByDateRange(startDate, endDate string) (*models.SalesReport, error) {
	if startDate == "" || endDate == "" {
		return nil, errors.New("start_date and end_date are required")
	}
	return s.repo.GetSalesReportByDateRange(startDate, endDate)
}
