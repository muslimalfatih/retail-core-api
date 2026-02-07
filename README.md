# Category & Product Management API

RESTful API for managing categories, products, transactions, and sales reports with layered architecture pattern, built with Go and PostgreSQL.

## Architecture

This project implements **Layered Architecture** (also known as N-Tier Architecture) with clear separation of concerns:

```
┌─────────────────────────────────────┐
│         Handler Layer               │  ← HTTP Request/Response handling
│  (Presentation/API Layer)           │
├─────────────────────────────────────┤
│         Service Layer               │  ← Business Logic & Validation
│  (Business Logic Layer)             │
├─────────────────────────────────────┤
│       Repository Layer              │  ← Data Access & Persistence
│  (Data Access Layer)                │
├─────────────────────────────────────┤
│           Database                  │  ← PostgreSQL (Supabase)
│  (Data Storage)                     │
└─────────────────────────────────────┘
```

### Layer Responsibilities

1. **Handler Layer** (`handlers/`)
   - Receives HTTP requests
   - Validates request format
   - Returns HTTP responses
   - Error: Request/Response issues

2. **Service Layer** (`services/`)
   - Business logic validation
   - Data transformation
   - Orchestrates repository calls
   - Error: Business logic issues

3. **Repository Layer** (`repositories/`)
   - Database queries
   - Data persistence
   - SQL operations
   - Error: Database issues

4. **Model Layer** (`models/`)
   - Data structures
   - Request/Response schemas

## Features

### Categories Management
- Get all categories
- Get category by ID
- Create new category
- Update existing category
- Delete category

### Products Management
- Get all products 
- Get product by ID
- Create new product
- Update existing product
- Delete product
- Optional category relationship (Foreign Key)
- Category validation on create/update

### Transactions (Checkout)
- Process multi-item checkout
- Automatic stock deduction
- Transaction with detail items
- Product availability validation

### Sales Reports
- Daily sales report (today)
- Sales report by date range
- Total revenue & transaction count
- Best selling product tracking

### Technical Features
- Layered Architecture with Dependency Injection
- PostgreSQL database with `pgx/v5` driver (optimized for Supabase)
- Configuration management with `spf13/viper`
- Connection pooling with lifecycle management
- Environment-based configuration (`APP_ENV` for production/local)
- Automatic database migrations
- SQL JOIN for product-category relationships
- Foreign Key constraints with ON DELETE SET NULL / ON DELETE CASCADE
- Database indexes for performance
- CORS enabled for all endpoints
- Swagger/OpenAPI documentation
- Standard JSON response format
- Production deployment support (Zeabur)

## Getting Started

### Prerequisites

- Go 1.24 or higher
- PostgreSQL database (or Supabase account)
- Supabase connection with SSL enabled

### Installation

1. Clone the repository
```bash
git clone <your-repo-url>
cd category-management-api
```

2. Install dependencies
```bash
go mod download
```

3. Configure environment variables
```bash
cp .env.example .env
# Edit .env with your Supabase credentials
# Important: Add ?sslmode=require to your DB_CONN
```

**Example `.env`:**
```env
DB_CONN=postgresql://postgres.[PROJECT_ID]:[PASSWORD]@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require
PORT=8080
APP_ENV=local
```

4. Run the application
```bash
go run main.go
```

The server will start on `http://localhost:8080` and automatically:
- Connect to PostgreSQL
- Run database migrations (create tables if needed)
- Set up all API routes

## API Documentation

### Swagger UI
Access interactive API documentation at: `http://localhost:8080/docs/index.html`

### Available Endpoints

#### Root & Health
```
GET /       - API information and available endpoints
GET /health - Check API status
```

#### Categories
```
GET    /categories     - Get all categories
POST   /categories     - Create a new category
GET    /categories/:id - Get category by ID
PUT    /categories/:id - Update category
DELETE /categories/:id - Delete category
```

#### Products
```
GET    /products     - Get all products (with category names)
POST   /products     - Create a new product
GET    /products/:id - Get product by ID (with category name)
PUT    /products/:id - Update product
DELETE /products/:id - Delete product
```

#### Transactions
```
POST   /api/checkout       - Process multi-item checkout
```

#### Reports
```
GET    /api/report/today   - Get today's sales report
GET    /api/report          - Get sales report by date range (?start_date=&end_date=)
```

### Request/Response Examples

#### Create Category
```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electronics",
    "description": "Electronic devices and gadgets"
  }'
```

Response:
```json
{
  "status": true,
  "message": "Category created successfully",
  "data": {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and gadgets",
    "created_at": "2024-01-30T12:00:00Z",
    "updated_at": "2024-01-30T12:00:00Z"
  }
}
```

#### Create Product (with category)
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro",
    "price": 15000000,
    "stock": 50,
    "category_id": 1
  }'
