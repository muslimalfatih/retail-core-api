package main

import (
	"fmt"
	"log"
	"net/http"
	"retail-core-api/config"
	"retail-core-api/database"
	"retail-core-api/docs"
	"retail-core-api/handlers"
	"retail-core-api/helpers"
	"retail-core-api/middleware"
	"retail-core-api/repositories"
	"retail-core-api/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Retail Core API
// @version 1.0
// @description RESTful API for managing categories, products, transactions, and POS operations
// @description
// @description ## Features:
// @description - Category Management (CRUD)
// @description - Product Management (CRUD with category relationship)
// @description - Product Search by Name (case-insensitive partial match)
// @description - Transaction / Checkout (multi-item checkout with stock deduction)
// @description - Sales Reports (daily summary, date range, best selling product)
// @description - Dashboard Statistics

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /
// @schemes http https

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Configure Swagger
	docs.SwaggerInfo.Host = cfg.SwaggerHost()
	docs.SwaggerInfo.Schemes = cfg.SwaggerSchemes()

	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// ============================================
	// DATABASE CONNECTION
	// ============================================
	db, err := database.InitDB(cfg.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Run database migrations
	err = database.RunMigrations(db)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// ============================================
	// DEPENDENCY INJECTION
	// ============================================

	// Repositories
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Services
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo, categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo)

	// Handlers
	categoryHandler := handlers.NewCategoryHandler(categoryService, productService)
	productHandler := handlers.NewProductHandler(productService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// ============================================
	// ROUTER SETUP
	// ============================================
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// ── Health & Info ──────────────────────────
	r.GET("/health", func(c *gin.Context) {
		helpers.OK(c, "Server is running successfully", gin.H{"status": "OK"})
	})

	r.GET("/", func(c *gin.Context) {
		helpers.OK(c, "Retail Core API", gin.H{
			"name":        "Retail Core API",
			"version":     "1.0",
			"status":      "running",
			"description": "RESTful API for managing categories, products, and transactions",
			"endpoints": gin.H{
				"documentation": "/docs/index.html",
				"health":        "/health",
				"categories":    "/categories",
				"products":      "/products",
			},
		})
	})

	// ── Swagger Documentation ─────────────────
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ── Categories (public read) ──────────────
	r.GET("/categories", categoryHandler.List)
	r.GET("/categories/:id", categoryHandler.GetByID)
	r.GET("/categories/:id/products", categoryHandler.GetProducts)
	r.POST("/categories", categoryHandler.Create)
	r.PUT("/categories/:id", categoryHandler.Update)
	r.DELETE("/categories/:id", categoryHandler.Delete)

	// ── Products (public read) ────────────────
	r.GET("/products", productHandler.List)
	r.GET("/products/:id", productHandler.GetByID)
	r.POST("/products", productHandler.Create)
	r.PUT("/products/:id", productHandler.Update)
	r.DELETE("/products/:id", productHandler.Delete)

	// ── API group ─────────────────────────────
	api := r.Group("/api")
	{
		// Transactions / Checkout
		api.POST("/checkout", transactionHandler.Checkout)
		api.GET("/transactions", transactionHandler.ListTransactions)
		api.GET("/transactions/:id", transactionHandler.GetTransactionByID)

		// Dashboard
		api.GET("/dashboard", transactionHandler.Dashboard)

		// Reports
		api.GET("/report/today", transactionHandler.DailyReport)
		api.GET("/report", transactionHandler.ReportByRange)
	}

	// ── Start Server ──────────────────────────
	addr := "0.0.0.0:" + cfg.Port
	fmt.Printf("Server running on %s\n", addr)
	fmt.Printf("API Documentation: http://localhost:%s/docs/index.html\n", cfg.Port)
	log.Println("Available endpoints:")
	log.Println("  GET    /health")
	log.Println("  GET    /categories")
	log.Println("  POST   /categories")
	log.Println("  GET    /categories/:id")
	log.Println("  PUT    /categories/:id")
	log.Println("  DELETE /categories/:id")
	log.Println("  GET    /categories/:id/products")
	log.Println("  GET    /products")
	log.Println("  POST   /products")
	log.Println("  GET    /products/:id")
	log.Println("  PUT    /products/:id")
	log.Println("  DELETE /products/:id")
	log.Println("  POST   /api/checkout")
	log.Println("  GET    /api/transactions")
	log.Println("  GET    /api/transactions/:id")
	log.Println("  GET    /api/dashboard")
	log.Println("  GET    /api/report/today")
	log.Println("  GET    /api/report?start_date=&end_date=")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
