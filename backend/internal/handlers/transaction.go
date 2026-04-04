package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/models"
	"github.com/vannyezha/sis-inv/internal/utils"
	"github.com/xuri/excelize/v2"
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

	// Pre-check: Ensure the actor (staff) exists in the database to avoid FK violations (stale sessions)
	var actorExists bool
	err := h.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userId).Scan(&actorExists)
	if err != nil || !actorExists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Sesi tidak valid atau akun tidak ditemukan. Silakan login kembali."))
		return
	}

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
	if req.ExpectedReturnDays > 0 {
		hours := req.ExpectedReturnDays * 24
		dueDate = time.Now().Add(time.Duration(hours) * time.Hour)
	} else {
		if req.BorrowerType == "STUDENT" {
			dueDate = time.Now().Add(12 * time.Hour) // Default 12 hours fallback
		} else {
			dueDate = time.Now().Add(7 * 24 * time.Hour) // Default 7 days fallback
		}
	}

	// Create transaction record
	var trxId string
	query := `
		INSERT INTO transactions 
		(item_id, borrower_type, user_id, student_nis, student_name, student_class, borrowed_by, status, due_date, purpose, borrow_photo_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	
	// Normalize IDs
	operatorId := userId.(string)
	var borrowerId interface{} = nil
	if req.BorrowerType == "STAFF" {
		borrowerId = operatorId
	}

	err = tx.QueryRow(ctx, query, 
		itemId, 
		req.BorrowerType, 
		borrowerId, 
		utils.NullString(req.StudentNIS), 
		utils.NullString(req.StudentName), 
		utils.NullString(req.StudentClass),
		operatorId, 
		"ACTIVE",
		dueDate, 
		req.Purpose, 
		utils.NullString(req.PhotoURL),
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

	actorId, _ := c.Get("userID")

	// Resolve operator (teacher/admin) name for audit
	var operatorName string
	_ = h.db.QueryRow(context.Background(), "SELECT full_name FROM users WHERE id = $1", actorId).Scan(&operatorName)
	if operatorName == "" {
		operatorName = "Unknown"
	}

	var borrowerInfo string
	if req.BorrowerType == "STUDENT" {
		borrowerInfo = fmt.Sprintf("Siswa: %s (%s, NIS: %s)", req.StudentName, req.StudentClass, req.StudentNIS)
	} else {
		borrowerInfo = "Staff/Ybs"
	}

	purpose := req.Purpose
	if purpose == "" {
		purpose = "Tidak disebutkan"
	}

	auditDesc := fmt.Sprintf("Peminjaman barang [%s] '%s' kepada %s. Durasi: %d hari. Tujuan: %s. Difasilitasi oleh: %s.", 
		req.ItemCode, itemName, borrowerInfo, int(req.ExpectedReturnDays), purpose, operatorName)
	utils.LogAudit(h.db, actorId.(string), "BORROW_ITEM", "ITEM", itemId, auditDesc, c.ClientIP())

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

	// Find the active transaction with borrower details for comprehensive logging
	var trxId string
	var borrowedBy string
	var borrowerType string
	var studentName, studentClass *string
	var facilitatorName string
	err = tx.QueryRow(ctx, `
		SELECT t.id, t.borrowed_by, t.borrower_type, t.student_name, t.student_class, u.full_name
		FROM transactions t
		JOIN users u ON t.borrowed_by = u.id
		WHERE t.item_id = $1 AND t.status = 'ACTIVE'
		ORDER BY t.borrowed_at DESC LIMIT 1
	`, itemId).Scan(&trxId, &borrowedBy, &borrowerType, &studentName, &studentClass, &facilitatorName)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Active borrow record not found"))
		return
	}

	// Any authenticated staff can accept the return
	
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

	// Build comprehensive audit description
	actorId, _ := c.Get("userID")

	// Resolve receiver name
	var receiverName string
	_ = h.db.QueryRow(context.Background(), "SELECT full_name FROM users WHERE id = $1", actorId).Scan(&receiverName)
	if receiverName == "" {
		receiverName = "Unknown"
	}

	// Build borrower info string
	var borrowerInfo string
	if borrowerType == "STUDENT" && studentName != nil {
		borrowerInfo = fmt.Sprintf("Siswa: %s", *studentName)
		if studentClass != nil {
			borrowerInfo += fmt.Sprintf(" (%s)", *studentClass)
		}
		borrowerInfo += fmt.Sprintf(" | Perantara Guru: %s", facilitatorName)
	} else {
		borrowerInfo = fmt.Sprintf("Staff: %s", facilitatorName)
	}

	auditDesc := fmt.Sprintf("Pengembalian barang [%s]. Peminjam: %s. Diterima oleh: %s. Kondisi: %s", 
		req.ItemCode, borrowerInfo, receiverName, req.Condition)
	if req.Notes != "" {
		auditDesc += fmt.Sprintf(". Catatan: %s", req.Notes)
	}
	utils.LogAudit(h.db, actorId.(string), "RETURN_ITEM", "ITEM", itemId, auditDesc, c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"new_status": newItemStatus}, "Item returned successfully"))
}

// MyBorrowings lists items currently borrowed by the logged-in user or by students via this user
func (h *TransactionHandler) MyBorrowings(c *gin.Context) {
	userId, _ := c.Get("userID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	offset := (page - 1) * pageSize
	search := c.Query("search")

	ctx := context.Background()

	baseQuery := " FROM transactions t JOIN items i ON t.item_id = i.id WHERE t.borrowed_by = $1 AND t.status = 'ACTIVE'"
	var args []interface{}
	args = append(args, userId)
	argIdx := 2

	if search != "" {
		searchStr := "%" + search + "%"
		baseQuery += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.code ILIKE $%d)", argIdx, argIdx)
		args = append(args, searchStr)
		argIdx++
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	err := h.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil { total = 0 }

	query := `
		SELECT 
			t.id, t.item_id, t.status, t.borrowed_at, t.due_date,
			t.borrower_type, COALESCE(t.student_name, ''), COALESCE(t.student_class, ''), t.borrow_photo_url,
			i.code, i.name
	` + baseQuery + `
		ORDER BY t.borrowed_at DESC
		LIMIT $` + fmt.Sprint(argIdx) + ` OFFSET $` + fmt.Sprint(argIdx+1)
	
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch borrowings"))
		return
	}
	defer rows.Close()

	borrowings := []map[string]interface{}{}
	for rows.Next() {
		var id, itemId, status, code, itemName, borrowerType string
		var borrowDate, dueDate time.Time
		var studentName, studentClass, borrowPhoto *string

		err := rows.Scan(&id, &itemId, &status, &borrowDate, &dueDate, &borrowerType, &studentName, &studentClass, &borrowPhoto, &code, &itemName)
		if err == nil {
			displayName := "Diri Sendiri"
			if borrowerType == "STUDENT" && studentName != nil {
				displayName = *studentName
				if studentClass != nil {
					displayName += " (" + *studentClass + ")"
				}
			}

			borrowings = append(borrowings, map[string]interface{}{
				"id":               id,
				"item_code":        code,
				"item_name":        itemName,
				"borrower_type":    borrowerType,
				"borrower_display": displayName,
				"student_name":     studentName,
				"student_class":    studentClass,
				"borrow_date":      borrowDate,
				"due_date":         dueDate,
				"borrow_photo_url": borrowPhoto,
				"status":           status,
				"is_overdue":       time.Now().After(dueDate),
			})
		}
	}

	if borrowings == nil { borrowings = []map[string]interface{}{} }

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": borrowings,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}, "Fetched active borrowings successfully"))
}

// MyBorrowingsHistory returns historical (RETURNED) transactions managed by this user
func (h *TransactionHandler) MyBorrowingsHistory(c *gin.Context) {
	userId, _ := c.Get("userID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	offset := (page - 1) * pageSize
	search := c.Query("search")

	ctx := context.Background()

	baseQuery := " FROM transactions t JOIN items i ON t.item_id = i.id WHERE t.borrowed_by = $1 AND t.status = 'RETURNED'"
	var args []interface{}
	args = append(args, userId)
	argIdx := 2

	if search != "" {
		searchStr := "%" + search + "%"
		baseQuery += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.code ILIKE $%d)", argIdx, argIdx)
		args = append(args, searchStr)
		argIdx++
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	err := h.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil { total = 0 }

	query := `
		SELECT 
			t.id, t.status, t.borrowed_at, t.returned_at, t.due_date,
			t.borrower_type, COALESCE(t.student_name, ''), COALESCE(t.student_class, ''),
			t.return_condition, t.return_notes, t.borrow_photo_url, t.return_photo_url,
			i.code, i.name
	` + baseQuery + `
		ORDER BY t.returned_at DESC
		LIMIT $` + fmt.Sprint(argIdx) + ` OFFSET $` + fmt.Sprint(argIdx+1)
	
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch history"))
		return
	}
	defer rows.Close()

	history := []map[string]interface{}{}
	for rows.Next() {
		var id, status, code, itemName, borrowerType string
		var borrowDate, returnedAt, dueDate time.Time
		var studentName, studentClass, returnCondition, returnNotes, borrowPhoto, returnPhoto *string

		err := rows.Scan(&id, &status, &borrowDate, &returnedAt, &dueDate, &borrowerType, &studentName, &studentClass, &returnCondition, &returnNotes, &borrowPhoto, &returnPhoto, &code, &itemName)
		if err == nil {
			displayName := "Diri Sendiri"
			if borrowerType == "STUDENT" && studentName != nil {
				displayName = *studentName
			}

			history = append(history, map[string]interface{}{
				"id":               id,
				"item_code":        code,
				"item_name":        itemName,
				"borrower_display": displayName,
				"student_name":     studentName,
				"student_class":    studentClass,
				"borrowed_at":      borrowDate,
				"returned_at":      returnedAt,
				"due_date":         dueDate,
				"return_condition": returnCondition,
				"return_notes":     returnNotes,
				"borrow_photo_url": borrowPhoto,
				"return_photo_url": returnPhoto,
				"status":           status,
			})
		}
	}

	if history == nil { history = []map[string]interface{}{} }

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": history,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}, "Fetched history successfully"))
}

// GetStudentHistory returns all transactions for a specific student (F04.6)
func (h *TransactionHandler) GetStudentHistory(c *gin.Context) {
	nis := c.Param("nis")
	if nis == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Student NIS is required"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	offset := (page - 1) * pageSize

	search := c.Query("search")

	ctx := context.Background()

	baseQuery := " FROM transactions t JOIN items i ON t.item_id = i.id LEFT JOIN users u ON t.borrowed_by = u.id WHERE t.student_nis = $1 OR t.student_name ILIKE $1"
	var args []interface{}
	args = append(args, nis)
	argIdx := 2

	if search != "" {
		searchStr := "%" + search + "%"
		baseQuery += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.code ILIKE $%d)", argIdx, argIdx)
		args = append(args, searchStr)
		argIdx++
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	err := h.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil { total = 0 }

	query := `
		SELECT 
			t.id, t.status, t.borrowed_at, t.returned_at, t.due_date,
			COALESCE(t.student_name, ''), COALESCE(t.student_class, ''), COALESCE(t.purpose, ''),
			t.borrow_photo_url, t.return_photo_url, t.return_condition,
			i.code as item_code, i.name as item_name,
			COALESCE(u.full_name, 'System') as teacher_name
	` + baseQuery + `
		ORDER BY t.borrowed_at DESC
		LIMIT $` + fmt.Sprint(argIdx) + ` OFFSET $` + fmt.Sprint(argIdx+1)

	args = append(args, pageSize, offset)

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch student history"))
		return
	}
	defer rows.Close()

	history := []map[string]interface{}{}
	for rows.Next() {
		var id, status, itemCode, itemName, teacherName string
		var studentName, studentClass, purpose, borrowPhoto, returnPhoto, returnCondition *string
		var borrowedAt time.Time
		var returnedAt, dueDate *time.Time

		err := rows.Scan(
			&id, &status, &borrowedAt, &returnedAt, &dueDate,
			&studentName, &studentClass, &purpose,
			&borrowPhoto, &returnPhoto, &returnCondition,
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
				"item_code":        itemCode,
				"item_name":        itemName,
				"teacher_name":     teacherName,
				"borrow_photo_url": borrowPhoto,
				"return_photo_url": returnPhoto,
				"return_condition": returnCondition,
			})
		}
	}

	if history == nil { history = []map[string]interface{}{} }

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": history,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}, "Fetched student history successfully"))
}

func (h *TransactionHandler) StudentHistoryExport(c *gin.Context) {
	nis := c.Param("nis")
	if nis == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Student NIS is required"))
		return
	}

	query := `
		SELECT 
			t.status, t.borrowed_at, t.returned_at,
			t.student_name, t.student_class, t.purpose,
			t.return_condition, 
			i.code as item_code, i.name as item_name,
			COALESCE(u.full_name, 'System') as teacher_name
		FROM transactions t
		JOIN items i ON t.item_id = i.id
		LEFT JOIN users u ON t.borrowed_by = u.id
		WHERE t.student_nis = $1 OR t.student_name ILIKE $1
		ORDER BY t.borrowed_at DESC
	`

	rows, err := h.db.Query(context.Background(), query, nis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data untuk ekspor"))
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing excel file: %v", err)
		}
	}()

	sheet := "Riwayat Peminjaman"
	index, _ := f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	// Headers
	hdrs := []string{"Tanggal Pinjam", "Tanggal Kembali", "Aset", "Kode Aset", "Guru Piket", "Keperluan", "Kondisi", "Status"}
	for i, hd := range hdrs {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, hd)
	}

	// Stylings
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4F46E5"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetRowStyle(sheet, 1, 1, style)

	row := 2
	var studentNameForTitle string = nis
	rowCount := 0
	for rows.Next() {
		var status, itemCode, itemName, teacherName string
		var stdName, stdClass, purpose, resCondition *string
		var borrowedAt time.Time
		var returnedAt *time.Time

		err := rows.Scan(
			&status, &borrowedAt, &returnedAt,
			&stdName, &stdClass, &purpose, &resCondition,
			&itemCode, &itemName, &teacherName,
		)
		if err == nil {
			if rowCount == 0 && stdName != nil && *stdName != "" { 
				studentNameForTitle = *stdName 
			}
			
			retDate := "-"
			if returnedAt != nil { retDate = returnedAt.Format("02 Jan 2006 15:04") }

			cond := "-"
			if resCondition != nil { cond = *resCondition }

			purp := "-"
			if purpose != nil { purp = *purpose }

			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), borrowedAt.Format("02 Jan 2006 15:04"))
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), retDate)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), itemName)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), itemCode)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), teacherName)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), purp)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), cond)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), status)
			row++
			rowCount++
		}
	}

	f.SetActiveSheet(index)
	
	// Sanitize filename to prevent header injection or filesystem issues
	sanitizedName := strings.Map(func(r rune) rune {
		if strings.ContainsRune("/\\?%*:|\"<>", r) {
			return '-'
		}
		if r < 32 || r > 126 { // Basic ASCII only for safe filename in headers
			return '_'
		}
		return r
	}, studentNameForTitle)

	if sanitizedName == "" {
		sanitizedName = "Student"
	}

	filename := fmt.Sprintf("Riwayat_Siswa_%s_%s.xlsx", sanitizedName, time.Now().Format("2006-01-02"))
	
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Transfer-Encoding", "binary")
	
	if err := f.Write(c.Writer); err != nil {
		log.Printf("Error writing excel to response: %v", err)
	}
}
