package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/vannyezha/sis-inv/internal/models"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type ItemHandler struct {
	db *pgxpool.Pool
}

func NewItemHandler(db *pgxpool.Pool) *ItemHandler {
	return &ItemHandler{db: db}
}

func (h *ItemHandler) List(c *gin.Context) {
	var filter models.ItemFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid filter parameters"))
		return
	}

	// Base query
	query := `
		SELECT i.id, i.code, i.qr_code_data, i.name, i.category_id, c.name as category_name,
		       i.location_id, l.name as location_name, i.location, i.condition, i.status, i.borrower_type, 
		       i.purchase_date::text, i.purchase_price, i.warranty_end_date::text,
		       i.notes, i.photo_url, i.created_at, i.updated_at
		FROM items i
		LEFT JOIN categories c ON i.category_id = c.id
		LEFT JOIN locations l ON i.location_id = l.id
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.code ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" AND i.status = $%d", argIdx)
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Condition != "" {
		query += fmt.Sprintf(" AND i.condition = $%d", argIdx)
		args = append(args, filter.Condition)
		argIdx++
	}
	if filter.Location != "" {
		query += fmt.Sprintf(" AND i.location = $%d", argIdx)
		args = append(args, filter.Location)
		argIdx++
	}
	if filter.CategoryID > 0 {
		query += fmt.Sprintf(" AND i.category_id = $%d", argIdx)
		args = append(args, filter.CategoryID)
		argIdx++
	}
	if filter.BorrowerType != "" {
		query += fmt.Sprintf(" AND i.borrower_type = $%d", argIdx)
		args = append(args, filter.BorrowerType)
		argIdx++
	}

	query += " ORDER BY i.created_at DESC"

	// Pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}
	offset := (filter.Page - 1) * filter.PageSize
	
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, filter.PageSize, offset)

	rows, err := h.db.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch items"))
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var i models.Item
		// Need temporary variables for pointers that might be scanned as null strings
		var pDate, wDate *string
		
		if err := rows.Scan(
			&i.ID, &i.Code, &i.QRCodeData, &i.Name, &i.CategoryID, &i.CategoryName,
			&i.LocationID, &i.LocationName, &i.Location, &i.Condition, &i.Status, &i.BorrowerType,
			&pDate, &i.PurchasePrice, &wDate,
			&i.Notes, &i.PhotoURL, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}
		
		if pDate != nil && *pDate != "" { 
			i.PurchaseDate = pDate 
		}
		if wDate != nil && *wDate != "" { 
			i.WarrantyEndDate = wDate 
		}
		
		items = append(items, i)
	}

	// Also get total count
	countQuery := "SELECT COUNT(*) FROM items i WHERE 1=1" 
	countArgs := []interface{}{}
	countArgIdx := 1
	
	if filter.Search != "" {
		countQuery += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.code ILIKE $%d)", countArgIdx, countArgIdx)
		countArgs = append(countArgs, "%"+filter.Search+"%")
		countArgIdx++
	}
	if filter.Status != "" {
		countQuery += fmt.Sprintf(" AND i.status = $%d", countArgIdx)
		countArgs = append(countArgs, filter.Status)
		countArgIdx++
	}
	if filter.Condition != "" {
		countQuery += fmt.Sprintf(" AND i.condition = $%d", countArgIdx)
		countArgs = append(countArgs, filter.Condition)
		countArgIdx++
	}
	if filter.Location != "" {
		countQuery += fmt.Sprintf(" AND i.location = $%d", countArgIdx)
		countArgs = append(countArgs, filter.Location)
		countArgIdx++
	}
	if filter.CategoryID > 0 {
		countQuery += fmt.Sprintf(" AND i.category_id = $%d", countArgIdx)
		countArgs = append(countArgs, filter.CategoryID)
		countArgIdx++
	}
	if filter.BorrowerType != "" {
		countQuery += fmt.Sprintf(" AND i.borrower_type = $%d", countArgIdx)
		countArgs = append(countArgs, filter.BorrowerType)
		countArgIdx++
	}

	var totalItems int
	err = h.db.QueryRow(context.Background(), countQuery, countArgs...).Scan(&totalItems)
	if err != nil {
		totalItems = 0
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": items,
		"meta": gin.H{
			"page":       filter.Page,
			"page_size":  filter.PageSize,
			"total":      totalItems,
			"total_pages": (totalItems + filter.PageSize - 1) / filter.PageSize,
		},
	}, ""))
}

func (h *ItemHandler) Get(c *gin.Context) {
	id := c.Param("id")

	var i models.Item
	var pDate, wDate *string

	err := h.db.QueryRow(context.Background(),
		`SELECT i.id, i.code, i.qr_code_data, i.name, i.category_id, c.name as category_name,
		        i.location_id, l.name as location_name, i.location, i.condition, i.status, i.borrower_type, 
		        i.purchase_date::text, i.purchase_price, i.warranty_end_date::text,
		        i.notes, i.photo_url, i.created_at, i.updated_at
		 FROM items i
		 LEFT JOIN categories c ON i.category_id = c.id
		 LEFT JOIN locations l ON i.location_id = l.id
		 WHERE i.id = $1`, id,
	).Scan(
		&i.ID, &i.Code, &i.QRCodeData, &i.Name, &i.CategoryID, &i.CategoryName,
		&i.LocationID, &i.LocationName, &i.Location, &i.Condition, &i.Status, &i.BorrowerType,
		&pDate, &i.PurchasePrice, &wDate,
		&i.Notes, &i.PhotoURL, &i.CreatedAt, &i.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Item not found"))
		return
	}
	
	if pDate != nil && *pDate != "" { i.PurchaseDate = pDate }
	if wDate != nil && *wDate != "" { i.WarrantyEndDate = wDate }

	c.JSON(http.StatusOK, utils.SuccessResponse(i, ""))
}

func (h *ItemHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	var i models.Item
	var pDate, wDate *string

	err := h.db.QueryRow(context.Background(),
		`SELECT i.id, i.code, i.qr_code_data, i.name, i.category_id, c.name as category_name,
		        i.location_id, l.name as location_name, i.location, i.condition, i.status, i.borrower_type, 
		        i.purchase_date::text, i.purchase_price, i.warranty_end_date::text,
		        i.notes, i.photo_url, i.created_at, i.updated_at,
		        t.borrow_photo_url, COALESCE(u.full_name, t.student_name) as current_borrower
		 FROM items i
		 LEFT JOIN categories c ON i.category_id = c.id
		 LEFT JOIN locations l ON i.location_id = l.id
		 LEFT JOIN transactions t ON i.id = t.item_id AND t.status = 'ACTIVE'
		 LEFT JOIN users u ON t.user_id = u.id
		 WHERE i.code = $1`, code,
	).Scan(
		&i.ID, &i.Code, &i.QRCodeData, &i.Name, &i.CategoryID, &i.CategoryName,
		&i.LocationID, &i.LocationName, &i.Location, &i.Condition, &i.Status, &i.BorrowerType,
		&pDate, &i.PurchasePrice, &wDate,
		&i.Notes, &i.PhotoURL, &i.CreatedAt, &i.UpdatedAt,
		&i.LastBorrowPhotoURL, &i.CurrentBorrower,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Item not found"))
		return
	}
	
	if pDate != nil && *pDate != "" { i.PurchaseDate = pDate }
	if wDate != nil && *wDate != "" { i.WarrantyEndDate = wDate }

	c.JSON(http.StatusOK, utils.SuccessResponse(i, ""))
}

func (h *ItemHandler) Create(c *gin.Context) {
	var req models.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid data: "+err.Error()))
		return
	}
	
	// Default values if empty
	if req.Condition == "" {
		req.Condition = "GOOD"
	}
	if req.BorrowerType == "" {
		req.BorrowerType = "STAFF_ONLY"
	}

	var itemID string
	err := h.db.QueryRow(context.Background(),
		`INSERT INTO items (code, name, category_id, location_id, location, condition, borrower_type, 
		                   purchase_date, purchase_price, warranty_end_date, notes, photo_url)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
		req.Code, req.Name, req.CategoryID, req.LocationID, req.Location, req.Condition, req.BorrowerType,
		req.PurchaseDate, req.PurchasePrice, req.WarrantyEndDate, req.Notes, req.PhotoURL,
	).Scan(&itemID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create item: " + err.Error()))
		return
	}

	actorId, _ := c.Get("userID")
	utils.LogAudit(h.db, actorId.(string), "CREATE_ITEM", "ITEM", itemID, "Created item: "+req.Code, c.ClientIP())

	c.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"id": itemID}, "Item successfully created"))
}

