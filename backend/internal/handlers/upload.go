package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

	// Compress if it is an image
	extLow := strings.ToLower(ext)
	if extLow == ".jpg" || extLow == ".jpeg" || extLow == ".png" {
		tempDst := dst + ".compressed"
		// Compress with quality 75 for JPEG, PNG will use default compression
		if err := utils.CompressImage(dst, tempDst, 75); err == nil {
			// If compression was successful and generated a file, replace the original
			if info, err := os.Stat(tempDst); err == nil && info.Size() > 0 {
				// Only replace if the compressed file is actually smaller (safeguard)
				originalInfo, _ := os.Stat(dst)
				if info.Size() < originalInfo.Size() {
					os.Remove(dst)
					os.Rename(tempDst, dst)
				} else {
					// Compressed is larger or equal (rare but possible), just keep original
					os.Remove(tempDst)
				}
			} else {
				os.Remove(tempDst)
			}
		} else {
			fmt.Printf("⚠️ Compression failed for %s: %v\n", filename, err)
			os.Remove(tempDst)
		}
	}

	// Return the relative URL
	photoURL := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"url": photoURL,
	}, "File uploaded successfully"))
}
