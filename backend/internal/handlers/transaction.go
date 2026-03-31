package handlers

import (
	"context"
	"log"
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

// Borrow handles both STAFF and STUDENT borrowing
func (h *TransactionHandler) Borrow(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	var req models.BorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body: "+err.Error()))
		return
	}

	ctx := context.Background()

	// Start a transaction
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
	var borrowerTypeRule string
	err = tx.QueryRow(ctx, "SELECT id, status, name, borrower_type FROM items WHERE code = $1", req.ItemCode).Scan(&itemId, &currentStatus, &itemName, &borrowerTypeRule)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Item not found"))
		return
	}

	if currentStatus != "AVAILABLE" {
		c.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, "Item is currently "+currentStatus+" and cannot be borrowed"))
		return
	}

	// Validation: If it's a student borrowing, check if the item is allowed for students
	if req.BorrowerType == "STUDENT" && borrowerTypeRule == "STAFF_ONLY" {
		c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, "This item is restricted to Staff only and cannot be borrowed by students"))
		return
	}

	// Update item status
	_, err = tx.Exec(ctx, "UPDATE items SET status = 'BORROWED' WHERE id = $1", itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update item status"))
		return
	}

	// Determine expected return date (due_date)
	var dueDate time.Time
	if req.BorrowerType == "STUDENT" {
		// Mandatory 6 hours for students as per PRD F04.4
		dueDate = time.Now().Add(6 * time.Hour)
	} else {
		if req.ExpectedReturnDays > 0 {
			dueDate = time.Now().AddDate(0, 0, req.ExpectedReturnDays)
		} else {
			dueDate = time.Now().AddDate(0, 0, 7) // Default 7 days
		}
	}

	// Create transaction record
	var trxId string
	query := `
		INSERT INTO transactions 
		(item_id, borrower_type, user_id, student_nis, student_name, student_class, borrowed_by, status, due_date, purpose)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'ACTIVE', $8, $9)
		RETURNING id
	`
	
	var borrowerUserId *string
	if req.BorrowerType == "STAFF" {
		uid := userId.(string)
		borrowerUserId = &uid
	}

	err = tx.QueryRow(ctx, query, 
		itemId, req.BorrowerType, borrowerUserId, 
		utils.NullString(req.StudentNIS), utils.NullString(req.StudentName), utils.NullString(req.StudentClass),
		userId, dueDate, req.Purpose,
	).Scan(&trxId)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create transaction record: "+err.Error()))
		return
	}

	// NEW: If student borrowing, auto-register/update student in students table
	if req.BorrowerType == "STUDENT" && req.StudentNIS != "" {
		_, err = tx.Exec(ctx, `
			INSERT INTO students (nis, full_name, class, updated_at)
			VALUES ($1, $2, $3, NOW())
			ON CONFLICT (nis) DO UPDATE SET 
				full_name = EXCLUDED.full_name,
				class = EXCLUDED.class,
				updated_at = NOW()
		`, req.StudentNIS, req.StudentName, req.StudentClass)
		if err != nil {
			log.Printf("⚠️  Failed to auto-register student: %v", err)
			// We don't fail the whole transaction for this, as the borrow itself succeeded
		}
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
		"borrower_type":  req.BorrowerType,
	}, "Item borrowed successfully"))
}

