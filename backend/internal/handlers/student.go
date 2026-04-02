package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/utils"
	"github.com/xuri/excelize/v2"
)

type StudentHandler struct {
	db *pgxpool.Pool
}

func NewStudentHandler(db *pgxpool.Pool) *StudentHandler {
	return &StudentHandler{db: db}
}

func (h *StudentHandler) GetUniqueClasses(c *gin.Context) {
	query := `SELECT DISTINCT class FROM students WHERE class IS NOT NULL AND class != '' ORDER BY class ASC`
	rows, err := h.db.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data kelas"))
		return
	}
	defer rows.Close()

	var classes []string
	for rows.Next() {
		var className string
		if err := rows.Scan(&className); err == nil {
			classes = append(classes, className)
		}
	}
	
	if classes == nil { classes = []string{} }
	c.JSON(http.StatusOK, utils.SuccessResponse(classes, ""))
}

func (h *StudentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	search := c.Query("search")
	class := c.Query("class")
	
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 15 }
	offset := (page - 1) * pageSize

	query := `SELECT id, nis, full_name, class, is_active, created_at FROM students WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM students WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if search != "" {
		searchStr := "%" + search + "%"
		query += fmt.Sprintf(" AND (nis ILIKE $%d OR full_name ILIKE $%d)", argIdx, argIdx)
		countQuery += fmt.Sprintf(" AND (nis ILIKE $%d OR full_name ILIKE $%d)", argIdx, argIdx)
		args = append(args, searchStr)
		argIdx++
	}

	if class != "" {
		query += fmt.Sprintf(" AND class = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND class = $%d", argIdx)
		args = append(args, class)
		argIdx++
	}

	// Get total count
	var total int
	err := h.db.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil { total = 0 }

	query += fmt.Sprintf(" ORDER BY full_name ASC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(context.Background(), query, args...)
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
				"nisn":       nis, // Alias for compatibility
				"full_name":  name,
				"class":      class,
				"is_active":  isActive,
				"created_at": createdAt,
			})
		}
	}
	
	if students == nil { students = []map[string]interface{}{} }
	
	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": students,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}, ""))
}

func (h *StudentHandler) Create(c *gin.Context) {
	var req struct {
		NIS      string `json:"nis"`
		NISN     string `json:"nisn"`
		FullName string `json:"full_name" binding:"required"`
		Class    string `json:"class"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Data tidak valid"))
		return
	}

	nis := req.NIS
	if nis == "" { nis = req.NISN }
	if nis == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "NIS atau NISN wajib diisi"))
		return
	}

	_, err := h.db.Exec(context.Background(), 
		`INSERT INTO students (nis, full_name, class) VALUES ($1, $2, $3)`,
		nis, req.FullName, req.Class)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal menyimpan data: NIS mungkin sudah terpakai"))
		return
	}

	actorId, _ := c.Get("userID")
	auditDesc := fmt.Sprintf("Mendaftarkan siswa baru: %s (NIS: %s, Kelas: %s). Data siswa kini tersedia untuk transaksi peminjaman.", 
		req.FullName, nis, req.Class)
	utils.LogAudit(h.db, actorId.(string), "CREATE_STUDENT", "STUDENT", "00000000-0000-0000-0000-000000000000", auditDesc, c.ClientIP())

	c.JSON(http.StatusCreated, utils.SuccessResponse(nil, "Siswa berhasil ditambahkan"))
}

func (h *StudentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		NIS      string `json:"nis"`
		NISN     string `json:"nisn"`
		FullName string `json:"full_name" binding:"required"`
		Class    string `json:"class"`
		IsActive bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Data tidak valid"))
		return
	}

	nis := req.NIS
	if nis == "" { nis = req.NISN }
	if nis == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "NIS atau NISN wajib diisi"))
		return
	}

	ctx := context.Background()

	// Fetch old data for audit
	var oldNis, oldName, oldClass string
	var oldActive bool
	err := h.db.QueryRow(ctx, "SELECT nis, full_name, COALESCE(class, ''), is_active FROM students WHERE id = $1", id).
		Scan(&oldNis, &oldName, &oldClass, &oldActive)
	
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Data siswa tidak ditemukan"))
		return
	}

	changes := []string{}
	if nis != oldNis { changes = append(changes, fmt.Sprintf("NIS: %s -> %s", oldNis, nis)) }
	if req.FullName != oldName { changes = append(changes, fmt.Sprintf("Nama: %s -> %s", oldName, req.FullName)) }
	if req.Class != oldClass { changes = append(changes, fmt.Sprintf("Kelas: %s -> %s", oldClass, req.Class)) }
	if req.IsActive != oldActive { changes = append(changes, fmt.Sprintf("Aktif: %v -> %v", oldActive, req.IsActive)) }

	_, err = h.db.Exec(ctx,
		`UPDATE students SET nis = $1, full_name = $2, class = $3, is_active = $4, updated_at = NOW() WHERE id = $5`,
		nis, req.FullName, req.Class, req.IsActive, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal memperbarui data"))
		return
	}

	actorId, _ := c.Get("userID")
	auditDesc := "Updated student profile."
	if len(changes) > 0 {
		auditDesc = "Updated student. Changes: " + strings.Join(changes, " | ")
	}
	utils.LogAudit(h.db, actorId.(string), "UPDATE_STUDENT", "STUDENT", id, auditDesc, c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Siswa berhasil diperbarui"))
}

