package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/models"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type TransactionHandler struct {
	db *pgxpool.Pool
}

func NewTransactionHandler(db *pgxpool.Pool) *TransactionHandler {
	return &TransactionHandler{db: db}
}

// BorrowStaff allows a logged-in user to borrow an item via its code
func (h *TransactionHandler) BorrowStaff(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	var req models.BorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body"))
		return
	}

	ctx := context.Background()

	// Start a transaction since we need to update item status AND create a borrowing record
	tx, err := h.db.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to start transaction"))
		return
	}
	defer tx.Rollback(ctx)

	// Fetch item by code to ensure it exists and is AVAILABLE
	var itemId string
	var currentStatus string
	var itemName string
	err = tx.QueryRow(ctx, "SELECT id, status, name FROM items WHERE code = $1", req.ItemCode).Scan(&itemId, &currentStatus, &itemName)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Item not found"))
		return
	}

	if currentStatus != "AVAILABLE" {
		c.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, "Item is currently "+currentStatus+" and cannot be borrowed"))
		return
	}

	// Update item status
	_, err = tx.Exec(ctx, "UPDATE items SET status = 'BORROWED' WHERE id = $1", itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update item status"))
		return
	}

	// Determine expected return date
	var expectedDate time.Time
	if req.ExpectedReturnDays > 0 {
		expectedDate = time.Now().AddDate(0, 0, req.ExpectedReturnDays)
	} else {
		expectedDate = time.Now().AddDate(0, 0, 7) // Default 7 days
	}

	// Create transaction record
	var trxId string
	query := `
		INSERT INTO transactions 
		(item_id, borrower_type, user_id, borrowed_by, status, due_date, purpose)
		VALUES ($1, 'STAFF', $2, $2, 'ACTIVE', $3, $4)
		RETURNING id
	`
	err = tx.QueryRow(ctx, query, itemId, userId, expectedDate, req.Notes).Scan(&trxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create transaction record"))
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to commit transaction"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"transaction_id": trxId,
		"item_name":      itemName,
		"item_code":      req.ItemCode,
	}, "Item borrowed successfully"))
}

// Return allows returning an item via its code
func (h *TransactionHandler) Return(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	userRole, _ := c.Get("role")

	var req models.ReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body"))
		return
	}

	ctx := context.Background()

	tx, err := h.db.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to start transaction"))
		return
	}
	defer tx.Rollback(ctx)

	// Get Item ID from Code
	var itemId string
	var currentStatus string
	err = tx.QueryRow(ctx, "SELECT id, status FROM items WHERE code = $1", req.ItemCode).Scan(&itemId, &currentStatus)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Item not found"))
		return
	}

	if currentStatus != "BORROWED" && currentStatus != "MAINTENANCE" {
		c.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, "Item is not currently borrowed (Status: "+currentStatus+")"))
		return
	}

	// Find the active transaction
	// If the user is not ADMIN, they can only return items they borrowed themselves OR the system allows any staff to return?
	// The PRD implies staff borrow for themselves. Let's strictly find the transaction for this item.
	var trxId string
	var borrowerId string
	err = tx.QueryRow(ctx, "SELECT id, user_id FROM transactions WHERE item_id = $1 AND status = 'ACTIVE' ORDER BY created_at DESC LIMIT 1", itemId).Scan(&trxId, &borrowerId)
	if err != nil {
		// Possibly borrowed outside the system or record missing
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Active borrow record not found for this item"))
		return
	}

	// For security, non-admins might only return their own items. Let's assume ADMIN can return anything, others only theirs.
	if userRole != "ADMIN" && borrowerId != userId.(string) {
		c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, "You can only return items you borrowed yourself"))
		return
	}

	// Update Transaction
	now := time.Now()
	_, err = tx.Exec(ctx, `
		UPDATE transactions 
		SET status = 'RETURNED', returned_at = $1, return_condition = $2, return_notes = $3, updated_at = NOW() 
		WHERE id = $4
	`, now, req.Condition, req.Notes, trxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update transaction record"))
		return
	}

	// Determine new Item status based on condition
	newItemStatus := "AVAILABLE"
	if req.Condition == "DAMAGED" || req.Condition == "IN_REPAIR" {
		newItemStatus = "MAINTENANCE"
	} else if req.Condition == "LOST" {
		newItemStatus = "LOST"
	}

	// Update Item status & condition
	_, err = tx.Exec(ctx, "UPDATE items SET status = $1, condition = $2, updated_at = NOW() WHERE id = $3", newItemStatus, req.Condition, itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update item status"))
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to commit transaction"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"new_status":    newItemStatus,
		"new_condition": req.Condition,
	}, "Item returned successfully"))
}

// MyBorrowings lists items currently borrowed by the logged-in user
func (h *TransactionHandler) MyBorrowings(c *gin.Context) {
	userId, _ := c.Get("userID")

	query := `
		SELECT 
			t.id, t.item_id, t.status, t.borrowed_at as borrow_date, t.due_date as expected_return_date,
			i.code, i.name, i.category_id, i.location
		FROM transactions t
		JOIN items i ON t.item_id = i.id
		WHERE t.user_id = $1 AND t.status = 'ACTIVE'
		ORDER BY t.borrowed_at DESC
	`

	rows, err := h.db.Query(context.Background(), query, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch borrowings"))
		return
	}
	defer rows.Close()

	var borrowings []map[string]interface{}
	for rows.Next() {
		var id, itemId, status, code, name string
		var borrowDate time.Time
		var expectedDate *time.Time
		var categoryId *string
		var location *string

		err := rows.Scan(&id, &itemId, &status, &borrowDate, &expectedDate, &code, &name, &categoryId, &location)
		if err == nil {
			borrowings = append(borrowings, map[string]interface{}{
				"id":                   id,
				"item_id":              itemId,
				"item_code":            code,
				"item_name":            name,
				"status":               status,
				"borrow_date":          borrowDate,
				"expected_return_date": expectedDate,
			})
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(borrowings, "Fetched active borrowings successfully"))
}
