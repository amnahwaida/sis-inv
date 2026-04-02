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

type LocationHandler struct {
	db *pgxpool.Pool
}

func NewLocationHandler(db *pgxpool.Pool) *LocationHandler {
	return &LocationHandler{db: db}
}

func (h *LocationHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))
	search := c.Query("search")

	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 1000 { pageSize = 12 }
	offset := (page - 1) * pageSize

	baseQuery := ` FROM locations l LEFT JOIN items i ON l.id = i.location_id WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if search != "" {
		searchStr := "%" + search + "%"
		baseQuery += fmt.Sprintf(" AND (l.name ILIKE $%d OR l.description ILIKE $%d)", argIdx, argIdx)
		args = append(args, searchStr)
		argIdx++
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(DISTINCT l.id)" + baseQuery
	err := h.db.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil { total = 0 }

	// Main query with pagination
	query := `
		SELECT l.id, l.name, l.description, l.created_at, COUNT(i.id) as item_count 
		` + baseQuery + `
		GROUP BY l.id 
		ORDER BY l.name ASC
		LIMIT $` + fmt.Sprint(argIdx) + ` OFFSET $` + fmt.Sprint(argIdx+1)
	
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch locations"))
		return
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var loc models.Location
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Description, &loc.CreatedAt, &loc.ItemCount); err != nil {
			continue
		}
		locations = append(locations, loc)
	}

	if locations == nil {
		locations = []models.Location{}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": locations,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}, ""))
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

	actorId, _ := c.Get("userID")
	auditDesc := fmt.Sprintf("Menambahkan lokasi penyimpanan baru: '%s'. Melalui menu kustomisasi branding/pengaturan.", req.Name)
	if req.Description != nil && *req.Description != "" {
		auditDesc += fmt.Sprintf(" Deskripsi: %s.", *req.Description)
	}
	utils.LogAudit(h.db, actorId.(string), "CREATE_LOCATION", "LOCATION", fmt.Sprint(id), auditDesc, c.ClientIP())

	c.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"id": id}, "Location created successfully"))
}

func (h *LocationHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid data"))
		return
	}

	_, err := h.db.Exec(context.Background(),
		"UPDATE locations SET name = $1, description = $2 WHERE id = $3",
		req.Name, req.Description, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update location"))
		return
	}

	actorId, _ := c.Get("userID")
	utils.LogAudit(h.db, actorId.(string), "UPDATE_LOCATION", "LOCATION", id, fmt.Sprintf("Memperbarui data lokasi penyimpanan: '%s'. Melalui menu kustomisasi.", req.Name), c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Location updated successfully"))
}

func (h *LocationHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "ID lokasi tidak valid"))
		return
	}

	// Check if any items are using this location
	// We check both location_id (foreign key) and location (string name backup)
	var count int
	err = h.db.QueryRow(context.Background(), 
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

	actorId, _ := c.Get("userID")
	// Fetch name for audit log
	var locName string
	_ = h.db.QueryRow(context.Background(), "SELECT name FROM locations WHERE id = $1", id).Scan(&locName)
	utils.LogAudit(h.db, actorId.(string), "DELETE_LOCATION", "LOCATION", idStr, fmt.Sprintf("Menghapus lokasi penyimpanan: '%s' dari sistem.", locName), c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Location deleted successfully"))
}
