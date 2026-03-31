package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/utils"
	"github.com/xuri/excelize/v2"
)

type ReportHandler struct {
	db *pgxpool.Pool
}

func NewReportHandler(db *pgxpool.Pool) *ReportHandler {
	return &ReportHandler{db: db}
}

func (h *ReportHandler) ExportItems(c *gin.Context) {
	ctx := context.Background()

	query := `
		SELECT i.code, i.name, c.name as category, i.location, i.condition, i.status, i.purchase_date, i.purchase_price
		FROM items i
		LEFT JOIN categories c ON i.category_id = c.id
		ORDER BY i.name ASC
	`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch data for report"))
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := "Daftar Barang"
	index, _ := f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	headers := []string{"No", "Kode Aset", "Nama Barang", "Kategori", "Lokasi", "Kondisi", "Status", "Tgl Pembelian", "Harga"}
	for i, head := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, head)
	}

	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"4F46E5"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetRowStyle(sheetName, 1, 1, style)

	rowIdx := 2
	for rows.Next() {
		var code, name, location, condition, status string
		var category *string
		var purchaseDate *time.Time
		var purchasePrice *float64

		err := rows.Scan(&code, &name, &category, &location, &condition, &status, &purchaseDate, &purchasePrice)
		if err != nil {
			continue
		}

		catName := "-"
		if category != nil {
			catName = *category
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx), rowIdx-1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx), code)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIdx), name)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIdx), catName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIdx), location)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIdx), condition)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowIdx), status)

		if purchaseDate != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowIdx), purchaseDate.Format("2006-01-02"))
		}
		if purchasePrice != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowIdx), *purchasePrice)
		}

		rowIdx++
	}

	f.SetActiveSheet(index)

	fileName := fmt.Sprintf("SIS-INV_Laporan-Barang_%s.xlsx", time.Now().Format("20060102_1504"))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to generate Excel file"))
	}
}

// ExportTransactions exports all transactions to Excel (F05.4)
func (h *ReportHandler) ExportTransactions(c *gin.Context) {
	ctx := context.Background()

	query := `
		SELECT 
			t.borrower_type, t.student_nis, t.student_name, t.student_class,
			t.status, t.borrowed_at, t.due_date, t.returned_at, t.purpose,
			t.return_condition,
			i.code as item_code, i.name as item_name,
			u.full_name as teacher_name
		FROM transactions t
		JOIN items i ON t.item_id = i.id
		JOIN users u ON t.borrowed_by = u.id
		ORDER BY t.borrowed_at DESC
	`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data"))
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Laporan Transaksi"
	idx, _ := f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	hdrs := []string{"No", "Kode Barang", "Nama Barang", "Tipe", "NIS", "Nama Peminjam", "Kelas", "Guru Proses", "Tgl Pinjam", "Jatuh Tempo", "Tgl Kembali", "Status", "Kondisi Kembali", "Tujuan"}
	for i, hd := range hdrs {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, hd)
	}
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"4F46E5"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetRowStyle(sheet, 1, 1, style)

	row := 2
	for rows.Next() {
		var borrowerType, status, itemCode, itemName, teacherName string
		var studentNIS, studentName, studentClass, purpose, returnCondition *string
		var borrowedAt time.Time
		var dueDate *time.Time
		var returnedAt *time.Time

		if err := rows.Scan(&borrowerType, &studentNIS, &studentName, &studentClass,
			&status, &borrowedAt, &dueDate, &returnedAt, &purpose, &returnCondition,
			&itemCode, &itemName, &teacherName); err != nil {
			continue
		}

		peminjam := teacherName
		if borrowerType == "STUDENT" && studentName != nil {
			peminjam = *studentName
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), row-1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), itemCode)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), itemName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), borrowerType)
		if studentNIS != nil {
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *studentNIS)
		}
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), peminjam)
		if studentClass != nil {
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *studentClass)
		}
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), teacherName)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), borrowedAt.Format("2006-01-02 15:04"))
		if dueDate != nil {
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), dueDate.Format("2006-01-02 15:04"))
		}
		if returnedAt != nil {
			f.SetCellValue(sheet, fmt.Sprintf("K%d", row), returnedAt.Format("2006-01-02 15:04"))
		}
		f.SetCellValue(sheet, fmt.Sprintf("L%d", row), status)
		if returnCondition != nil {
			f.SetCellValue(sheet, fmt.Sprintf("M%d", row), *returnCondition)
		}
		if purpose != nil {
			f.SetCellValue(sheet, fmt.Sprintf("N%d", row), *purpose)
		}
		row++
	}

	f.SetActiveSheet(idx)
	fileName := fmt.Sprintf("SIS-INV_Laporan-Transaksi_%s.xlsx", time.Now().Format("20060102_1504"))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	f.Write(c.Writer)
}

