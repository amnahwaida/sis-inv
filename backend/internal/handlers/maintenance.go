package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type MaintenanceHandler struct {
	db *pgxpool.Pool
}

func NewMaintenanceHandler(db *pgxpool.Pool) *MaintenanceHandler {
	return &MaintenanceHandler{db: db}
}

// List returns all maintenance logs with item info
func (h *MaintenanceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	statusFilter := c.Query("status")

	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	offset := (page - 1) * pageSize

	baseQuery := `
		FROM maintenance_logs m
		JOIN items i ON m.item_id = i.id
		LEFT JOIN users u ON m.reported_by = u.id
		WHERE 1=1
	`
	var args []interface{}
	argIdx := 1

	if statusFilter != "" {
		baseQuery += fmt.Sprintf(" AND m.status = $%d", argIdx)
		args = append(args, statusFilter)
		argIdx++
	}

	if search != "" {
		searchStr := "%" + search + "%"
		baseQuery += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.code ILIKE $%d OR m.issue_description ILIKE $%d OR m.vendor ILIKE $%d)", argIdx, argIdx, argIdx, argIdx)
		args = append(args, searchStr)
		argIdx++
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	err := h.db.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil { total = 0 }

	// Main query with pagination
	query := `
		SELECT m.id, m.item_id, i.code as item_code, i.name as item_name,
		       COALESCE(u.full_name, 'System') as reported_by_name,
		       m.reported_at, m.issue_description, m.cost, m.vendor,
		       m.status, m.completed_at, m.notes, m.created_at
		` + baseQuery + `
		ORDER BY m.created_at DESC
		LIMIT $` + fmt.Sprint(argIdx) + ` OFFSET $` + fmt.Sprint(argIdx+1)
	
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data maintenance"))
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int
		var itemId, itemCode, itemName, reportedByName, issueDescription, status string
		var reportedAt, createdAt time.Time
		var completedAt *time.Time
		var cost *float64
		var vendor, notes *string

		if err := rows.Scan(&id, &itemId, &itemCode, &itemName,
			&reportedByName, &reportedAt, &issueDescription,
			&cost, &vendor, &status, &completedAt, &notes, &createdAt); err != nil {
			fmt.Printf("Scan error maintenance: %v\n", err)
			continue
		}

		result = append(result, map[string]interface{}{
			"id":                id,
			"item_id":           itemId,
			"item_code":         itemCode,
			"item_name":         itemName,
			"reported_by_name":  reportedByName,
			"reported_at":       reportedAt,
			"issue_description": issueDescription,
			"cost":              cost,
			"vendor":            vendor,
			"status":            status,
			"completed_at":      completedAt,
			"notes":             notes,
			"created_at":        createdAt,
		})
	}

	// Get summary stats
	var pendingCount, inProgressCount, doneCount, cancelledCount int
	statsQuery := `
		SELECT 
			COUNT(*) FILTER (WHERE status = 'PENDING'),
			COUNT(*) FILTER (WHERE status = 'IN_PROGRESS'),
			COUNT(*) FILTER (WHERE status = 'DONE'),
			COUNT(*) FILTER (WHERE status = 'CANCELLED')
		FROM maintenance_logs
	`
	h.db.QueryRow(context.Background(), statsQuery).Scan(&pendingCount, &inProgressCount, &doneCount, &cancelledCount)

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": result,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
		"summary": gin.H{
			"pending":     pendingCount,
			"in_progress": inProgressCount,
			"done":        doneCount,
			"cancelled":   cancelledCount,
		},
	}, ""))
}

// Create adds a new maintenance log entry
func (h *MaintenanceHandler) Create(c *gin.Context) {
	var req struct {
		ItemCode         string  `json:"item_code" binding:"required"`
		IssueDescription string  `json:"issue_description" binding:"required"`
		Cost             float64 `json:"cost"`
		Vendor           string  `json:"vendor"`
		Notes            string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Data tidak valid: "+err.Error()))
		return
	}

	ctx := context.Background()
	userId, _ := c.Get("userID")

	// Find item by code
	var itemId, itemName string
	err := h.db.QueryRow(ctx, "SELECT id, name FROM items WHERE code = $1", req.ItemCode).Scan(&itemId, &itemName)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Barang dengan kode "+req.ItemCode+" tidak ditemukan"))
		return
	}

	// Update item status to MAINTENANCE
	_, err = h.db.Exec(ctx, "UPDATE items SET status = 'MAINTENANCE', condition = 'IN_REPAIR', updated_at = NOW() WHERE id = $1", itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengupdate status barang"))
		return
	}

	// Insert maintenance log
	var logId int
	err = h.db.QueryRow(ctx,
		`INSERT INTO maintenance_logs (item_id, reported_by, issue_description, cost, vendor, notes, status)
		 VALUES ($1, $2, $3, $4, $5, $6, 'PENDING') RETURNING id`,
		itemId, userId, req.IssueDescription, req.Cost, utils.NullString(req.Vendor), utils.NullString(req.Notes),
	).Scan(&logId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal menyimpan log perbaikan: "+err.Error()))
		return
	}

	utils.LogAudit(h.db, userId.(string), "CREATE_MAINTENANCE", "MAINTENANCE", itemId, 
		fmt.Sprintf("Mencatatkan laporan perbaikan untuk barang [%s] '%s'. Masalah: %s. Biaya awal: %.2f. Vendor: %s. Melalui menu Pemeliharaan.", 
			req.ItemCode, itemName, req.IssueDescription, req.Cost, req.Vendor), c.ClientIP())

	c.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"id": logId, "item_name": itemName}, "Log perbaikan berhasil dibuat"))
}

