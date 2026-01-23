# Category Management API

Managing categories operations

## Features

- Get all categories
- Get category by ID
- Create new category
- Update existing category
- Delete category
- Swagger/OpenAPI documentation
- Hot reload for development

## Prerequisites

- Go 1.25+ installed
- Air (for hot reload) - will be installed automatically

## Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```

3. Install Air for hot reload (if not already installed):
```bash
go install github.com/air-verse/air@latest
```

## Running the Application

### Development Mode (with hot reload)

```bash
# Using air
air

# Or using the full path
~/go/bin/air
```

The server will automatically restart when you make changes to any `.go` files.

### Production Mode (without hot reload)

```bash
go run main.go
```

## API Documentation

Once the server is running, access the Swagger UI at:
- **Swagger UI**: http://localhost:8080/docs/index.html
- **Swagger JSON**: http://localhost:8080/docs/doc.json

## Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/categories` | Get all categories |
| POST | `/categories` | Create a new category |
| GET | `/categories/{id}` | Get a category by ID |
| PUT | `/categories/{id}` | Update a category |
| DELETE | `/categories/{id}` | Delete a category |

## Example Requests

### Get All Categories
```bash
curl http://localhost:8080/categories
```

### Create a Category
```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Technology",
    "description": "Technology and gadgets"
  }'
```

### Get Category by ID
```bash
curl http://localhost:8080/categories/1
```

### Update a Category
```bash
curl -X PUT http://localhost:8080/categories/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electronics & Tech",
    "description": "Updated description"
  }'
```

### Delete a Category
```bash
curl -X DELETE http://localhost:8080/categories/1
```

## Response Format

All endpoints return a standard response:
```json
{
  "status": true,
  "message": "Success message",
  "data": { }
}
```

## Initial Data

The application comes with 5 pre-loaded categories:
1. Electronics
2. Clothing
3. Books
4. Home & Garden
5. Sports

## Development

### Regenerate Swagger Documentation

After making changes to API endpoints or models:
```bash
~/go/bin/swag init
```

## Project Structure

```
category-management-api/
├── main.go              # Main application entry point
├── handlers/            # HTTP request handlers
│   └── category_handler.go
├── models/              # Data models
│   └── category.go
├── database/            # In-memory database
│   └── memori.go
├── docs/                # Swagger documentation (auto-generated)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── .air.toml            # Air configuration for hot reload
├── go.mod               # Go module dependencies
└── README.md            # This file
```