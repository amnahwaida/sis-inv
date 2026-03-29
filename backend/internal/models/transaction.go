package models

import (
	"time"
)

type Transaction struct {
	ID              string     `json:"id"`
	ItemID          string     `json:"item_id"`
	ItemName        *string    `json:"item_name,omitempty"`
	ItemCode        *string    `json:"item_code,omitempty"`
	BorrowerType    string     `json:"borrower_type"`
	UserID          *string    `json:"user_id,omitempty"`
	UserName        *string    `json:"user_name,omitempty"`
	StudentNIS      *string    `json:"student_nis,omitempty"`
	StudentName     *string    `json:"student_name,omitempty"`
	StudentClass    *string    `json:"student_class,omitempty"`
	BorrowedBy      string     `json:"borrowed_by"`
	BorrowedByName  *string    `json:"borrowed_by_name,omitempty"`
	BorrowedAt      time.Time  `json:"borrowed_at"`
	DueDate         time.Time  `json:"due_date"`
	ReturnedAt      *time.Time `json:"returned_at,omitempty"`
	Status          string     `json:"status"`
	Purpose         *string    `json:"purpose,omitempty"`
	ReturnCondition *string    `json:"return_condition,omitempty"`
	ReturnNotes     *string    `json:"return_notes,omitempty"`
	ReturnPhotoURL  *string    `json:"return_photo_url,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type BorrowStaffRequest struct {
	ItemCode string  `json:"item_code" binding:"required"`
	Purpose  *string `json:"purpose"`
	DueDate  string  `json:"due_date" binding:"required"`
}

type BorrowStudentRequest struct {
	ItemCode     string  `json:"item_code" binding:"required"`
	StudentNIS   string  `json:"student_nis" binding:"required"`
	StudentName  string  `json:"student_name" binding:"required"`
	StudentClass string  `json:"student_class" binding:"required"`
	Purpose      *string `json:"purpose"`
}

type ReturnRequest struct {
	ItemCode        string  `json:"item_code" binding:"required"`
	ReturnCondition string  `json:"return_condition" binding:"required,oneof=GOOD DAMAGED"`
	ReturnNotes     *string `json:"return_notes"`
}

type Student struct {
	ID        int       `json:"id"`
	NIS       string    `json:"nis"`
	FullName  string    `json:"full_name"`
	Class     *string   `json:"class,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
