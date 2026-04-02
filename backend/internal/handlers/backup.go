package handlers

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/database"
	"github.com/vannyezha/sis-inv/internal/utils"
)

type BackupHandler struct {
	db *pgxpool.Pool
}

func NewBackupHandler(db *pgxpool.Pool) *BackupHandler {
	return &BackupHandler{db: db}
}

func (h *BackupHandler) Backup(c *gin.Context) {
	timestamp := time.Now().Format("20060102_150405")
	zipFileName := fmt.Sprintf("sis_inv_backup_%s.zip", timestamp)

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFileName))

	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	// 1. Export Database Tables
	tables := []string{"categories", "locations", "items", "students", "users", "transactions", "maintenance_logs", "audit_logs", "audit_sessions", "audit_items", "settings"}
	
	for _, table := range tables {
		data, err := h.exportTable(table)
		if err != nil {
			continue // Skip failed tables but log or handle
		}
		
		f, err := zipWriter.Create(fmt.Sprintf("data/%s.json", table))
		if err != nil {
			continue
		}
		f.Write(data)
	}

	// 2. Export Uploads Folder
	_ = filepath.Walk("uploads", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		
		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		f, err := zipWriter.Create(path)
		if err != nil {
			return nil
		}
		
		_, _ = io.Copy(f, file)
		return nil
	})
}

func (h *BackupHandler) Restore(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "No file uploaded"))
		return
	}

	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("restore_%d", time.Now().Unix()))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create temp directory"))
		return
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, "backup.zip")
	if err := c.SaveUploadedFile(file, zipPath); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to save upload"))
		return
	}

	// Unzip
	if err := h.unzip(zipPath, tempDir); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to unzip backup"))
		return
	}

	// Start Transaction for Restore
	ctx := context.Background()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to start transaction"))
		return
	}
	defer tx.Rollback(ctx)

	// Truncate tables (Order matters for FK)
	tables := []string{"audit_items", "audit_sessions", "maintenance_logs", "transactions", "items", "students", "users", "categories", "locations", "audit_logs", "settings"}
	for _, table := range tables {
		if _, err := tx.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)); err != nil {
			log.Printf("❌ Failed to truncate table %s: %v", table, err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to clear table %s: %v", table, err)))
			return
		}
	}

	// Import Data in strict order to respect Foreign Keys
	dataDir := filepath.Join(tempDir, "data")
	importOrder := []string{"categories", "locations", "users", "students", "items", "transactions", "maintenance_logs", "audit_sessions", "audit_items", "audit_logs", "settings"}
	
	for _, tableName := range importOrder {
		filePath := filepath.Join(dataDir, tableName+".json")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}

		content, _ := os.ReadFile(filePath)
		var rows []map[string]interface{}
		if err := json.Unmarshal(content, &rows); err == nil {
			log.Printf("📥 Importing %d rows into %s...", len(rows), tableName)
			for _, row := range rows {
				if err := h.importRow(ctx, tx, tableName, row); err != nil {
					log.Printf("❌ Failed to import row into %s: %v", tableName, err)
					c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to import %s: %v", tableName, err)))
					return
				}
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("❌ Failed to commit restore: %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to commit restore: "+err.Error()))
		return
	}

	// Restore Uploads
	_ = os.RemoveAll("uploads")
	_ = os.MkdirAll("uploads", 0755)
	h.copyDir(filepath.Join(tempDir, "uploads"), "uploads")

	// Ensure Admin exists
	_ = database.SeedDefaultAdmin(h.db)

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "System restored successfully"))
}

func (h *BackupHandler) Reset(c *gin.Context) {
	ctx := context.Background()
	
	// Truncate all tables
	tables := []string{"audit_items", "audit_sessions", "maintenance_logs", "transactions", "items", "students", "users", "categories", "locations", "audit_logs", "settings"}
	for _, table := range tables {
		_, err := h.db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to truncate table: "+table))
			return
		}
	}

	// Clear uploads
	_ = os.RemoveAll("uploads")
	_ = os.MkdirAll("uploads", 0755)

	// Re-seed
	if err := database.SeedDefaultAdmin(h.db); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to re-seed system"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "System reset successfully"))
}

// Helpers

func (h *BackupHandler) exportTable(table string) ([]byte, error) {
	rows, err := h.db.Query(context.Background(), fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	cols := rows.FieldDescriptions()
	
	for rows.Next() {
		values, _ := rows.Values()
		row := make(map[string]interface{})
		for i, col := range cols {
			row[string(col.Name)] = values[i]
		}
		result = append(result, row)
	}
	return json.Marshal(result)
}

func (h *BackupHandler) importRow(ctx context.Context, tx pgx.Tx, table string, row map[string]interface{}) error {
	// Generic insert logic
	keys := make([]string, 0, len(row))
	values := make([]interface{}, 0, len(row))
	placeholders := make([]string, 0, len(row))
	
	i := 1
	for k, v := range row {
		keys = append(keys, k)
		
		// Fix UUID encoding: if it's a slice of interface{} (from JSON), 
		// and it has 16 elements, it's likely a UUID that needs to be []byte.
		if slice, ok := v.([]interface{}); ok && len(slice) == 16 {
			bytes := make([]byte, 16)
			allNumeric := true
			for idx, val := range slice {
				if n, ok := val.(float64); ok {
					bytes[idx] = byte(n)
				} else {
					allNumeric = false
					break
				}
			}
			if allNumeric {
				v = bytes
			}
		}

		// Handle JSONB: if it's a map or slice (not 16-byte UUID), 
		// convert it back to a JSON string for pgx to handle as JSONB.
		if _, ok := v.(map[string]interface{}); ok {
			if jsonBytes, err := json.Marshal(v); err == nil {
				v = string(jsonBytes)
			}
		} else if slice, ok := v.([]interface{}); ok && len(slice) != 16 {
			if jsonBytes, err := json.Marshal(v); err == nil {
				v = string(jsonBytes)
			}
		}
		
		values = append(values, v)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		i++
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO NOTHING", 
		table, strings.Join(keys, ","), strings.Join(placeholders, ","))
		
	_, err := tx.Exec(ctx, query, values...)
	return err
}

func (h *BackupHandler) unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *BackupHandler) copyDir(src, dest string) {
	_ = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		relPath, _ := filepath.Rel(src, path)
		destPath := filepath.Join(dest, relPath)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		
		in, _ := os.Open(path)
		defer in.Close()
		out, _ := os.Create(destPath)
		defer out.Close()
		_, _ = io.Copy(out, in)
		return nil
	})
}
