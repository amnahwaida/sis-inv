package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type StudentHandler struct {
	db *pgxpool.Pool
}

func NewStudentHandler(db *pgxpool.Pool) *StudentHandler {
	return &StudentHandler{db: db}
}

func (h *StudentHandler) List(c *gin.Context) {
	query := `SELECT id, nis, full_name, class, is_active, created_at FROM students ORDER BY full_name ASC`
	rows, err := h.db.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data siswa"))
		return
	}
	defer rows.Close()

	var students []map[string]interface{}
	for rows.Next() {
		var id int
		var nis, name, class string
		var isActive bool
		var createdAt interface{}
		if err := rows.Scan(&id, &nis, &name, &class, &isActive, &createdAt); err == nil {
			students = append(students, map[string]interface{}{
				"id":         id,
				"nis":        nis,
				"full_name":  name,
				"class":      class,
				"is_active":  isActive,
				"created_at": createdAt,
			})
		}
	}
	
	if students == nil { students = []map[string]interface{}{} }
	c.JSON(http.StatusOK, utils.SuccessResponse(students, ""))
}

func (h *StudentHandler) Create(c *gin.Context) {
	var req struct {
		NIS      string `json:"nis" binding:"required"`
		FullName string `json:"full_name" binding:"required"`
		Class    string `json:"class"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Data tidak valid"))
		return
	}

	_, err := h.db.Exec(context.Background(), 
		`INSERT INTO students (nis, full_name, class) VALUES ($1, $2, $3)`,
		req.NIS, req.FullName, req.Class)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal menyimpan data: NIS mungkin sudah terpakai"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse(nil, "Siswa berhasil ditambahkan"))
}

func (h *StudentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		NIS      string `json:"nis" binding:"required"`
		FullName string `json:"full_name" binding:"required"`
		Class    string `json:"class"`
		IsActive bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Data tidak valid"))
		return
	}

	_, err := h.db.Exec(context.Background(),
		`UPDATE students SET nis = $1, full_name = $2, class = $3, is_active = $4, updated_at = NOW() WHERE id = $5`,
		req.NIS, req.FullName, req.Class, req.IsActive, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal memperbarui data"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Siswa berhasil diperbarui"))
}

func (h *StudentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Exec(context.Background(), "DELETE FROM students WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal menghapus data: Siswa mungkin sudah memiliki transaksi"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Siswa berhasil dihapus"))
}

func (h *StudentHandler) Search(c *gin.Context) {
	queryParam := c.Query("q")
	if queryParam == "" {
		c.JSON(http.StatusOK, utils.SuccessResponse([]interface{}{}, ""))
		return
	}

	searchTerm := "%" + queryParam + "%"
	
	query := `
		SELECT nis, full_name, class 
		FROM students 
		WHERE (nis ILIKE $1 OR full_name ILIKE $1) AND is_active = true
		LIMIT 10
	`

	rows, err := h.db.Query(context.Background(), query, searchTerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mencari data siswa"))
		return
	}
	defer rows.Close()

	var students []map[string]string
	for rows.Next() {
		var nis, name, class string
		if err := rows.Scan(&nis, &name, &class); err == nil {
			students = append(students, map[string]string{
				"nis":       nis,
				"full_name": name,
				"class":     class,
			})
		}
	}

	if students == nil {
		students = []map[string]string{}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(students, ""))
}
