package handlers

import (
	"category-management-api/models"
	"category-management-api/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	service services.ProductService
}

// NewProductHandler creates a new product handler instance
func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Retrieve a list of all products with their category names
// @Tags Products
// @Produce json
// @Success 200 {object} models.Response{data=[]models.Product} "Successfully retrieved all products"
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Failed to retrieve products: " + err.Error(),
		})
		return
	}

	response := models.Response{
		Status:  true,
		Message: "Successfully retrieved all products",
		Data:    products,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Add a new product to the database
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.ProductInput true "Product object that needs to be added"
// @Success 201 {object} models.Response{data=models.Product} "Product created successfully"
// @Failure 400 {object} models.Response "Invalid request body or validation error"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Invalid request body",
		})
		return
	}

	// Call service to create product (validation is in service layer)
	createdProduct, err := h.service.CreateProduct(newProduct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{
		Status:  true,
		Message: "Product created successfully",
		Data:    createdProduct,
	})
}

// HandleProducts handles /products endpoint
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.GetAllProducts(w, r)
	case http.MethodPost:
		h.CreateProduct(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Method not allowed",
		})
	}
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Retrieve details of a specific product by its ID with category name
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Response{data=models.Product} "Product retrieved successfully"
// @Failure 400 {object} models.Response "Invalid product ID"
// @Failure 404 {object} models.Response "Product not found"
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request, id int) {
	product, err := h.service.GetProductByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Failed to retrieve product: " + err.Error(),
		})
		return
	}

	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Product not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{
		Status:  true,
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.ProductInput true "Updated product object"
// @Success 200 {object} models.Response{data=models.Product} "Product updated successfully"
// @Failure 400 {object} models.Response "Invalid request body or validation error"
// @Failure 404 {object} models.Response "Product not found"
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request, id int) {
	var updatedProduct models.Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Invalid request body",
		})
		return
	}

	// Call service to update product
	result, err := h.service.UpdateProduct(id, updatedProduct)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{
		Status:  true,
		Message: "Product updated successfully",
		Data:    result,
	})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Response "Product deleted successfully"
// @Failure 400 {object} models.Response "Invalid product ID"
// @Failure 404 {object} models.Response "Product not found"
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request, id int) {
	err := h.service.DeleteProduct(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to delete product: " + err.Error()

		if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
			message = "Product not found"
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: message,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{
		Status:  true,
		Message: "Product deleted successfully",
	})
}

// HandleProductByID handles /products/{id} endpoint
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Invalid product ID",
		})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetProductByID(w, r, id)
	case http.MethodPut:
		h.UpdateProduct(w, r, id)
	case http.MethodDelete:
		h.DeleteProduct(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Method not allowed",
		})
	}
}
