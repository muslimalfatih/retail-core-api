package models

import "time"

// Product represents a product entity
// @Description Product information with ID, name, price, stock, and category relationship
type Product struct {
	ID           int       `json:"id" example:"1"`
	Name         string    `json:"name" example:"iPhone 15 Pro" binding:"required"`
	Price        int       `json:"price" example:"15000000" binding:"required"`
	Stock        int       `json:"stock" example:"50" binding:"required"`
	CategoryID   *int      `json:"category_id" example:"1"`
	CategoryName string    `json:"category_name,omitempty" example:"Electronics"`
	CreatedAt    time.Time `json:"created_at" example:"2024-01-30T12:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2024-01-30T12:00:00Z"`
}

// ProductInput represents the input for creating/updating a product
// @Description Input model for creating or updating a product (ID is auto-generated)
type ProductInput struct {
	Name       string `json:"name" example:"iPhone 15 Pro" binding:"required"`
	Price      int    `json:"price" example:"15000000" binding:"required"`
	Stock      int    `json:"stock" example:"50" binding:"required"`
	CategoryID *int   `json:"category_id" example:"1"`
}
