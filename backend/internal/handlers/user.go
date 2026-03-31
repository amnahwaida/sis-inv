package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/models"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type UserHandler struct {
	db *pgxpool.Pool
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) List(c *gin.Context) {
	rows, err := h.db.Query(context.Background(),
		`SELECT id, username, full_name, role, nip, email, phone, is_active, last_login, created_at, updated_at
		 FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data user"))
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.FullName, &u.Role,
			&u.NIP, &u.Email, &u.Phone, &u.IsActive, &u.LastLogin,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			continue
		}
		users = append(users, u)
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(users, ""))
}

func (h *UserHandler) Create(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Data tidak valid: "+err.Error()))
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal hash password"))
		return
	}

	var userID string
	err = h.db.QueryRow(context.Background(),
		`INSERT INTO users (username, password_hash, full_name, role, nip, email, phone)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		req.Username, hash, req.FullName, req.Role, req.NIP, req.Email, req.Phone,
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, "Username sudah digunakan"))
		return
	}

	actorId, _ := c.Get("userID")
	utils.LogAudit(h.db, actorId.(string), "CREATE_USER", "USER", userID, "Created user: "+req.Username, c.ClientIP())

	c.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"id": userID}, "User berhasil dibuat"))
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Data tidak valid"))
		return
	}

	query := "UPDATE users SET updated_at = $1"
	args := []interface{}{time.Now()}
	argIdx := 2

	if req.FullName != nil {
		query += fmt.Sprintf(", full_name = $%d", argIdx)
		args = append(args, *req.FullName)
		argIdx++
	}
	if req.Role != nil {
		query += fmt.Sprintf(", role = $%d", argIdx)
		args = append(args, *req.Role)
		argIdx++
	}
	if req.NIP != nil {
		query += fmt.Sprintf(", nip = $%d", argIdx)
		args = append(args, *req.NIP)
		argIdx++
	}
	if req.Email != nil {
		query += fmt.Sprintf(", email = $%d", argIdx)
		args = append(args, *req.Email)
		argIdx++
	}
	if req.Phone != nil {
		query += fmt.Sprintf(", phone = $%d", argIdx)
		args = append(args, *req.Phone)
		argIdx++
	}
	if req.IsActive != nil {
		query += fmt.Sprintf(", is_active = $%d", argIdx)
		args = append(args, *req.IsActive)
		argIdx++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	result, err := h.db.Exec(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal update user"))
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "User tidak ditemukan"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "User berhasil diupdate"))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Soft delete - set is_active to false
	result, err := h.db.Exec(context.Background(),
		"UPDATE users SET is_active = false, updated_at = $1 WHERE id = $2", time.Now(), id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal menghapus user"))
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "User tidak ditemukan"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "User berhasil dinonaktifkan"))
}
