package models

import "time"

type Transaction struct {
	ID           string     `json:"id" db:"id"`
	ItemID       string     `json:"item_id" db:"item_id"`
	UserID       string     `json:"user_id" db:"user_id"`
	Status       string     `json:"status" db:"status"` // BORROWED, RETURNED, CANCELLED
	BorrowDate   time.Time  `json:"borrow_date" db:"borrow_date"`
	ExpectedDate *time.Time `json:"expected_return_date,omitempty" db:"expected_return_date"`
	ReturnDate   *time.Time `json:"actual_return_date,omitempty" db:"actual_return_date"`
	ReturnCond   *string    `json:"return_condition,omitempty" db:"return_condition"`
	BorrowNotes  *string    `json:"borrow_notes,omitempty" db:"borrow_notes"`
	ReturnNotes  *string    `json:"return_notes,omitempty" db:"return_notes"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`

	// Joined fields
	Item   *Item `json:"item,omitempty"`
	User   *User `json:"user,omitempty"`
}

type BorrowRequest struct {
	ItemCode         string  `json:"item_code" binding:"required"`
	ExpectedReturnDays int   `json:"expected_return_days"` // e.g. 7 days
	Notes            string  `json:"notes"`
}

type ReturnRequest struct {
	ItemCode      string `json:"item_code" binding:"required"`
	Condition     string `json:"condition" binding:"required,oneof=GOOD DAMAGED IN_REPAIR LOST"`
	Notes         string `json:"notes"`
}
