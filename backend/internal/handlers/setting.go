package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/models"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type SettingHandler struct {
	db *pgxpool.Pool
}

func NewSettingHandler(db *pgxpool.Pool) *SettingHandler {
	return &SettingHandler{db: db}
}

func (h *SettingHandler) List(c *gin.Context) {
	rows, err := h.db.Query(context.Background(), "SELECT key, value, updated_at FROM settings")
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch settings"))
		return
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var s models.Setting
		if err := rows.Scan(&s.Key, &s.Value, &s.UpdatedAt); err != nil {
			continue
		}
		settings[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(settings, ""))
}

func (h *SettingHandler) Update(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid data"))
		return
	}

	tx, err := h.db.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to start transaction"))
		return
	}
	defer tx.Rollback(context.Background())

	for key, value := range req {
		_, err := tx.Exec(context.Background(), 
			"INSERT INTO settings (key, value, updated_at) VALUES ($1, $2, CURRENT_TIMESTAMP) ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = CURRENT_TIMESTAMP",
			key, value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update setting: "+key))
			return
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to commit transaction"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Settings updated successfully"))
}