// UpdateStatus updates the status of a maintenance log (e.g., IN_PROGRESS, DONE, CANCELLED)
func (h *MaintenanceHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req struct {
		Status string  `json:"status" binding:"required"`
		Cost   float64 `json:"cost"`
		Vendor string  `json:"vendor"`
		Notes  string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Data tidak valid"))
		return
	}

	ctx := context.Background()

	// Validate status
	validStatuses := map[string]bool{"PENDING": true, "IN_PROGRESS": true, "DONE": true, "CANCELLED": true}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Status tidak valid"))
		return
	}

	// Get current maintenance log for audit and itemId
	var currentStatus, currentVendor, currentNotes string
	var currentCost float64
	var itemId string
	err := h.db.QueryRow(ctx, "SELECT status, COALESCE(vendor, ''), COALESCE(notes, ''), COALESCE(cost, 0), item_id FROM maintenance_logs WHERE id = $1", id).
		Scan(&currentStatus, &currentVendor, &currentNotes, &currentCost, &itemId)
	
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Log perbaikan tidak ditemukan"))
		return
	}

	changes := []string{}
	if req.Status != currentStatus { changes = append(changes, fmt.Sprintf("Status: %s -> %s", currentStatus, req.Status)) }
	if req.Vendor != currentVendor && req.Vendor != "" { changes = append(changes, fmt.Sprintf("Vendor: %s -> %s", currentVendor, req.Vendor)) }
	if req.Cost != currentCost { changes = append(changes, fmt.Sprintf("Biaya: %.2f -> %.2f", currentCost, req.Cost)) }

	// Update maintenance log
	query := `UPDATE maintenance_logs SET status = $1, notes = $2, cost = $3, vendor = $4`
	args := []interface{}{req.Status, utils.NullString(req.Notes), req.Cost, utils.NullString(req.Vendor)}

	if req.Status == "DONE" || req.Status == "CANCELLED" {
		query += `, completed_at = $5 WHERE id = $6`
		args = append(args, time.Now(), id)
	} else {
		query += ` WHERE id = $5`
		args = append(args, id)
	}

	_, err = h.db.Exec(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengupdate status perbaikan"))
		return
	}

	// If DONE or CANCELLED, update item status back
	if req.Status == "DONE" {
		_, _ = h.db.Exec(ctx, "UPDATE items SET status = 'AVAILABLE', condition = 'GOOD', updated_at = NOW() WHERE id = $1", itemId)
	} else if req.Status == "CANCELLED" {
		_, _ = h.db.Exec(ctx, "UPDATE items SET status = 'AVAILABLE', updated_at = NOW() WHERE id = $1", itemId)
	}

	userId, _ := c.Get("userID")
	// Fetch item name for better log
	var itemName string
	_ = h.db.QueryRow(ctx, "SELECT name FROM items WHERE id = $1", itemId).Scan(&itemName)
	
	auditDesc := fmt.Sprintf("Memperbarui status perbaikan barang '%s'. Status: %s.", itemName, req.Status)
	if len(changes) > 0 {
		auditDesc += " Perubahan: " + strings.Join(changes, " | ")
	}
	utils.LogAudit(h.db, userId.(string), "UPDATE_MAINTENANCE", "MAINTENANCE", itemId, auditDesc, c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Status perbaikan berhasil diupdate"))
}

// Delete removes a maintenance log
func (h *MaintenanceHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "ID log tidak valid"))
		return
	}
	ctx := context.Background()

	// Get item_id first
	var itemId string
	var status string
	err = h.db.QueryRow(ctx, "SELECT item_id, status FROM maintenance_logs WHERE id = $1", id).Scan(&itemId, &status)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Log perbaikan tidak ditemukan"))
		return
	}

	_, err = h.db.Exec(ctx, "DELETE FROM maintenance_logs WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal menghapus log perbaikan"))
		return
	}

	// If the maintenance was active, set item back to AVAILABLE
	if status == "PENDING" || status == "IN_PROGRESS" {
		_, _ = h.db.Exec(ctx, "UPDATE items SET status = 'AVAILABLE', updated_at = NOW() WHERE id = $1", itemId)
	}

	actorId, _ := c.Get("userID")
	// Fetch item details for audit
	var itemCode, itemName string
	h.db.QueryRow(context.Background(), "SELECT code, name FROM items WHERE id = $1", itemId).Scan(&itemCode, &itemName)
	
	auditDesc := fmt.Sprintf("Menghapus catatan riwayat perbaikan #%d untuk barang [%s] '%s' dari sistem.", id, itemCode, itemName)
	utils.LogAudit(h.db, actorId.(string), "DELETE_MAINTENANCE", "MAINTENANCE", itemId, auditDesc, c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Log perbaikan berhasil dihapus"))
}
