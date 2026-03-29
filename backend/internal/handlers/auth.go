package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/config"
	"github.com/vannyezha/sis-inv/internal/models"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type AuthHandler struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

func NewAuthHandler(db *pgxpool.Pool, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Username dan password diperlukan"))
		return
	}

	var user models.User
	err := h.db.QueryRow(context.Background(),
		`SELECT id, username, password_hash, full_name, role, nip, email, phone, is_active, created_at, updated_at
		 FROM users WHERE username = $1`,
		req.Username,
	).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.FullName, &user.Role,
		&user.NIP, &user.Email, &user.Phone, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Username atau password salah"))
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, "Akun tidak aktif. Hubungi admin."))
		return
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Username atau password salah"))
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Role, h.cfg.JWTAccessExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal generate token"))
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, h.cfg.JWTRefreshExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal generate refresh token"))
		return
	}

	// Update last login
	_, _ = h.db.Exec(context.Background(),
		"UPDATE users SET last_login = $1 WHERE id = $2",
		time.Now(), user.ID,
	)

	c.JSON(http.StatusOK, utils.SuccessResponse(models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, "Login berhasil"))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT is stateless, client handles token removal
	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Logout berhasil"))
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user models.User
	err := h.db.QueryRow(context.Background(),
		`SELECT id, username, full_name, role, nip, email, phone, is_active, last_login, created_at, updated_at
		 FROM users WHERE id = $1`,
		userID,
	).Scan(
		&user.ID, &user.Username, &user.FullName, &user.Role,
		&user.NIP, &user.Email, &user.Phone, &user.IsActive, &user.LastLogin,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "User tidak ditemukan"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(user, ""))
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Data tidak valid"))
		return
	}

	var currentHash string
	err := h.db.QueryRow(context.Background(),
		"SELECT password_hash FROM users WHERE id = $1", userID,
	).Scan(&currentHash)

	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "User tidak ditemukan"))
		return
	}

	if !utils.CheckPassword(req.OldPassword, currentHash) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Password lama salah"))
		return
	}

	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal hash password"))
		return
	}

	_, err = h.db.Exec(context.Background(),
		"UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3",
		newHash, time.Now(), userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal update password"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Password berhasil diubah"))
}
