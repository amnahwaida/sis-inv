package models

import "time"

type Transaction struct {
	ID              string     `json:"id" db:"id"`
	ItemID          string     `json:"item_id" db:"item_id"`
	BorrowerType    string     `json:"borrower_type" db:"borrower_type"` // STAFF or STUDENT
	UserID          *string    `json:"user_id,omitempty" db:"user_id"`
	StudentNIS      *string    `json:"student_nis,omitempty" db:"student_nis"`
	StudentName     *string    `json:"student_name,omitempty" db:"student_name"`
	StudentClass    *string    `json:"student_class,omitempty" db:"student_class"`
	BorrowedBy      string     `json:"borrowed_by" db:"borrowed_by"` // The staff/teacher who processed it
	BorrowedAt      time.Time  `json:"borrowed_at" db:"borrowed_at"`
	DueDate         time.Time  `json:"due_date" db:"due_date"`
	ReturnedAt      *time.Time `json:"returned_at,omitempty" db:"returned_at"`
	ReturnedTo      *string    `json:"returned_to,omitempty" db:"returned_to"`
	Status          string     `json:"status" db:"status"` // ACTIVE, RETURNED, OVERDUE
	Purpose         string     `json:"purpose" db:"purpose"`
	ReturnCondition *string    `json:"return_condition,omitempty" db:"return_condition"`
	ReturnNotes     *string    `json:"return_notes,omitempty" db:"return_notes"`
	ReturnPhotoURL  *string    `json:"return_photo_url,omitempty" db:"return_photo_url"`
	BorrowPhotoURL  *string    `json:"borrow_photo_url,omitempty" db:"borrow_photo_url"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`

	// Joined fields
	Item       *Item  `json:"item,omitempty"`
	Borrower   string `json:"borrower_name,omitempty"` // Derived name
}

type BorrowRequest struct {
	ItemCode           string `json:"item_code" binding:"required"`
	BorrowerType       string `json:"borrower_type" binding:"required,oneof=STAFF STUDENT"`
	StudentNIS         string `json:"student_nis"`
	StudentName        string `json:"student_name"`
	StudentClass       string  `json:"student_class"`
	ExpectedReturnDays float64 `json:"expected_return_days"`
	Purpose            string  `json:"purpose"`
	PhotoURL           string  `json:"photo_url"`
}

type ReturnRequest struct {
	ItemCode      string `json:"item_code" binding:"required"`
	Condition     string `json:"condition" binding:"required,oneof=GOOD DAMAGED IN_REPAIR LOST"`
	Notes         string `json:"notes"`
	PhotoURL      string `json:"photo_url"`
}
