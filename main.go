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
	"os"

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

// Global CORS middleware that wraps ALL handlers
type corsMiddleware struct {
	handler http.Handler
}

func (c *corsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for all requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
	w.Header().Set("Access-Control-Max-Age", "3600")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Call the next handler
	c.handler.ServeHTTP(w, r)
}

func corsMiddlewareWrapper(handler http.Handler) http.Handler {
	return &corsMiddleware{handler: handler}
}

func main() {
	// Initialize default categories
	initializeDefaultCategories()

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", HealthCheck)

	// Category endpoints
	mux.HandleFunc("/categories", handlers.CategoriesHandler)
	mux.HandleFunc("/categories/", handlers.CategoryHandler)

	// API Documentation endpoint
	mux.HandleFunc("/docs/", httpSwagger.WrapHandler)

	// Get port from environment 
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server with CORS middleware wrapping all routes
	fmt.Printf("Server running on port %s\n", port)
	fmt.Printf("API Documentation: http://localhost:%s/docs/index.html\n", port)
	log.Println("Available endpoints:")
	log.Println("  GET    /health")
	log.Println("  GET    /categories")
	log.Println("  POST   /categories")
	log.Println("  GET    /categories/{id}")
	log.Println("  PUT    /categories/{id}")
	log.Println("  DELETE /categories/{id}")

	// Wrap the entire mux with CORS middleware
	handler := corsMiddlewareWrapper(mux)

	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
