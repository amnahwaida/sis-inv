package handlers

import (
	"context"
	"net/http"

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
	rows, err := h.db.Query(context.Background(), "SELECT id, name, description, color_code, created_at FROM categories ORDER BY name ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch categories"))
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.ColorCode, &cat.CreatedAt); err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	if categories == nil {
		categories = []models.Category{}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(categories, ""))
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

	c.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"id": id}, "Category created successfully"))
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Check if any items are using this category
	var count int
	err := h.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM items WHERE category_id = $1", id).Scan(&count)
	if err == nil && count > 0 {
		c.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, "Cannot delete category being used by items"))
		return
	}

	_, err = h.db.Exec(context.Background(), "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to delete category"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Category deleted successfully"))
}
