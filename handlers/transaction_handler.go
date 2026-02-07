package handlers

import (
	"category-management-api/models"
	"category-management-api/services"
	"encoding/json"
	"net/http"
)

// TransactionHandler handles HTTP requests for transactions and reports
type TransactionHandler struct {
	service services.TransactionService
}

// NewTransactionHandler creates a new transaction handler instance
func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// HandleCheckout routes /api/checkout requests by HTTP method
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Method not allowed",
		})
	}
}

// Checkout godoc
// @Summary Process checkout
// @Description Process a checkout with multiple items. Validates product availability, deducts stock, and creates a transaction.
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body models.CheckoutRequest true "Checkout request with items and quantities"
// @Success 201 {object} models.Response{data=models.Transaction} "Checkout successful"
// @Failure 400 {object} models.Response "Invalid request body or validation error"
// @Failure 500 {object} models.Response "Server error or insufficient stock"
// @Router /api/checkout [post]
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Invalid request body",
		})
		return
	}

	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Response{
		Status:  true,
		Message: "Checkout successful",
		Data:    transaction,
	})
}

// HandleDailyReport godoc
// @Summary Get today's sales report
// @Description Get sales summary for today including total revenue, transaction count, and best selling product
// @Tags Reports
// @Produce json
// @Success 200 {object} models.Response{data=models.SalesReport} "Daily sales report retrieved successfully"
// @Failure 500 {object} models.Response "Failed to get daily report"
// @Router /api/report/today [get]
func (h *TransactionHandler) HandleDailyReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Method not allowed",
		})
		return
	}

	report, err := h.service.GetDailySalesReport()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Failed to get daily report: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(models.Response{
		Status:  true,
		Message: "Daily sales report retrieved successfully",
		Data:    report,
	})
}

// HandleReportByRange godoc
// @Summary Get sales report by date range
// @Description Get sales summary for a specific date range including total revenue, transaction count, and best selling product
// @Tags Reports
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD format)" example("2026-01-01")
// @Param end_date query string true "End date (YYYY-MM-DD format)" example("2026-02-01")
// @Success 200 {object} models.Response{data=models.SalesReport} "Sales report retrieved successfully"
// @Failure 400 {object} models.Response "Missing required query parameters"
// @Failure 500 {object} models.Response "Failed to get report"
// @Router /api/report [get]
func (h *TransactionHandler) HandleReportByRange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Method not allowed",
		})
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "start_date and end_date query parameters are required",
		})
		return
	}

	report, err := h.service.GetSalesReportByDateRange(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Failed to get report: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(models.Response{
		Status:  true,
		Message: "Sales report retrieved successfully",
		Data:    report,
	})
}