func (h *StudentHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	// Convert ID string to int since students.id is SERIAL
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "ID siswa tidak valid"))
		return
	}

	// Get student info for audit
	var stName, stNis string
	err = h.db.QueryRow(context.Background(), "SELECT full_name, nis FROM students WHERE id = $1", id).Scan(&stName, &stNis)

	_, err = h.db.Exec(context.Background(), "DELETE FROM students WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal menghapus data: Siswa mungkin sudah memiliki transaksi"))
		return
	}

	actorId, _ := c.Get("userID")
	
	// Fetch student class for audit
	var stClass string
	h.db.QueryRow(context.Background(), "SELECT class FROM students WHERE id = $1", id).Scan(&stClass)

	auditDesc := fmt.Sprintf("Menghapus data siswa: %s (NIS: %s, Kelas: %s) dari sistem.", stName, stNis, stClass)
	utils.LogAudit(h.db, actorId.(string), "DELETE_STUDENT", "STUDENT", "00000000-0000-0000-0000-000000000000", auditDesc, c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(nil, "Siswa berhasil dihapus"))
}

func (h *StudentHandler) Search(c *gin.Context) {
	queryParam := c.Query("q")
	if queryParam == "" {
		c.JSON(http.StatusOK, utils.SuccessResponse([]interface{}{}, ""))
		return
	}

	searchTerm := "%" + queryParam + "%"
	includeInactive := c.Query("include_inactive") == "true"
	
	query := `
		SELECT nis, full_name, class 
		FROM students 
		WHERE (nis ILIKE $1 OR full_name ILIKE $1)`
	
	if !includeInactive {
		query += ` AND is_active = true`
	}
	query += ` ORDER BY full_name ASC LIMIT 10`

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
				"nisn":      nis, // Alias for compatibility
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

func (h *StudentHandler) ExportExcel(c *gin.Context) {
	rows, err := h.db.Query(context.Background(), "SELECT nis, full_name, class FROM students ORDER BY nis ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengekspor data"))
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Data Siswa"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	// Headers
	hdrs := []string{"NIS", "Nama Lengkap", "Kelas"}
	for i, hd := range hdrs {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, hd)
	}

	// Header Style
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4F46E5"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetRowStyle(sheet, 1, 1, style)

	row := 2
	for rows.Next() {
		var nis, name, class string
		if err := rows.Scan(&nis, &name, &class); err == nil {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), nis)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), name)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), class)
			row++
		}
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=SIS-INV_Daftar_Siswa.xlsx")
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal menulis file excel"))
	}
}

func (h *StudentHandler) ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "File tidak ditemukan"))
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal membuka file"))
		return
	}
	defer openedFile.Close()

	f, err := excelize.OpenReader(openedFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Format file Excel tidak valid"))
		return
	}
	defer f.Close()

	// Get all rows from the first sheet
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil || len(rows) < 2 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Data Excel kosong atau tidak terbaca"))
		return
	}

	ctx := context.Background()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal memulai transaksi"))
		return
	}
	defer tx.Rollback(ctx)

	imported := 0
	for i, row := range rows {
		// Minimum validation: Row must have at least 2 columns (NIS, FullName)
		if i == 0 || len(row) < 2 { continue } 

		nis := strings.TrimSpace(row[0])
		name := strings.TrimSpace(row[1])
		class := ""
		if len(row) >= 3 { class = strings.TrimSpace(row[2]) }

		// Basic data validation
		if nis == "" || name == "" || strings.EqualFold(nis, "NIS") { 
			continue 
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO students (nis, full_name, class, is_active)
			VALUES ($1, $2, $3, true)
			ON CONFLICT (nis) DO UPDATE SET 
				full_name = EXCLUDED.full_name,
				class = EXCLUDED.class,
				is_active = true,
				updated_at = NOW()
		`, nis, name, class)
		
		if err == nil {
			imported++
		}
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal menyimpan data"))
		return
	}

	actorId, _ := c.Get("userID")
	utils.LogAudit(h.db, actorId.(string), "IMPORT_STUDENTS", "STUDENT", "00000000-0000-0000-0000-000000000000", fmt.Sprintf("Imported %d students via Excel", imported), c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"total": imported}, "Berhasil mengimpor data siswa"))
}
