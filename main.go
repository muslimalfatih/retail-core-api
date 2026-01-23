package main

import (
	"category-management-api/database"
	_ "category-management-api/docs"
	"category-management-api/handlers"
	"category-management-api/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Category Management API
// @version 1.0
// @description RESTful API for managing categories with full CRUD operations
// @description
// @description ## Features:
// @description - Get all categories
// @description - Get category by ID
// @description - Create new category
// @description - Update existing category
// @description - Delete category
// @description
// @description ## Response Format:
// @description All endpoints return a standard response with:
// @description - status (bool): Request success status
// @description - message (string): Response message
// @description - data (object): Response data (when applicable)

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /
// @schemes http https

func initializeDefaultCategories() {
	// Check if categories already exist
	if len(database.Categories) == 0 {
		// Add initial categories
		database.Categories = []models.Category{
			{ID: 1, Name: "Electronics", Description: "Electronic devices and gadgets"},
			{ID: 2, Name: "Clothing", Description: "Apparel and fashion items"},
			{ID: 3, Name: "Books", Description: "Books, magazines, and publications"},
			{ID: 4, Name: "Home & Garden", Description: "Home improvement and gardening supplies"},
			{ID: 5, Name: "Sports", Description: "Sports equipment and accessories"},
		}
		database.NextID = 6
		fmt.Println("✅ Initial categories created successfully")
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Check if the API server is running
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "Server is running"
// @Router /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "Server is running successfully",
	})
}

func main() {
	// Initialize default categories
	initializeDefaultCategories()
	
	// Health check endpoint
	http.HandleFunc("/health", HealthCheck)

	// Category endpoints
	http.HandleFunc("/categories", handlers.CategoriesHandler)
	http.HandleFunc("/categories/", handlers.CategoryHandler)

	// API Documentation endpoint
	http.HandleFunc("/docs/", httpSwagger.WrapHandler)

	// Start server
	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("API Documentation: http://localhost:8080/docs/index.html")
	log.Println("Available endpoints:")
	log.Println("  GET    /health")
	log.Println("  GET    /categories")
	log.Println("  POST   /categories")
	log.Println("  GET    /categories/{id}")
	log.Println("  PUT    /categories/{id}")
	log.Println("  DELETE /categories/{id}")
	
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