func (h *ItemHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid data"))
		return
	}

	ctx := context.Background()

	// 1. Fetch current data for audit comparison
	var (
		oldName, oldCode, oldCond, oldStatus, oldBorrowerType, oldLocation string
		oldPrice                                                           float64
	)
	err := h.db.QueryRow(ctx, "SELECT name, code, condition, status, borrower_type, COALESCE(location, ''), COALESCE(purchase_price, 0) FROM items WHERE id = $1", id).
		Scan(&oldName, &oldCode, &oldCond, &oldStatus, &oldBorrowerType, &oldLocation, &oldPrice)

	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Item not found"))
		return
	}

	changes := []string{}
	query := "UPDATE items SET updated_at = $1"
	args := []interface{}{time.Now()}
	argIdx := 2

	if req.Name != nil {
		query += fmt.Sprintf(", name = $%d", argIdx)
		args = append(args, *req.Name)
		argIdx++
		if *req.Name != oldName {
			changes = append(changes, fmt.Sprintf("Name: %s -> %s", oldName, *req.Name))
		}
	}
	if req.CategoryID != nil {
		query += fmt.Sprintf(", category_id = $%d", argIdx)
		args = append(args, *req.CategoryID)
		argIdx++
		changes = append(changes, "Category updated")
	}
	if req.LocationID != nil {
		query += fmt.Sprintf(", location_id = $%d", argIdx)
		args = append(args, *req.LocationID)
		argIdx++
		changes = append(changes, "Location ID updated")
	}
	if req.Location != nil {
		query += fmt.Sprintf(", location = $%d", argIdx)
		args = append(args, *req.Location)
		argIdx++
		if *req.Location != oldLocation {
			changes = append(changes, fmt.Sprintf("Location: %s -> %s", oldLocation, *req.Location))
		}
	}
	if req.Condition != nil {
		query += fmt.Sprintf(", condition = $%d", argIdx)
		args = append(args, *req.Condition)
		argIdx++
		if *req.Condition != oldCond {
			changes = append(changes, fmt.Sprintf("Condition: %s -> %s", oldCond, *req.Condition))
		}
	}
	if req.Status != nil {
		query += fmt.Sprintf(", status = $%d", argIdx)
		args = append(args, *req.Status)
		argIdx++
		if *req.Status != oldStatus {
			changes = append(changes, fmt.Sprintf("Status: %s -> %s", oldStatus, *req.Status))
		}
	}
	if req.BorrowerType != nil {
		query += fmt.Sprintf(", borrower_type = $%d", argIdx)
		args = append(args, *req.BorrowerType)
		argIdx++
		if *req.BorrowerType != oldBorrowerType {
			changes = append(changes, fmt.Sprintf("Borrower: %s -> %s", oldBorrowerType, *req.BorrowerType))
		}
	}
	if req.PurchaseDate != nil {
		query += fmt.Sprintf(", purchase_date = $%d", argIdx)
		args = append(args, *req.PurchaseDate)
		argIdx++
		changes = append(changes, "Purchase date updated")
	}
	if req.PurchasePrice != nil {
		query += fmt.Sprintf(", purchase_price = $%d", argIdx)
		args = append(args, *req.PurchasePrice)
		argIdx++
		if *req.PurchasePrice != oldPrice {
			changes = append(changes, fmt.Sprintf("Price: %.2f -> %.2f", oldPrice, *req.PurchasePrice))
		}
	}
	if req.Notes != nil {
		query += fmt.Sprintf(", notes = $%d", argIdx)
		args = append(args, *req.Notes)
		argIdx++
		changes = append(changes, "Notes updated")
	}
	if req.PhotoURL != nil {
		query += fmt.Sprintf(", photo_url = $%d", argIdx)
		args = append(args, *req.PhotoURL)
		argIdx++
		changes = append(changes, "Photo updated")
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	result, err := h.db.Exec(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update item"))
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Item not found"))
		return
	}

	actorId, _ := c.Get("userID")
	auditDesc := "Updated item."
	if len(changes) > 0 {
		auditDesc = "Updated item. Changes: " + strings.Join(changes, " | ")
	}
	utils.LogAudit(h.db, actorId.(string), "UPDATE_ITEM", "ITEM", id, auditDesc, c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Item successfully updated"))
}

func (h *ItemHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Hard delete for items, but first check if it's currently borrowed
	var status string
	err := h.db.QueryRow(context.Background(), "SELECT status FROM items WHERE id = $1", id).Scan(&status)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Item not found"))
		return
	}
	
	if status == "BORROWED" {
		c.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, "Cannot delete currently borrowed item"))
		return
	}

	result, err := h.db.Exec(context.Background(), "DELETE FROM items WHERE id = $1", id)
	if err != nil {
		// Could be a foreign key constraint violation (has transaction history)
		// Fallback to soft delete: update status to LOST or create a deleted flag if schema supported it
		// For MVP, we'll try to just update status if hard delete fails due to relations
		_, updateErr := h.db.Exec(context.Background(), "UPDATE items SET status = 'LOST', notes = notes || '\n[DELETED FROM SYSTEM]' WHERE id = $1", id)
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to delete item completely, foreign key restriction."))
			return
		}
		c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Item archived (has history)"))
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Item not found"))
		return
	}

	actorId, _ := c.Get("userID")
	utils.LogAudit(h.db, actorId.(string), "DELETE_ITEM", "ITEM", id, "Archived/Deleted item", c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Item successfully deleted"))
}

