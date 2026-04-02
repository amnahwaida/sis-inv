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

type CategoryHandler struct {
	db *pgxpool.Pool
}

func NewCategoryHandler(db *pgxpool.Pool) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))
	search := c.Query("search")

	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 1000 { pageSize = 12 }
	offset := (page - 1) * pageSize

	baseQuery := ` FROM categories c LEFT JOIN items i ON c.id = i.category_id WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if search != "" {
		searchStr := "%" + search + "%"
		baseQuery += fmt.Sprintf(" AND (c.name ILIKE $%d OR c.description ILIKE $%d)", argIdx, argIdx)
		args = append(args, searchStr)
		argIdx++
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(DISTINCT c.id)" + baseQuery
	err := h.db.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil { total = 0 }

	// Main query with pagination
	query := `
		SELECT c.id, c.name, c.description, c.color_code, c.created_at, COUNT(i.id) as item_count 
		` + baseQuery + `
		GROUP BY c.id 
		ORDER BY c.name ASC
		LIMIT $` + fmt.Sprint(argIdx) + ` OFFSET $` + fmt.Sprint(argIdx+1)
	
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch categories"))
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.ColorCode, &cat.CreatedAt, &cat.ItemCount); err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	if categories == nil {
		categories = []models.Category{}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": categories,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}, ""))
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description"`
		ColorCode   *string `json:"color_code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid data"))
		return
	}

	var id int
	err := h.db.QueryRow(context.Background(), 
		"INSERT INTO categories (name, description, color_code) VALUES ($1, $2, $3) RETURNING id",
		req.Name, req.Description, req.ColorCode,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, "Category already exists or invalid data"))
		return
	}

	actorId, _ := c.Get("userID")
	auditDesc := fmt.Sprintf("Menambahkan kategori barang baru: '%s'. Melalui menu kustomisasi branding/pengaturan.", req.Name)
	if req.Description != nil && *req.Description != "" {
		auditDesc += fmt.Sprintf(" Deskripsi: %s.", *req.Description)
	}
	utils.LogAudit(h.db, actorId.(string), "CREATE_CATEGORY", "CATEGORY", fmt.Sprint(id), auditDesc, c.ClientIP())

	c.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"id": id}, "Category created successfully"))
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description"`
		ColorCode   *string `json:"color_code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid data"))
		return
	}

	_, err := h.db.Exec(context.Background(),
		"UPDATE categories SET name = $1, description = $2, color_code = $3 WHERE id = $4",
		req.Name, req.Description, req.ColorCode, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update category"))
		return
	}

	actorId, _ := c.Get("userID")
	utils.LogAudit(h.db, actorId.(string), "UPDATE_CATEGORY", "CATEGORY", id, fmt.Sprintf("Memperbarui data kategori barang: '%s'. Melalui menu kustomisasi.", req.Name), c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Category updated successfully"))
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "ID kategori tidak valid"))
		return
	}

	// Check if any items are using this category
	var count int
	err = h.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM items WHERE category_id = $1", id).Scan(&count)
	if err == nil && count > 0 {
		c.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, "Cannot delete category being used by items"))
		return
	}

	_, err = h.db.Exec(context.Background(), "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to delete category"))
		return
	}

	actorId, _ := c.Get("userID")
	// Fetch name for audit log
	var catName string
	_ = h.db.QueryRow(context.Background(), "SELECT name FROM categories WHERE id = $1", id).Scan(&catName)
	utils.LogAudit(h.db, actorId.(string), "DELETE_CATEGORY", "CATEGORY", idStr, fmt.Sprintf("Menghapus kategori barang: '%s' dari sistem.", catName), c.ClientIP())
	
	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Category deleted successfully"))
}
