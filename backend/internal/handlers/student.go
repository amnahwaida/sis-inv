package handlers

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"

	"strings"

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

	actorId, _ := c.Get("userID")
	utils.LogAudit(h.db, actorId.(string), "CREATE_STUDENT", "STUDENT", "00000000-0000-0000-0000-000000000000", "Added student NIS: "+req.NIS, c.ClientIP())

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
	if req.NIS != oldNis { changes = append(changes, fmt.Sprintf("NIS: %s -> %s", oldNis, req.NIS)) }
	if req.FullName != oldName { changes = append(changes, fmt.Sprintf("Nama: %s -> %s", oldName, req.FullName)) }
	if req.Class != oldClass { changes = append(changes, fmt.Sprintf("Kelas: %s -> %s", oldClass, req.Class)) }
	if req.IsActive != oldActive { changes = append(changes, fmt.Sprintf("Aktif: %v -> %v", oldActive, req.IsActive)) }

	_, err = h.db.Exec(ctx,
		`UPDATE students SET nis = $1, full_name = $2, class = $3, is_active = $4, updated_at = NOW() WHERE id = $5`,
		req.NIS, req.FullName, req.Class, req.IsActive, id)
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
	id := c.Param("id")
	// Get student info for audit
	var stName, stNis string
	_ = h.db.QueryRow(context.Background(), "SELECT full_name, nis FROM students WHERE id = $1", id).Scan(&stName, &stNis)

	_, err := h.db.Exec(context.Background(), "DELETE FROM students WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal menghapus data: Siswa mungkin sudah memiliki transaksi"))
		return
	}

	actorId, _ := c.Get("userID")
	// Since students.id is SERIAL (int) and audit_logs.entity_id is UUID,
	// we use a NULL or a fixed UUID for student entity_id in audit log to prevent cast errors,
	// or we pass the string ID only if our LogAudit handles non-UUID strings (currently it doesn't).
	// For now, we'll use a placeholder UUID: '00000000-0000-0000-0000-000000000000'
	utils.LogAudit(h.db, actorId.(string), "DELETE_STUDENT", "STUDENT", "00000000-0000-0000-0000-000000000000", fmt.Sprintf("Deleted student: %s (NIS: %s)", stName, stNis), c.ClientIP())

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

func (h *StudentHandler) ExportCSV(c *gin.Context) {
	rows, err := h.db.Query(context.Background(), "SELECT nis, full_name, class FROM students ORDER BY nis ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengekspor data"))
		return
	}
	defer rows.Close()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=daftar_siswa.csv")
	c.Header("Content-Type", "text/csv")

	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{"NIS", "Nama Lengkap", "Kelas"})

	for rows.Next() {
		var nis, name, class string
		if err := rows.Scan(&nis, &name, &class); err == nil {
			_ = w.Write([]string{nis, name, class})
		}
	}
	w.Flush()
}

func (h *StudentHandler) ImportCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "File tidak ditemukan"))
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal membuka file"))
		return
	}
	defer src.Close()

	reader := csv.NewReader(src)
	// Skip header
	_, _ = reader.Read()

	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Format CSV tidak valid"))
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
	for _, record := range records {
		if len(record) < 3 { continue }
		nis := record[0]
		name := record[1]
		class := record[2]

		if nis == "" || name == "" { continue }

		_, err = tx.Exec(ctx, `
			INSERT INTO students (nis, full_name, class)
			VALUES ($1, $2, $3)
			ON CONFLICT (nis) DO UPDATE SET 
				full_name = EXCLUDED.full_name,
				class = EXCLUDED.class,
				updated_at = NOW()
		`, nis, name, class)
		if err == nil {
			imported++
		}
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal menyimpan data ke database"))
		return
	}

	actorId, _ := c.Get("userID")
	utils.LogAudit(h.db, actorId.(string), "IMPORT_STUDENTS", "STUDENT", "00000000-0000-0000-0000-000000000000", fmt.Sprintf("Imported %d students via CSV", imported), c.ClientIP())

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"total": imported}, "Berhasil mengimpor data siswa"))
}
