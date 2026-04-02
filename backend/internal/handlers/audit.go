package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/models"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type AuditHandler struct {
	db *pgxpool.Pool
}

func NewAuditHandler(db *pgxpool.Pool) *AuditHandler {
	return &AuditHandler{db: db}
}

// StartSession creates a new audit session
func (h *AuditHandler) StartSession(c *gin.Context) {
	var req models.CreateAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body"))
		return
	}

	userId, _ := c.Get("userID")
	ctx := context.Background()

	var sessionId string
	err := h.db.QueryRow(ctx,
		`INSERT INTO audit_sessions (location_id, user_id, notes, status)
		 VALUES ($1, $2, $3, 'OPEN') RETURNING id`,
		req.LocationID, userId, utils.NullString(req.Notes),
	).Scan(&sessionId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create audit session: "+err.Error()))
		return
	}

	utils.LogAudit(h.db, userId.(string), "START_AUDIT", "AUDIT_SESSION", sessionId, "Started audit session for location ID: "+fmt.Sprint(req.LocationID), c.ClientIP())

	c.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"id": sessionId}, "Audit session started"))
}

// ScanItem registers an item as found in a session
func (h *AuditHandler) ScanItem(c *gin.Context) {
	sessionId := c.Param("id")
	var req models.ScanAuditItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body"))
		return
	}

	ctx := context.Background()

	// Find item id by code
	var itemId string
	err := h.db.QueryRow(ctx, "SELECT id FROM items WHERE code = $1", req.ItemCode).Scan(&itemId)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Item not found"))
		return
	}

	// Update item condition and updated_at
	_, err = h.db.Exec(ctx, "UPDATE items SET condition = $1, updated_at = NOW() WHERE id = $2", req.Condition, itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update item condition"))
		return
	}

	// Upsert audit item record based on (session_id, item_id)
	// We check if it exists first because we don't have a unique constraint on (session_id, item_id) yet
	var existingId string
	err = h.db.QueryRow(ctx, "SELECT id FROM audit_items WHERE session_id = $1 AND item_id = $2", sessionId, itemId).Scan(&existingId)
	
	if err != nil {
		// New scan
		_, err = h.db.Exec(ctx,
			`INSERT INTO audit_items (session_id, item_id, found_condition, notes)
			 VALUES ($1, $2, $3, $4)`,
			sessionId, itemId, req.Condition, utils.NullString(req.Notes),
		)
	} else {
		// Update existing scan
		_, err = h.db.Exec(ctx,
			`UPDATE audit_items SET found_condition = $1, notes = $2, scanned_at = NOW()
			 WHERE id = $3`,
			req.Condition, utils.NullString(req.Notes), existingId,
		)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to record audited item: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Item scanned successfully"))
}

