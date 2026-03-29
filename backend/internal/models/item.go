package models

import (
	"time"
)

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	ColorCode   *string   `json:"color_code,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Item struct {
	ID              string     `json:"id"`
	Code            string     `json:"code"`
	QRCodeData      *string    `json:"qr_code_data,omitempty"`
	Name            string     `json:"name"`
	CategoryID      *int       `json:"category_id,omitempty"`
	CategoryName    *string    `json:"category_name,omitempty"`
	Location        *string    `json:"location,omitempty"`
	Condition       string     `json:"condition"`
	Status          string     `json:"status"`
	BorrowerType    string     `json:"borrower_type"`
	PurchaseDate    *string    `json:"purchase_date,omitempty"`
	PurchasePrice   *float64   `json:"purchase_price,omitempty"`
	WarrantyEndDate *string    `json:"warranty_end_date,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	PhotoURL        *string    `json:"photo_url,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type CreateItemRequest struct {
	Code            string   `json:"code" binding:"required"`
	Name            string   `json:"name" binding:"required"`
	CategoryID      *int     `json:"category_id"`
	Location        *string  `json:"location"`
	Condition       string   `json:"condition" binding:"omitempty,oneof=GOOD DAMAGED LOST IN_REPAIR"`
	BorrowerType    string   `json:"borrower_type" binding:"omitempty,oneof=STAFF_ONLY STUDENT_ALLOWED"`
	PurchaseDate    *string  `json:"purchase_date"`
	PurchasePrice   *float64 `json:"purchase_price"`
	WarrantyEndDate *string  `json:"warranty_end_date"`
	Notes           *string  `json:"notes"`
	PhotoURL        *string  `json:"photo_url"`
}

type UpdateItemRequest struct {
	Name            *string  `json:"name"`
	CategoryID      *int     `json:"category_id"`
	Location        *string  `json:"location"`
	Condition       *string  `json:"condition" binding:"omitempty,oneof=GOOD DAMAGED LOST IN_REPAIR"`
	Status          *string  `json:"status" binding:"omitempty,oneof=AVAILABLE BORROWED MAINTENANCE LOST"`
	BorrowerType    *string  `json:"borrower_type" binding:"omitempty,oneof=STAFF_ONLY STUDENT_ALLOWED"`
	PurchaseDate    *string  `json:"purchase_date"`
	PurchasePrice   *float64 `json:"purchase_price"`
	WarrantyEndDate *string  `json:"warranty_end_date"`
	Notes           *string  `json:"notes"`
	PhotoURL        *string  `json:"photo_url"`
}

type ItemFilter struct {
	Search       string `form:"search"`
	Status       string `form:"status"`
	Condition    string `form:"condition"`
	CategoryID   int    `form:"category_id"`
	Location     string `form:"location"`
	BorrowerType string `form:"borrower_type"`
	Page         int    `form:"page,default=1"`
	PageSize     int    `form:"page_size,default=20"`
}
