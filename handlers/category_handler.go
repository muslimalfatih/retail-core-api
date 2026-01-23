package handlers

import (
	"category-management-api/database"
	"category-management-api/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// GetAllCategories godoc
// @Summary Get all categories
// @Description Retrieve a list of all categories
// @Tags Categories
// @Produce json
// @Success 200 {object} models.Response{data=[]models.Category} "Successfully retrieved all categories"
// @Router /categories [get]
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories := database.GetAllCategories()
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
func CreateCategory(w http.ResponseWriter, r *http.Request) {
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

	// Validation
	if newCategory.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Category name is required",
		})
		return
	}

	// Add category to database
	createdCategory := database.AddCategory(newCategory)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{
		Status:  true,
		Message: "Category created successfully",
		Data:    createdCategory,
	})
}

// CategoriesHandler handles /categories endpoint
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetAllCategories(w, r)
	case http.MethodPost:
		CreateCategory(w, r)
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
func GetCategoryByID(w http.ResponseWriter, r *http.Request, id int) {
	category, found := database.GetCategoryByID(id)
	if !found {
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
func UpdateCategory(w http.ResponseWriter, r *http.Request, id int) {
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

	// Validation
	if updatedCategory.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Category name is required",
		})
		return
	}

	// Update category
	category, found := database.UpdateCategory(id, updatedCategory)
	if !found {
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
func DeleteCategory(w http.ResponseWriter, r *http.Request, id int) {
	deleted := database.DeleteCategory(id)
	if !deleted {
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
		Message: "Category deleted successfully",
	})
}

// CategoryHandler handles /categories/{id} endpoint
func CategoryHandler(w http.ResponseWriter, r *http.Request) {
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
		GetCategoryByID(w, r, id)
	case http.MethodPut:
		UpdateCategory(w, r, id)
	case http.MethodDelete:
		DeleteCategory(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Response{
			Status:  false,
			Message: "Method not allowed",
		})
	}
}
