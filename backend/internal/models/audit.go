package models

import "time"

type AuditSession struct {
	ID           string     `json:"id"`
	LocationID   int        `json:"location_id"`
	LocationName string     `json:"location_name,omitempty"`
	UserID       string     `json:"user_id"`
	UserName     string     `json:"user_name,omitempty"`
	Status       string     `json:"status"` // OPEN, CLOSED
	StartedAt    time.Time  `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Derived
	TotalExpected int `json:"total_expected"`
	TotalFound    int `json:"total_found"`
}

type AuditItem struct {
	ID             string    `json:"id"`
	SessionID      string    `json:"session_id"`
	ItemID         string    `json:"item_id"`
	ItemCode       string    `json:"item_code,omitempty"`
	ItemName       string    `json:"item_name,omitempty"`
	FoundCondition string    `json:"found_condition"`
	ScannedAt      time.Time `json:"scanned_at"`
	Notes          *string   `json:"notes,omitempty"`
}

type CreateAuditRequest struct {
	LocationID int    `json:"location_id" binding:"required"`
	Notes      string `json:"notes"`
}

type ScanAuditItemRequest struct {
	ItemCode  string `json:"item_code" binding:"required"`
	Condition string `json:"condition" binding:"required,oneof=GOOD DAMAGED IN_REPAIR LOST"`
	Notes     string `json:"notes"`
}
