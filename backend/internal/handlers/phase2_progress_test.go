package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/vannyezha/sis-inv/internal/config"
	"github.com/vannyezha/sis-inv/internal/database"
	"github.com/vannyezha/sis-inv/internal/utils"
)

func setupTestRouter() (*gin.Engine, *config.Config) {
	gin.SetMode(gin.TestMode)
	_ = godotenv.Load("../../.env")
	cfg := config.Load()
	utils.InitJWT(cfg.JWTSecret)
	return gin.Default(), cfg
}

func TestPhase2Progress(t *testing.T) {
	router, cfg := setupTestRouter()
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Initialize Handlers
	itemHandler := NewItemHandler(db)
	dashHandler := NewDashboardHandler(db)
	userHandler := NewUserHandler(db)

	// Mock Admin Token
	token, _ := utils.GenerateAccessToken("test-admin-id", "admin", "ADMIN", time.Hour)

	t.Run("Check Dashboard Summary (Phase 2)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/dashboard/summary", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		w := httptest.NewRecorder()
		router.GET("/api/v1/dashboard/summary", dashHandler.GetSummary)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
	})

	t.Run("Check Items List & History (Phase 2)", func(t *testing.T) {
		// 1. Check List
		req, _ := http.NewRequest("GET", "/api/v1/items", nil)
		w := httptest.NewRecorder()
		router.GET("/api/v1/items", itemHandler.List)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// 2. Check History Endpoint (Added in Phase 2)
		// Assuming we have at least one item or using a dummy ID
		reqHist, _ := http.NewRequest("GET", "/api/v1/items/any-id/history", nil)
		wHist := httptest.NewRecorder()
		router.GET("/api/v1/items/:id/history", itemHandler.GetHistory)
		router.ServeHTTP(wHist, reqHist)
		
		// Should be 200 even if empty (returns []), or 404 if ID doesn't exist but we return 200 [] for any ID usually if no rows found
		assert.Equal(t, http.StatusOK, wHist.Code)
	})

	t.Run("Check User Management CRUD (Phase 2)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users", nil)
		w := httptest.NewRecorder()
		router.GET("/api/v1/users", userHandler.List)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