// Return handles returning an item
func (h *TransactionHandler) Return(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized"))
		return
	}

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
	var trxId string
	var borrowedBy string
	var borrowerType string
	err = tx.QueryRow(ctx, "SELECT id, borrowed_by, borrower_type FROM transactions WHERE item_id = $1 AND status = 'ACTIVE' ORDER BY borrowed_at DESC LIMIT 1", itemId).Scan(&trxId, &borrowedBy, &borrowerType)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Active borrow record not found"))
		return
	}

	// RELAXED PERMISSION:
	// - Admins can return anything.
	// - ANY staff/teacher can return an item borrowed for a STUDENT.
	// - Staff return for THEMSELVES: Should strictly be the same person OR any other teacher can receive it too (as per user request).
	// We'll allow any teacher to accept the return of ANY active transaction for practicality.
	
	// Update Transaction
	_, err = tx.Exec(ctx, `
		UPDATE transactions 
		SET status = 'RETURNED', 
		    returned_at = NOW(), 
		    returned_to = $1,
		    return_condition = $2, 
		    return_notes = $3, 
		    return_photo_url = $4,
		    updated_at = NOW() 
		WHERE id = $5
	`, userId, req.Condition, req.Notes, utils.NullString(req.PhotoURL), trxId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update transaction"))
		return
	}

	newItemStatus := "AVAILABLE"
	if req.Condition == "DAMAGED" || req.Condition == "IN_REPAIR" {
		newItemStatus = "MAINTENANCE"
	} else if req.Condition == "LOST" {
		newItemStatus = "LOST"
	}

	_, err = tx.Exec(ctx, "UPDATE items SET status = $1, condition = $2, updated_at = NOW() WHERE id = $3", newItemStatus, req.Condition, itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update item"))
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to commit transaction"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"new_status": newItemStatus}, "Item returned successfully"))
}

// MyBorrowings lists items currently borrowed by the logged-in user or by students via this user
func (h *TransactionHandler) MyBorrowings(c *gin.Context) {
	userId, _ := c.Get("userID")

	query := `
		SELECT 
			t.id, t.item_id, t.status, t.borrowed_at, t.due_date,
			t.borrower_type, t.student_name, t.student_class,
			i.code, i.name
		FROM transactions t
		JOIN items i ON t.item_id = i.id
		WHERE t.borrowed_by = $1 AND t.status = 'ACTIVE'
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
		var id, itemId, status, code, itemName, borrowerType string
		var borrowDate, dueDate time.Time
		var studentName, studentClass *string

		err := rows.Scan(&id, &itemId, &status, &borrowDate, &dueDate, &borrowerType, &studentName, &studentClass, &code, &itemName)
		if err == nil {
			displayName := "Diri Sendiri"
			if borrowerType == "STUDENT" && studentName != nil {
				displayName = *studentName
				if studentClass != nil {
					displayName += " (" + *studentClass + ")"
				}
			}

			borrowings = append(borrowings, map[string]interface{}{
				"id":                   id,
				"item_code":            code,
				"item_name":            itemName,
				"borrower_type":        borrowerType,
				"borrower_display":     displayName,
				"borrow_date":          borrowDate,
				"due_date":             dueDate,
				"status":               status,
			})
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(borrowings, "Fetched active borrowings successfully"))
}

// GetStudentHistory returns all transactions for a specific student (F04.6)
func (h *TransactionHandler) GetStudentHistory(c *gin.Context) {
	nis := c.Param("nis")
	if nis == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Student NIS is required"))
		return
	}

	query := `
		SELECT 
			t.id, t.status, t.borrowed_at, t.returned_at, t.due_date,
			t.student_name, t.student_class, t.purpose,
			i.code as item_code, i.name as item_name,
			u.full_name as teacher_name
		FROM transactions t
		JOIN items i ON t.item_id = i.id
		JOIN users u ON t.borrowed_by = u.id
		WHERE t.student_nis = $1
		ORDER BY t.borrowed_at DESC
	`

	rows, err := h.db.Query(context.Background(), query, nis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch student history"))
		return
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var id, status, itemCode, itemName, teacherName string
		var studentName, studentClass, purpose *string
		var borrowedAt time.Time
		var returnedAt, dueDate *time.Time

		err := rows.Scan(
			&id, &status, &borrowedAt, &returnedAt, &dueDate,
			&studentName, &studentClass, &purpose,
			&itemCode, &itemName, &teacherName,
		)
		if err == nil {
			history = append(history, map[string]interface{}{
				"id":            id,
				"status":        status,
				"borrowed_at":   borrowedAt,
				"returned_at":   returnedAt,
				"due_date":      dueDate,
				"student_name":  studentName,
				"student_class": studentClass,
				"purpose":       purpose,
				"item_code":     itemCode,
				"item_name":     itemName,
				"teacher_name":  teacherName,
			})
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(history, "Fetched student history successfully"))
}
