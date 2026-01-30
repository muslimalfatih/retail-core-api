package handlers

import (
	"category-management-api/models"
	"category-management-api/services"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	service services.CategoryService
}

// NewCategoryHandler creates a new category handler instance
func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Retrieve a list of all categories
// @Tags Categories
// @Produce json
// @Success 200 {object} models.Response{data=[]models.Category} "Successfully retrieved all categories"
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Failed to retrieve categories: " + err.Error(),
		})
		return
	}

	response := models.Response{
		Status:  true,
		Message: "Successfully retrieved all categories",
		Data:    categories,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Add a new category to the database
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.CategoryInput true "Category object that needs to be added"
// @Success 201 {object} models.Response{data=models.Category} "Category created successfully"
// @Failure 400 {object} models.Response "Invalid request body or validation error"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory models.Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Invalid request body",
		})
		return
	}

	// Call service to create category (validation is in service layer)
	createdCategory, err := h.service.CreateCategory(newCategory)
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
		Message: "Category created successfully",
		Data:    createdCategory,
	})
}

// HandleCategories handles /categories endpoint
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.GetAllCategories(w, r)
	case http.MethodPost:
		h.CreateCategory(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Method not allowed",
		})
	}
}

// GetCategoryByID godoc
// @Summary Get a category by ID
// @Description Retrieve details of a specific category by its ID
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Response{data=models.Category} "Category retrieved successfully"
// @Failure 400 {object} models.Response "Invalid category ID"
// @Failure 404 {object} models.Response "Category not found"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request, id int) {
	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Failed to retrieve category: " + err.Error(),
		})
		return
	}

	if category == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Category not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{
		Status:  true,
		Message: "Category retrieved successfully",
		Data:    category,
	})
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category by its ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.CategoryInput true "Updated category object"
// @Success 200 {object} models.Response{data=models.Category} "Category updated successfully"
// @Failure 400 {object} models.Response "Invalid request body or validation error"
// @Failure 404 {object} models.Response "Category not found"
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request, id int) {
	var updatedCategory models.Category
	err := json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Invalid request body",
		})
		return
	}

	// Call service to update category (validation is in service layer)
	category, err := h.service.UpdateCategory(id, updatedCategory)
	if err != nil {
		if err.Error() == "category not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{
		Status:  true,
		Message: "Category updated successfully",
		Data:    category,
	})
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by its ID
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Response "Category deleted successfully"
// @Failure 400 {object} models.Response "Invalid category ID"
// @Failure 404 {object} models.Response "Category not found"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request, id int) {
	err := h.service.DeleteCategory(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.Response{
				Status:  false,
				Message: "Category not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Failed to delete category: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{
		Status:  true,
		Message: "Category deleted successfully",
	})
}

// HandleCategoryByID handles /categories/{id} endpoint
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Invalid category ID",
		})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetCategoryByID(w, r, id)
	case http.MethodPut:
		h.UpdateCategory(w, r, id)
	case http.MethodDelete:
		h.DeleteCategory(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Method not allowed",
		})
	}
}
