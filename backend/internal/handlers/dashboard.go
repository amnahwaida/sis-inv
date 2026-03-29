package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type DashboardHandler struct {
	db *pgxpool.Pool
}

func NewDashboardHandler(db *pgxpool.Pool) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// GetSummary returns counts of items by status and condition, as well as active borrowing count
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	ctx := context.Background()

	// 1. Total Items Count
	var totalItems int
	err := h.db.QueryRow(ctx, "SELECT COUNT(*) FROM items").Scan(&totalItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to count total items"))
		return
	}

	// 2. Counts Grouped by Status
	statusQuery := "SELECT status, COUNT(*) FROM items GROUP BY status"
	statusRows, err := h.db.Query(ctx, statusQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to count items by status"))
		return
	}
	defer statusRows.Close()

	itemsByStatus := map[string]int{
		"AVAILABLE":   0,
		"BORROWED":    0,
		"MAINTENANCE": 0,
		"LOST":        0,
	}
	for statusRows.Next() {
		var status string
		var count int
		if err := statusRows.Scan(&status, &count); err == nil {
			itemsByStatus[status] = count
		}
	}

	// 3. Counts Grouped by Condition
	conditionQuery := "SELECT condition, COUNT(*) FROM items GROUP BY condition"
	conditionRows, err := h.db.Query(ctx, conditionQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to count items by condition"))
		return
	}
	defer conditionRows.Close()

	itemsByCondition := map[string]int{
		"GOOD":      0,
		"DAMAGED":   0,
		"LOST":      0,
		"IN_REPAIR": 0,
	}
	for conditionRows.Next() {
		var cond string
		var count int
		if err := conditionRows.Scan(&cond, &count); err == nil {
			itemsByCondition[cond] = count
		}
	}

	// 4. Overdue transactions count
	var overdueCount int
	err = h.db.QueryRow(ctx, "SELECT COUNT(*) FROM transactions WHERE status = 'ACTIVE' AND due_date < NOW()").Scan(&overdueCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to count overdue transactions"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"total_items":        totalItems,
		"items_by_status":    itemsByStatus,
		"items_by_condition": itemsByCondition,
		"overdue_count":      overdueCount,
	}, "Dashboard summary fetched successfully"))
}