// CloseSession closes an audit session
func (h *AuditHandler) CloseSession(c *gin.Context) {
	sessionId := c.Param("id")
	ctx := context.Background()

	_, err := h.db.Exec(ctx,
		"UPDATE audit_sessions SET status = 'CLOSED', finished_at = NOW(), updated_at = NOW() WHERE id = $1",
		sessionId,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to close audit session"))
		return
	}

	userId, _ := c.Get("userID")
	utils.LogAudit(h.db, userId.(string), "CLOSE_AUDIT", "AUDIT_SESSION", sessionId, "Closed audit session", c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Audit session closed"))
}

// ListSessions returns a list of audit sessions
func (h *AuditHandler) ListSessions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")

	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	offset := (page - 1) * pageSize

	baseQuery := `
		FROM audit_sessions s
		JOIN locations l ON s.location_id = l.id
		JOIN users u ON s.user_id = u.id
		WHERE 1=1
	`
	var args []interface{}
	argIdx := 1

	if search != "" {
		searchStr := "%" + search + "%"
		baseQuery += fmt.Sprintf(" AND (l.name ILIKE $%d OR u.full_name ILIKE $%d OR s.notes ILIKE $%d)", argIdx, argIdx, argIdx)
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
		SELECT s.id, s.location_id, l.name as location_name, s.user_id, u.full_name as user_name,
		       s.status, s.started_at, s.finished_at, s.notes, s.created_at,
		       (SELECT COUNT(*) FROM items WHERE location_id = l.id) as total_expected,
		       (SELECT COUNT(*) FROM audit_items WHERE session_id = s.id) as total_found
		` + baseQuery + `
		ORDER BY s.created_at DESC
		LIMIT $` + fmt.Sprint(argIdx) + ` OFFSET $` + fmt.Sprint(argIdx+1)
	
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch audit sessions"))
		return
	}
	defer rows.Close()

	sessions := []models.AuditSession{}
	for rows.Next() {
		var s models.AuditSession
		err := rows.Scan(
			&s.ID, &s.LocationID, &s.LocationName, &s.UserID, &s.UserName,
			&s.Status, &s.StartedAt, &s.FinishedAt, &s.Notes, &s.CreatedAt,
			&s.TotalExpected, &s.TotalFound,
		)
		if err != nil {
			fmt.Printf("⚠️ Scan error in ListSessions: %v\n", err)
			continue
		}
		sessions = append(sessions, s)
	}

	if sessions == nil { sessions = []models.AuditSession{} }

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": sessions,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}, ""))
}

// GetSessionDetail returns detailed report of a session
func (h *AuditHandler) GetSessionDetail(c *gin.Context) {
	sessionId := c.Param("id")
	ctx := context.Background()

	// 1. Get Session Info
	var s models.AuditSession
	err := h.db.QueryRow(ctx,
		`SELECT s.id, s.location_id, l.name as location_name, s.user_id, u.full_name as user_name,
		        s.status, s.started_at, s.finished_at, s.notes, s.created_at,
		        (SELECT COUNT(*) FROM items WHERE location_id = l.id) as total_expected,
		        (SELECT COUNT(*) FROM audit_items WHERE session_id = s.id) as total_found
		 FROM audit_sessions s
		 JOIN locations l ON s.location_id = l.id
		 JOIN users u ON s.user_id = u.id
		 WHERE s.id = $1`, sessionId).Scan(
		&s.ID, &s.LocationID, &s.LocationName, &s.UserID, &s.UserName,
		&s.Status, &s.StartedAt, &s.FinishedAt, &s.Notes, &s.CreatedAt,
		&s.TotalExpected, &s.TotalFound,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Session not found"))
		return
	}

	// 2. Get Found Items
	foundRows, err := h.db.Query(ctx,
		`SELECT ai.item_id, i.code, i.name, ai.found_condition, ai.scanned_at, ai.notes
		 FROM audit_items ai
		 JOIN items i ON ai.item_id = i.id
		 WHERE ai.session_id = $1`, sessionId)
	
	foundItems := []models.AuditItem{}
	if err == nil {
		defer foundRows.Close()
		for foundRows.Next() {
			var ai models.AuditItem
			ai.SessionID = sessionId
			foundRows.Scan(&ai.ItemID, &ai.ItemCode, &ai.ItemName, &ai.FoundCondition, &ai.ScannedAt, &ai.Notes)
			foundItems = append(foundItems, ai)
		}
	}

	// 3. Get Missing Items (Items in location but not in audit_items)
	missingRows, err := h.db.Query(ctx,
		`SELECT i.id, i.code, i.name, i.condition
		 FROM items i
		 WHERE i.location_id = $1
		 AND i.id NOT IN (SELECT item_id FROM audit_items WHERE session_id = $2)`,
		s.LocationID, sessionId)
	
	missingItems := []map[string]interface{}{}
	if err == nil {
		defer missingRows.Close()
		for missingRows.Next() {
			var id, code, name, condition string
			missingRows.Scan(&id, &code, &name, &condition)
			missingItems = append(missingItems, map[string]interface{}{
				"id": id, "code": code, "name": name, "condition": condition,
			})
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"session": s,
		"found_items": foundItems,
		"missing_items": missingItems,
	}, ""))
}
