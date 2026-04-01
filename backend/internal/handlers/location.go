package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/models"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type LocationHandler struct {
	db *pgxpool.Pool
}

func NewLocationHandler(db *pgxpool.Pool) *LocationHandler {
	return &LocationHandler{db: db}
}

func (h *LocationHandler) List(c *gin.Context) {
	rows, err := h.db.Query(context.Background(), "SELECT id, name, description, created_at FROM locations ORDER BY name ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch locations"))
		return
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var loc models.Location
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Description, &loc.CreatedAt); err != nil {
			continue
		}
		locations = append(locations, loc)
	}

	if locations == nil {
		locations = []models.Location{}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(locations, ""))
}

func (h *LocationHandler) Create(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid data"))
		return
	}

	var id int
	err := h.db.QueryRow(context.Background(), 
		"INSERT INTO locations (name, description) VALUES ($1, $2) RETURNING id",
		req.Name, req.Description,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, "Location already exists or invalid data"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"id": id}, "Location created successfully"))
}

func (h *LocationHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	// Check if any items are using this location
	// We check both location_id (foreign key) and location (string name backup)
	var count int
	err := h.db.QueryRow(context.Background(), 
		"SELECT COUNT(*) FROM items WHERE location_id = $1 OR location = (SELECT name FROM locations WHERE id = $1)", 
		id).Scan(&count)
	if err == nil && count > 0 {
		c.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, "Cannot delete location being used by items"))
		return
	}

	_, err = h.db.Exec(context.Background(), "DELETE FROM locations WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to delete location"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Location deleted successfully"))
}
