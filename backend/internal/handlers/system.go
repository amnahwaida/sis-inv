package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type SystemHandler struct {
	db *pgxpool.Pool
}

func NewSystemHandler(db *pgxpool.Pool) *SystemHandler {
	return &SystemHandler{db: db}
}

func (h *SystemHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"status":  "ok",
		"service": "SIS-INV API",
		"version": "1.0.0",
	}, "Server berjalan"))
}