// ActiveBorrowings returns all currently active borrowings with optional class filter (F05.2, F05.3)
func (h *ReportHandler) ActiveBorrowings(c *gin.Context) {
	classFilter := c.Query("class")
	ctx := context.Background()

	query := `
		SELECT 
			t.id, t.borrower_type, t.student_nis, t.student_name, t.student_class,
			t.borrowed_at, t.due_date, t.purpose,
			i.code as item_code, i.name as item_name,
			u.full_name as teacher_name
		FROM transactions t
		JOIN items i ON t.item_id = i.id
		JOIN users u ON t.borrowed_by = u.id
		WHERE t.status = 'ACTIVE'
	`
	args := []interface{}{}
	if classFilter != "" {
		query += " AND t.student_class = $1"
		args = append(args, classFilter)
	}
	query += " ORDER BY t.borrowed_at DESC"

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data"))
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, borrowerType, itemCode, itemName, teacherName string
		var studentNIS, studentName, studentClass, purpose *string
		var borrowedAt time.Time
		var dueDate *time.Time

		if err := rows.Scan(&id, &borrowerType, &studentNIS, &studentName, &studentClass,
			&borrowedAt, &dueDate, &purpose, &itemCode, &itemName, &teacherName); err != nil {
			continue
		}

		overdue := false
		if dueDate != nil && time.Now().After(*dueDate) {
			overdue = true
		}

		result = append(result, map[string]interface{}{
			"id":            id,
			"borrower_type": borrowerType,
			"student_nis":   studentNIS,
			"student_name":  studentName,
			"student_class": studentClass,
			"borrowed_at":   borrowedAt,
			"due_date":      dueDate,
			"purpose":       purpose,
			"item_code":     itemCode,
			"item_name":     itemName,
			"teacher_name":  teacherName,
			"is_overdue":    overdue,
		})
	}

	if result == nil {
		result = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(result, ""))
}

// OverdueReport returns only overdue borrowings (F05.2)
func (h *ReportHandler) OverdueReport(c *gin.Context) {
	ctx := context.Background()

	query := `
		SELECT 
			t.id, t.borrower_type, t.student_nis, t.student_name, t.student_class,
			t.borrowed_at, t.due_date, t.purpose,
			i.code as item_code, i.name as item_name,
			u.full_name as teacher_name
		FROM transactions t
		JOIN items i ON t.item_id = i.id
		JOIN users u ON t.borrowed_by = u.id
		WHERE t.status = 'ACTIVE' AND t.due_date < NOW()
		ORDER BY t.due_date ASC
	`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data"))
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, borrowerType, itemCode, itemName, teacherName string
		var studentNIS, studentName, studentClass, purpose *string
		var borrowedAt time.Time
		var dueDate *time.Time

		if err := rows.Scan(&id, &borrowerType, &studentNIS, &studentName, &studentClass,
			&borrowedAt, &dueDate, &purpose, &itemCode, &itemName, &teacherName); err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"id":            id,
			"borrower_type": borrowerType,
			"student_nis":   studentNIS,
			"student_name":  studentName,
			"student_class": studentClass,
			"borrowed_at":   borrowedAt,
			"due_date":      dueDate,
			"purpose":       purpose,
			"item_code":     itemCode,
			"item_name":     itemName,
			"teacher_name":  teacherName,
		})
	}

	if result == nil {
		result = []map[string]interface{}{}
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(result, ""))
}

// TransactionHistory returns all transactions with a limit of 100 (F05.2)
func (h *ReportHandler) TransactionHistory(c *gin.Context) {
	ctx := context.Background()

	query := `
		SELECT 
			t.borrower_type, t.student_name, t.student_class,
			t.status, t.borrowed_at, t.returned_at,
			i.code as item_code, i.name as item_name,
			u.full_name as teacher_name
		FROM transactions t
		JOIN items i ON t.item_id = i.id
		JOIN users u ON t.borrowed_by = u.id
		ORDER BY t.borrowed_at DESC
		LIMIT 100
	`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil riwayat"))
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var borrowerType, status, itemCode, itemName, teacherName string
		var studentName, studentClass *string
		var borrowedAt time.Time
		var returnedAt *time.Time

		if err := rows.Scan(&borrowerType, &studentName, &studentClass,
			&status, &borrowedAt, &returnedAt,
			&itemCode, &itemName, &teacherName); err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"borrower_type": borrowerType,
			"student_name":  studentName,
			"student_class": studentClass,
			"status":        status,
			"borrowed_at":   borrowedAt,
			"returned_at":   returnedAt,
			"item_code":     itemCode,
			"item_name":     itemName,
			"teacher_name":  teacherName,
		})
	}

	if result == nil { result = []map[string]interface{}{} }
	c.JSON(http.StatusOK, utils.SuccessResponse(result, ""))
}

// AuditLogs returns the last 100 system audit logs (F06.3)
func (h *ReportHandler) AuditLogs(c *gin.Context) {
	ctx := context.Background()

	query := `
		SELECT 
			a.id, a.action, a.entity_type, a.entity_id, a.description, 
			a.ip_address, a.created_at,
			u.full_name as user_name
		FROM audit_logs a
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY a.created_at DESC
		LIMIT 100
	`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil log audit"))
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int
		var action, entityType, description, userName string
		var entityId *string
		var ipAddress *string
		var createdAt time.Time

		if err := rows.Scan(&id, &action, &entityType, &entityId, &description, &ipAddress, &createdAt, &userName); err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"id":          id,
			"action":      action,
			"entity_type": entityType,
			"entity_id":   entityId,
			"description": description,
			"ip_address":  ipAddress,
			"created_at":  createdAt,
			"user_name":   userName,
		})
	}

	if result == nil { result = []map[string]interface{}{} }
	c.JSON(http.StatusOK, utils.SuccessResponse(result, ""))
}
