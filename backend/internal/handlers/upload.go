package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("❌ Upload error: %v\n", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "No file uploaded"))
		return
	}

	// Create uploads directory if not exists
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadDir, 0755); err != nil {
			fmt.Printf("❌ Failed to create directory: %v\n", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Storage initialization failed"))
			return
		}
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)
	dst := filepath.Join(uploadDir, filename)

	// Save the file
	if err := c.SaveUploadedFile(file, dst); err != nil {
		fmt.Printf("❌ Save error: %v\n", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to save file on server"))
		return
	}

	// Return the relative URL
	photoURL := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"url": photoURL,
	}, "File uploaded successfully"))
}