func (h *ItemHandler) GenerateQRCode(c *gin.Context) {
	id := c.Param("id")

	// Get item code
	var code string
	err := h.db.QueryRow(context.Background(), "SELECT code FROM items WHERE id = $1", id).Scan(&code)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Item not found"))
		return
	}

	// Generate QR Code. We'll store the item code in the QR.
	// You might want to encode a full deep link structure, e.g. "sisinv://item/TEST-123" 
	// or "https://school.edu/items/TEST-123", but for MVP let's just encode the item code.
	qrData := code

	// Use github.com/skip2/go-qrcode in a real app, generate base64
	png, err := qrcode.Encode(qrData, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to generate QR code"))
		return
	}

	base64Image := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)

	// Update the database with the generated data (optional, but good for caching)
	_, _ = h.db.Exec(context.Background(), "UPDATE items SET qr_code_data = $1 WHERE id = $2", qrData, id)

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"qr_code": base64Image,
		"data": qrData,
	}, "QR Code generated successfully"))
}

func (h *ItemHandler) GetHistory(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT t.id, t.borrower_type, 
		       COALESCE(u.full_name, t.student_name) as borrower_name,
		       t.student_class,
		       staff.full_name as staff_name,
		       t.borrowed_at, t.returned_at, t.status, t.return_condition, t.purpose, t.return_notes, t.return_photo_url
		FROM transactions t
		LEFT JOIN users u ON t.user_id = u.id
		LEFT JOIN users staff ON t.borrowed_by = staff.id
		WHERE t.item_id = $1
		ORDER BY t.borrowed_at DESC
		LIMIT 20
	`
	rows, err := h.db.Query(context.Background(), query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch item history"))
		return
	}
	defer rows.Close()

	type TransactionHistoryLite struct {
		ID              string     `json:"id"`
		BorrowerType    string     `json:"borrower_type"`
		BorrowerName    string     `json:"borrower_name"`
		StudentClass    *string    `json:"student_class"`
		StaffName       string     `json:"staff_name"`
		BorrowedAt      time.Time  `json:"borrowed_at"`
		ReturnedAt      *time.Time `json:"returned_at"`
		Status          string     `json:"status"`
		ReturnCondition *string    `json:"return_condition"`
		Purpose         string     `json:"purpose"`
		ReturnNotes     *string    `json:"return_notes"`
		ReturnPhotoURL  *string    `json:"return_photo_url"`
	}

	var history []TransactionHistoryLite
	for rows.Next() {
		var hLite TransactionHistoryLite
		if err := rows.Scan(
			&hLite.ID, &hLite.BorrowerType, &hLite.BorrowerName, &hLite.StudentClass, &hLite.StaffName,
			&hLite.BorrowedAt, &hLite.ReturnedAt, &hLite.Status, &hLite.ReturnCondition, &hLite.Purpose, &hLite.ReturnNotes, &hLite.ReturnPhotoURL,
		); err != nil {
			fmt.Printf("Scan error in history: %v\n", err)
			continue
		}
		history = append(history, hLite)
	}

	if history == nil {
		history = []TransactionHistoryLite{}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(history, ""))
}
func (h *ItemHandler) ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "File tidak ditemukan"))
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal membuka file"))
		return
	}
	defer openedFile.Close()

	f, err := excelize.OpenReader(openedFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Format file Excel tidak valid"))
		return
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil || len(rows) < 2 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Data Excel kosong atau tidak terbaca"))
		return
	}

	ctx := context.Background()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal memulai transaksi"))
		return
	}
	defer tx.Rollback(ctx)

	imported := 0
	for i, row := range rows {
		if i == 0 { continue } // Skip header
		if len(row) < 2 { continue }

		code := strings.TrimSpace(row[0])
		name := strings.TrimSpace(row[1])
		if code == "" || name == "" { continue }

		categoryName := ""
		if len(row) >= 3 { categoryName = strings.TrimSpace(row[2]) }
		
		locationName := ""
		if len(row) >= 4 { locationName = strings.TrimSpace(row[3]) }

		condition := "GOOD"
		if len(row) >= 5 {
			cond := strings.ToUpper(strings.TrimSpace(row[4]))
			if cond == "BAIK" || cond == "GOOD" {
				condition = "GOOD"
			} else if cond == "RUSAK" || cond == "DAMAGED" {
				condition = "DAMAGED"
			} else if cond == "HILANG" || cond == "LOST" {
				condition = "LOST"
			}
		}

		borrowerType := "STAFF_ONLY"
		if len(row) >= 6 {
			bt := strings.ToUpper(strings.TrimSpace(row[5]))
			if strings.Contains(bt, "SISWA") || strings.Contains(bt, "STUDENT") || strings.Contains(bt, "ALLOWED") {
				borrowerType = "STUDENT_ALLOWED"
			}
		}

		notes := ""
		if len(row) >= 9 { notes = strings.TrimSpace(row[8]) }

		// Resolve Category ID
		var categoryId *int
		if categoryName != "" {
			var id int
			err := tx.QueryRow(ctx, "SELECT id FROM categories WHERE name ILIKE $1", categoryName).Scan(&id)
			if err != nil {
				// Create category if not exists
				_ = tx.QueryRow(ctx, "INSERT INTO categories (name) VALUES ($1) RETURNING id", categoryName).Scan(&id)
			}
			categoryId = &id
		}

		// Resolve Location ID
		var locationId *int
		if locationName != "" {
			var id int
			err := tx.QueryRow(ctx, "SELECT id FROM locations WHERE name ILIKE $1", locationName).Scan(&id)
			if err != nil {
				// Create location if not exists
				_ = tx.QueryRow(ctx, "INSERT INTO locations (name) VALUES ($1) RETURNING id", locationName).Scan(&id)
			}
			locationId = &id
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO items (code, name, category_id, location_id, condition, borrower_type, notes, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 'AVAILABLE')
			ON CONFLICT (code) DO UPDATE SET 
				name = EXCLUDED.name,
				category_id = EXCLUDED.category_id,
				location_id = EXCLUDED.location_id,
				condition = EXCLUDED.condition,
				borrower_type = EXCLUDED.borrower_type,
				notes = EXCLUDED.notes,
				updated_at = NOW()
		`, code, name, categoryId, locationId, condition, borrowerType, notes)
		
		if err == nil {
			imported++
		}
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal menyimpan data"))
		return
	}

	actorId, _ := c.Get("userID")
	utils.LogAudit(h.db, actorId.(string), "IMPORT_ITEMS", "ITEM", "00000000-0000-0000-0000-000000000000", fmt.Sprintf("Imported %d items via Excel", imported), c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"total": imported}, "Berhasil mengimpor data barang"))
}