```

Response:
```json
{
  "status": true,
  "message": "Product created successfully",
  "data": {
    "id": 1,
    "name": "iPhone 15 Pro",
    "price": 15000000,
    "stock": 50,
    "category_id": 1,
    "category_name": "Electronics",
    "created_at": "2024-01-30T12:00:00Z",
    "updated_at": "2024-01-30T12:00:00Z"
  }
}
```

#### Get Product by ID (with JOIN)
```bash
curl http://localhost:8080/products/1
```

Response shows product with category name fetched via SQL JOIN:
```json
{
  "status": true,
  "message": "Product retrieved successfully",
  "data": {
    "id": 1,
    "name": "iPhone 15 Pro",
    "price": 15000000,
    "stock": 50,
    "category_id": 1,
    "category_name": "Electronics",
    "created_at": "2024-01-30T12:00:00Z",
    "updated_at": "2024-01-30T12:00:00Z"
  }
}
```

#### Checkout (Create Transaction)
```bash
curl -X POST http://localhost:8080/api/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      { "product_id": 1, "quantity": 2 },
      { "product_id": 3, "quantity": 5 }
    ]
  }'
```

Response:
```json
{
  "status": true,
  "message": "Checkout successful",
  "data": {
    "id": 1,
    "total_amount": 45000000,
    "created_at": "2026-02-08T12:00:00Z",
    "details": [
      {
        "id": 1,
        "transaction_id": 1,
        "product_id": 1,
        "product_name": "iPhone 15 Pro",
        "quantity": 2,
        "subtotal": 30000000
      },
      {
        "id": 2,
        "transaction_id": 1,
        "product_id": 3,
        "product_name": "Indomie Goreng",
        "quantity": 5,
        "subtotal": 15000
      }
    ]
  }
}
```

#### Get Today's Sales Report
```bash
curl http://localhost:8080/api/report/today
```

Response:
```json
{
  "status": true,
  "message": "Daily sales report retrieved successfully",
  "data": {
    "total_revenue": 45000000,
    "total_transactions": 5,
    "best_selling_product": {
      "name": "Indomie Goreng",
      "qty_sold": 12
    }
  }
}
```

#### Get Sales Report by Date Range
```bash
curl "http://localhost:8080/api/report?start_date=2026-01-01&end_date=2026-02-08"
```

Response:
```json
{
  "status": true,
  "message": "Sales report retrieved successfully",
  "data": {
    "total_revenue": 120000000,
    "total_transactions": 15,
    "best_selling_product": {
      "name": "Indomie Goreng",
      "qty_sold": 30
    }
  }
}
```

## Database Schema

### Categories Table
```sql
CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Products Table
```sql
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  price INTEGER NOT NULL DEFAULT 0,
  stock INTEGER NOT NULL DEFAULT 0,
  category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_category_id ON products(category_id);
```

**Foreign Key Behavior:**
- `category_id` references `categories(id)`
- `ON DELETE SET NULL`: If a category is deleted, products in that category will have `category_id` set to NULL

### Transactions Table
```sql
CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  total_amount INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Transaction Details Table
```sql
CREATE TABLE transaction_details (
  id SERIAL PRIMARY KEY,
  transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
  product_id INT REFERENCES products(id),
  quantity INT NOT NULL,
  subtotal INT NOT NULL
);
```

**Foreign Key Behavior:**
- `transaction_id` references `transactions(id)` with `ON DELETE CASCADE`: If a transaction is deleted, all its details are also deleted
- `product_id` references `products(id)`

## Development

### Project Structure
```
category-management-api/
├── main.go                 # Application entry point & dependency injection
├── .env                    # Environment configuration
├── .env.example           # Example environment file
├── go.mod                 # Go module dependencies
├── database/
│   ├── postgres.go        # PostgreSQL connection
│   └── migration.go       # Database migrations
├── models/
│   ├── category.go        # Category data structures
│   ├── product.go         # Product data structures
│   └── transaction.go     # Transaction & report data structures
├── repositories/
│   ├── category_repository.go     # Category data access
│   ├── product_repository.go      # Product data access (with JOINs)
│   └── transaction_repository.go  # Transaction & report data access
├── services/
│   ├── category_service.go        # Category business logic
│   ├── product_service.go         # Product business logic
│   └── transaction_service.go     # Transaction & report business logic
├── handlers/
│   ├── category_handler.go        # Category HTTP handlers
│   ├── product_handler.go         # Product HTTP handlers
│   └── transaction_handler.go     # Transaction & report HTTP handlers
└── docs/                  # Swagger documentation (auto-generated)
```

### Regenerate Swagger Docs
If you modify API annotations in code:
```bash
swag init
# or
~/go/bin/swag init
```

**Note:** Only `docs/docs.go` is required. The `swagger.json` and `swagger.yaml` files are optional exports.

### Build
```bash
go build
```

- [Layered Architecture - Martin Fowler](https://martinfowler.com/bliki/PresentationDomainDataLayering.html)
- [Clean Architecture - Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
