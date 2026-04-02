package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PWAHandler struct {
	db *pgxpool.Pool
}

func NewPWAHandler(db *pgxpool.Pool) *PWAHandler {
	return &PWAHandler{db: db}
}

func (h *PWAHandler) GetManifest(c *gin.Context) {
	// Fetch necessary settings
	rows, err := h.db.Query(c, "SELECT key, value FROM settings WHERE key IN ('app_name', 'app_theme', 'app_pwa_icon_192', 'app_pwa_icon_512', 'app_description')")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err == nil {
			settings[key] = value
		}
	}

	// Default values
	appName := settings["app_name"]
	if appName == "" {
		appName = "SIS-INV"
	}
	appDesc := settings["app_description"]
	if appDesc == "" {
		appDesc = "Sistem Manajemen Inventaris Sekolah"
	}
	icon192 := settings["app_pwa_icon_192"]
	if icon192 == "" {
		icon192 = "/pwa-192x192.png"
	}
	icon512 := settings["app_pwa_icon_512"]
	if icon512 == "" {
		icon512 = "/pwa-512x512.png"
	}

	// Dynamic theme color (optional, can be mapped from app_theme)
	themeColor := "#4F46E5" // Default Indigo 600
	switch settings["app_theme"] {
	case "green":
		themeColor = "#10B981"
	case "purple":
		themeColor = "#8B5CF6"
	case "red":
		themeColor = "#EF4444"
	case "orange":
		themeColor = "#F59E0B"
	}

	manifest := gin.H{
		"name":             appName,
		"short_name":        appName,
		"description":       appDesc,
		"theme_color":       themeColor,
		"background_color":  "#ffffff",
		"display":           "standalone",
		"start_url":         "/",
		"icons": []interface{}{
			gin.H{
				"src":   icon192,
				"sizes": "192x192",
				"type":  "image/png",
			},
			gin.H{
				"src":   icon512,
				"sizes": "512x512",
				"type":  "image/png",
			},
			gin.H{
				"src":     icon512,
				"sizes":   "512x512",
				"type":    "image/png",
				"purpose": "any maskable",
			},
		},
	}

	c.JSON(http.StatusOK, manifest)
}
