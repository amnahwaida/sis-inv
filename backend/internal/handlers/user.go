package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"strconv"

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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	role := c.Query("role")

	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	offset := (page - 1) * pageSize

	query := `SELECT id, username, full_name, role, nip, email, phone, is_active, last_login, created_at, updated_at FROM users WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM users WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if search != "" {
		searchStr := "%" + search + "%"
		query += fmt.Sprintf(" AND (username ILIKE $%d OR full_name ILIKE $%d)", argIdx, argIdx)
		countQuery += fmt.Sprintf(" AND (username ILIKE $%d OR full_name ILIKE $%d)", argIdx, argIdx)
		args = append(args, searchStr)
		argIdx++
	}

	if role != "" {
		query += fmt.Sprintf(" AND role = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, role)
		argIdx++
	}

	// Get total count
	var total int
	err := h.db.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil { total = 0 }

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(context.Background(), query, args...)
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

	if users == nil { users = []models.User{} }

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": users,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}, ""))
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
	auditDesc := fmt.Sprintf("Menambahkan pengguna baru: %s (Username: %s, Role: %s, NIP: %s). Melalui panel administrasi.", 
		req.FullName, req.Username, req.Role, req.NIP)
	utils.LogAudit(h.db, actorId.(string), "CREATE_USER", "USER", userID, auditDesc, c.ClientIP())

	c.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"id": userID}, "User berhasil dibuat"))
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Data tidak valid"))
		return
	}

	ctx := context.Background()

	// 1. Fetch current user data for audit
	var oldUsername, oldName, oldRole, oldNip, oldEmail, oldPhone string
	var oldActive bool
	err := h.db.QueryRow(ctx, "SELECT username, full_name, role, COALESCE(nip, ''), COALESCE(email, ''), COALESCE(phone, ''), is_active FROM users WHERE id = $1", id).
		Scan(&oldUsername, &oldName, &oldRole, &oldNip, &oldEmail, &oldPhone, &oldActive)

	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "User tidak ditemukan"))
		return
	}

	// Validate if user is disabling the last active admin
	if req.IsActive != nil && !*req.IsActive && oldActive && oldRole == "ADMIN" {
		var activeAdminCount int
		err := h.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE role = 'ADMIN' AND is_active = true").Scan(&activeAdminCount)
		if err == nil && activeAdminCount <= 1 {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Pencegahan Sistem: Minimal harus ada 1 akun Admin yang aktif"))
			return
		}
	}

	changes := []string{}
	query := "UPDATE users SET updated_at = $1"
	args := []interface{}{time.Now()}
	argIdx := 2

	if req.Username != nil {
		query += fmt.Sprintf(", username = $%d", argIdx)
		args = append(args, *req.Username)
		argIdx++
		if *req.Username != oldUsername { changes = append(changes, fmt.Sprintf("Username: %s -> %s", oldUsername, *req.Username)) }
	}
	if req.FullName != nil {
		query += fmt.Sprintf(", full_name = $%d", argIdx)
		args = append(args, *req.FullName)
		argIdx++
		if *req.FullName != oldName { changes = append(changes, fmt.Sprintf("Name: %s -> %s", oldName, *req.FullName)) }
	}
	if req.Role != nil {
		query += fmt.Sprintf(", role = $%d", argIdx)
		args = append(args, *req.Role)
		argIdx++
		if *req.Role != oldRole { changes = append(changes, fmt.Sprintf("Role: %s -> %s", oldRole, *req.Role)) }
	}
	if req.NIP != nil {
		query += fmt.Sprintf(", nip = $%d", argIdx)
		args = append(args, *req.NIP)
		argIdx++
		if *req.NIP != oldNip { changes = append(changes, fmt.Sprintf("NIP: %s -> %s", oldNip, *req.NIP)) }
	}
	if req.Email != nil {
		query += fmt.Sprintf(", email = $%d", argIdx)
		args = append(args, *req.Email)
		argIdx++
		if *req.Email != oldEmail { changes = append(changes, fmt.Sprintf("Email: %s -> %s", oldEmail, *req.Email)) }
	}
	if req.Phone != nil {
		query += fmt.Sprintf(", phone = $%d", argIdx)
		args = append(args, *req.Phone)
		argIdx++
		if *req.Phone != oldPhone { changes = append(changes, fmt.Sprintf("Phone: %s -> %s", oldPhone, *req.Phone)) }
	}
	if req.IsActive != nil {
		query += fmt.Sprintf(", is_active = $%d", argIdx)
		args = append(args, *req.IsActive)
		argIdx++
		if *req.IsActive != oldActive { changes = append(changes, fmt.Sprintf("Active: %v -> %v", oldActive, *req.IsActive)) }
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

	actorId, _ := c.Get("userID")
	auditDesc := "Updated user profile/status for: " + oldName
	if len(changes) > 0 {
		auditDesc = fmt.Sprintf("Updated user %s. Changes: %s", oldName, strings.Join(changes, " | "))
	}
	utils.LogAudit(h.db, actorId.(string), "UPDATE_USER", "USER", id, auditDesc, c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "User berhasil diupdate"))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" || id == "undefined" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "ID user tidak valid"))
		return
	}

	// Soft delete - set is_active to false
	// First check if target is an active admin
	var targetRole string
	var targetActive bool
	var targetName string
	err := h.db.QueryRow(context.Background(), "SELECT full_name, role, is_active FROM users WHERE id = $1", id).Scan(&targetName, &targetRole, &targetActive)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "User tidak ditemukan"))
		return
	}

	if targetRole == "ADMIN" && targetActive {
		var activeAdminCount int
		err := h.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE role = 'ADMIN' AND is_active = true").Scan(&activeAdminCount)
		if err == nil && activeAdminCount <= 1 {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Pencegahan Sistem: Minimal harus ada 1 akun Admin yang aktif"))
			return
		}
	}

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

	if targetName == "" {
		targetName = id
	}

	actorId, _ := c.Get("userID")
	auditDesc := fmt.Sprintf("Menonaktifkan akses login pengguna: %s (Username: %s). Akun tidak lagi dapat digunakan untuk masuk ke sistem.", 
		targetName, id)
	utils.LogAudit(h.db, actorId.(string), "DISABLE_USER", "USER", id, auditDesc, c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "User berhasil dinonaktifkan"))
}
