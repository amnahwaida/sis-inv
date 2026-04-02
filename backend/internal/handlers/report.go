package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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
		SELECT i.code, i.name, c.name as category, COALESCE(l.name, i.location, ''), i.condition, i.status, i.purchase_date, i.purchase_price
		FROM items i
		LEFT JOIN categories c ON i.category_id = c.id
		LEFT JOIN locations l ON i.location_id = l.id
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
	classFilter := c.Query("class")
	typeFilter := c.Query("type")

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
		WHERE 1=1
	`
	
	args := []interface{}{}
	argIndex := 1

	if classFilter != "" {
		query += fmt.Sprintf(" AND t.student_class = $%d", argIndex)
		args = append(args, classFilter)
		argIndex++
	}

	if typeFilter == "active" {
		query += " AND t.status = 'ACTIVE'"
	} else if typeFilter == "overdue" {
		query += " AND t.status = 'ACTIVE' AND t.due_date < NOW()"
	} else if typeFilter == "history" {
		// History usually implies all or returned, but we'll show all if history is requested
		// as per the frontend's current logic for history tab.
	}

	query += " ORDER BY t.borrowed_at DESC"

	rows, err := h.db.Query(ctx, query, args...)
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

// ImportItemsTemplate generates a template for bulk item import
func (h *ReportHandler) ImportItemsTemplate(c *gin.Context) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Template Import Barang"
	idx, _ := f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	// Headers and Sample Data
	hdrs := []string{"Kode Aset", "Nama Barang", "Kategori", "Lokasi", "Kondisi", "Tipe Peminjam", "Tanggal Pembelian", "Harga", "Catatan"}
	sample := []string{"LAB-001", "Laptop ASUS Vivobook", "Elektronik", "Lab Komputer", "GOOD", "STUDENT_ALLOWED", "2024-01-20", "7500000", "Bantuan BOS 2024"}

	for i, hd := range hdrs {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, hd)
		
		// Style header
		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"4F46E5"}, Pattern: 1},
		})
		f.SetCellStyle(sheet, cell, cell, style)
		
		// Fill sample
		sampleCell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheet, sampleCell, sample[i])
	}

	// Add instructions as comments or a separate sheet (optional, keeping it simple for now)
	f.SetActiveSheet(idx)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=SIS-INV_Template_Barang.xlsx")
	f.Write(c.Writer)
}

// ImportStudentsTemplate generates a template for bulk student import
func (h *ReportHandler) ImportStudentsTemplate(c *gin.Context) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Template Import Siswa"
	idx, _ := f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	hdrs := []string{"NIS", "Nama Lengkap", "Kelas"}
	sample := []string{"2024001", "Ahmad Yusuf", "X-MIPA-1"}

	for i, hd := range hdrs {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, hd)
		
		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"4F46E5"}, Pattern: 1},
		})
		f.SetCellStyle(sheet, cell, cell, style)
		
		sampleCell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheet, sampleCell, sample[i])
	}

	f.SetActiveSheet(idx)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=SIS-INV_Template_Siswa.xlsx")
	f.Write(c.Writer)
}

// ActiveBorrowings returns all currently active borrowings with optional class filter (F05.2, F05.3)
func (h *ReportHandler) ActiveBorrowings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	search := c.Query("search")
	classFilter := c.Query("class")

	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 15 }
	offset := (page - 1) * pageSize

	ctx := context.Background()

	baseQuery := `
		FROM transactions t
		JOIN items i ON t.item_id = i.id
		JOIN users u ON t.borrowed_by = u.id
		WHERE t.status = 'ACTIVE'
	`
	args := []interface{}{}
	argIdx := 1

	if classFilter != "" {
		baseQuery += fmt.Sprintf(" AND t.student_class = $%d", argIdx)
		args = append(args, classFilter)
		argIdx++
	}

	if search != "" {
		searchStr := "%" + search + "%"
		baseQuery += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.code ILIKE $%d OR t.student_name ILIKE $%d OR u.full_name ILIKE $%d)", argIdx, argIdx, argIdx, argIdx)
		args = append(args, searchStr)
		argIdx++
	}

	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	err := h.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil { total = 0 }

	query := `
		SELECT 
			t.id, t.borrower_type, t.student_nis, t.student_name, t.student_class,
			t.borrowed_at, t.due_date, t.purpose, t.borrow_photo_url,
			i.code as item_code, i.name as item_name,
			u.full_name as teacher_name
	` + baseQuery + `
		ORDER BY t.borrowed_at DESC
		LIMIT $` + fmt.Sprint(argIdx) + ` OFFSET $` + fmt.Sprint(argIdx+1)
	
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data"))
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, borrowerType, itemCode, itemName, teacherName string
		var studentNIS, studentName, studentClass, purpose, borrowPhoto *string
		var borrowedAt time.Time
		var dueDate *time.Time

		if err := rows.Scan(&id, &borrowerType, &studentNIS, &studentName, &studentClass,
			&borrowedAt, &dueDate, &purpose, &borrowPhoto, &itemCode, &itemName, &teacherName); err != nil {
			continue
		}

		overdue := false
		if dueDate != nil && time.Now().After(*dueDate) {
			overdue = true
		}

		result = append(result, map[string]interface{}{
			"id":               id,
			"borrower_type":    borrowerType,
			"student_nis":      studentNIS,
			"student_name":     studentName,
			"student_class":    studentClass,
			"borrowed_at":      borrowedAt,
			"due_date":         dueDate,
			"purpose":          purpose,
			"borrow_photo_url": borrowPhoto,
			"item_code":        itemCode,
			"item_name":        itemName,
			"teacher_name":     teacherName,
			"is_overdue":       overdue,
		})
	}

	if result == nil { result = []map[string]interface{}{} }

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": result,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}, ""))
}

// OverdueReport returns only overdue borrowings (F05.2)
func (h *ReportHandler) OverdueReport(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	search := c.Query("search")
	classFilter := c.Query("class")

	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 15 }
	offset := (page - 1) * pageSize

	ctx := context.Background()

	baseQuery := `
		FROM transactions t
		JOIN items i ON t.item_id = i.id
		JOIN users u ON t.borrowed_by = u.id
		WHERE t.status = 'ACTIVE' AND t.due_date < NOW()
	`
	args := []interface{}{}
	argIdx := 1

	if classFilter != "" {
		baseQuery += fmt.Sprintf(" AND t.student_class = $%d", argIdx)
		args = append(args, classFilter)
		argIdx++
	}

	if search != "" {
		searchStr := "%" + search + "%"
		baseQuery += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.code ILIKE $%d OR t.student_name ILIKE $%d OR u.full_name ILIKE $%d)", argIdx, argIdx, argIdx, argIdx)
		args = append(args, searchStr)
		argIdx++
	}

	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	err := h.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil { total = 0 }

	query := `
		SELECT 
			t.id, t.borrower_type, t.student_nis, t.student_name, t.student_class,
			t.borrowed_at, t.due_date, t.purpose, t.borrow_photo_url,
			i.code as item_code, i.name as item_name,
			u.full_name as teacher_name
	` + baseQuery + `
		ORDER BY t.due_date ASC
		LIMIT $` + fmt.Sprint(argIdx) + ` OFFSET $` + fmt.Sprint(argIdx+1)
	
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data"))
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, borrowerType, itemCode, itemName, teacherName string
		var studentNIS, studentName, studentClass, purpose, borrowPhoto *string
		var borrowedAt time.Time
		var dueDate *time.Time

		if err := rows.Scan(&id, &borrowerType, &studentNIS, &studentName, &studentClass,
			&borrowedAt, &dueDate, &purpose, &borrowPhoto, &itemCode, &itemName, &teacherName); err != nil {
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

	if result == nil { result = []map[string]interface{}{} }

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": result,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}, ""))
}

// TransactionHistory returns all transactions with a limit of 100 (F05.2)
func (h *ReportHandler) TransactionHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	classFilter := c.Query("class")

	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }
	offset := (page - 1) * pageSize

	ctx := context.Background()

	baseQuery := `
		FROM transactions t
		JOIN items i ON t.item_id = i.id
		JOIN users u ON t.borrowed_by = u.id
		WHERE 1=1
	`
	var args []interface{}
	argIdx := 1

	if classFilter != "" {
		baseQuery += fmt.Sprintf(" AND t.student_class = $%d", argIdx)
		args = append(args, classFilter)
		argIdx++
	}

	if search != "" {
		searchStr := "%" + search + "%"
		baseQuery += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.code ILIKE $%d OR t.student_name ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, searchStr)
		argIdx++
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	err := h.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil { total = 0 }

	// Main query
	query := `
		SELECT 
			t.borrower_type, COALESCE(t.student_name, ''), COALESCE(t.student_class, ''),
			t.status, t.borrowed_at, t.returned_at, t.borrow_photo_url, t.return_photo_url,
			i.code as item_code, i.name as item_name,
			u.full_name as teacher_name
	` + baseQuery + `
		ORDER BY t.borrowed_at DESC
		LIMIT $` + fmt.Sprint(argIdx) + ` OFFSET $` + fmt.Sprint(argIdx+1)
	
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil riwayat"))
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var borrowerType, status, itemCode, itemName, teacherName string
		var studentName, studentClass, borrowPhoto, returnPhoto *string
		var borrowedAt time.Time
		var returnedAt *time.Time

		if err := rows.Scan(&borrowerType, &studentName, &studentClass,
			&status, &borrowedAt, &returnedAt, &borrowPhoto, &returnPhoto,
			&itemCode, &itemName, &teacherName); err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"borrower_type":    borrowerType,
			"student_name":     studentName,
			"student_class":    studentClass,
			"status":           status,
			"borrowed_at":      borrowedAt,
			"returned_at":      returnedAt,
			"borrow_photo_url": borrowPhoto,
			"return_photo_url": returnPhoto,
			"item_code":        itemCode,
			"item_name":        itemName,
			"teacher_name":     teacherName,
		})
	}

	if result == nil { result = []map[string]interface{}{} }
	
	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": result,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}, ""))
}

// AuditLogs returns the last 100 system audit logs (F06.3)
func (h *ReportHandler) AuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")
	action := c.Query("action")

	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	offset := (page - 1) * pageSize

	ctx := context.Background()

	baseQuery := `
		FROM audit_logs a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE 1=1
	`
	var args []interface{}
	argIdx := 1

	if action != "" {
		baseQuery += fmt.Sprintf(" AND a.action = $%d", argIdx)
		args = append(args, action)
		argIdx++
	}

	if search != "" {
		searchStr := "%" + search + "%"
		baseQuery += fmt.Sprintf(" AND (a.description ILIKE $%d OR u.full_name ILIKE $%d OR a.entity_type ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, searchStr)
		argIdx++
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	err := h.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil { total = 0 }

	// Main query
	query := `
		SELECT 
			a.id, a.action, COALESCE(a.entity_type, ''), a.entity_id, COALESCE(a.description, ''), 
			a.ip_address, a.created_at,
			COALESCE(u.full_name, 'System') as user_name
	` + baseQuery + `
		ORDER BY a.created_at DESC
		LIMIT $` + fmt.Sprint(argIdx) + ` OFFSET $` + fmt.Sprint(argIdx+1)
	
	args = append(args, pageSize, offset)

	rows, err := h.db.Query(ctx, query, args...)
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
	
	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"items": result,
		"meta": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + pageSize - 1) / pageSize,
		},
	}, ""))
}

// ExportAuditLogs exports the audit logs to Excel
func (h *ReportHandler) ExportAuditLogs(c *gin.Context) {
	ctx := context.Background()

	query := `
		SELECT 
			a.id, a.action, a.entity_type, a.description, 
			a.ip_address, a.created_at,
			u.full_name as user_name
		FROM audit_logs a
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY a.created_at DESC
	`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil log audit untuk export"))
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet := "Log Audit"
	idx, _ := f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	hdrs := []string{"No", "Waktu", "Pelaku", "Aksi", "Entitas", "Deskripsi", "Alamat IP"}
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
		var id int
		var action, entityType, description string
		var userName *string
		var ipAddress *string
		var createdAt time.Time

		if err := rows.Scan(&id, &action, &entityType, &description, &ipAddress, &createdAt, &userName); err != nil {
			continue
		}

		pelaku := "System"
		if userName != nil {
			pelaku = *userName
		}
		
		ip := "127.0.0.1"
		if ipAddress != nil {
			ip = *ipAddress
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), row-1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), createdAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), pelaku)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), action)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), entityType)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), description)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), ip)
		row++
	}

	f.SetActiveSheet(idx)
	fileName := fmt.Sprintf("SIS-INV_Log-Audit_%s.xlsx", time.Now().Format("20060102_1504"))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	f.Write(c.Writer)
}

// ExportMaintenanceLogs exports the maintenance logs to Excel
func (h *ReportHandler) ExportMaintenanceLogs(c *gin.Context) {
	ctx := context.Background()
	statusFilter := c.Query("status")

	query := `
		SELECT 
			m.id, i.code, i.name, u.full_name as reported_by, 
			m.issue_description, m.vendor, m.cost, m.status, 
			m.reported_at, m.completed_at, m.notes
		FROM maintenance_logs m
		JOIN items i ON m.item_id = i.id
		LEFT JOIN users u ON m.reported_by = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	if statusFilter != "" {
		query += " AND m.status = $1"
		args = append(args, statusFilter)
	}
	query += " ORDER BY m.reported_at DESC"

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Gagal mengambil data perbaikan untuk export"))
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet := "Log Perbaikan"
	idx, _ := f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	hdrs := []string{"No", "Kode Barang", "Nama Barang", "Pelapor", "Deskripsi Keluhan", "Vendor/Teknisi", "Biaya", "Status", "Tgl Lapor", "Tgl Selesai", "Catatan"}
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
		var id int
		var itemCode, itemName, issueDescription, status string
		var reportedBy, vendor, notes *string
		var cost *float64
		var reportedAt time.Time
		var completedAt *time.Time

		if err := rows.Scan(&id, &itemCode, &itemName, &reportedBy, &issueDescription, &vendor, &cost, &status, &reportedAt, &completedAt, &notes); err != nil {
			continue
		}

		pelapor := "-"
		if reportedBy != nil { pelapor = *reportedBy }
		vnd := "-"
		if vendor != nil && *vendor != "" { vnd = *vendor }
		cst := 0.0
		if cost != nil { cst = *cost }
		nts := "-"
		if notes != nil && *notes != "" { nts = *notes }

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), row-1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), itemCode)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), itemName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), pelapor)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), issueDescription)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), vnd)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), cst)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), status)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), reportedAt.Format("2006-01-02 15:04:05"))
		if completedAt != nil {
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), completedAt.Format("2006-01-02 15:04:05"))
		} else {
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), "-")
		}
		f.SetCellValue(sheet, fmt.Sprintf("K%d", row), nts)
		
		row++
	}

	f.SetActiveSheet(idx)
	fileName := fmt.Sprintf("SIS-INV_Log-Perbaikan_%s.xlsx", time.Now().Format("20060102_1504"))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	f.Write(c.Writer)
}
func (h *ReportHandler) ExportAuditSession(c *gin.Context) {
	sessionId := c.Param("id")
	ctx := context.Background()

	// 1. Get Session Info
	var s struct {
		ID           string
		LocationName string
		UserName     string
		Status       string
		StartedAt    time.Time
		FinishedAt   *time.Time
		Notes        *string
	}
	err := h.db.QueryRow(ctx,
		`SELECT s.id, l.name as location_name, u.full_name as user_name,
		        s.status, s.started_at, s.finished_at, s.notes
		 FROM audit_sessions s
		 JOIN locations l ON s.location_id = l.id
		 JOIN users u ON s.user_id = u.id
		 WHERE s.id = $1`, sessionId).Scan(
		&s.ID, &s.LocationName, &s.UserName, &s.Status, &s.StartedAt, &s.FinishedAt, &s.Notes,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Audit session not found"))
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Hasil Audit"
	idx, _ := f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	// Header Info
	f.SetCellValue(sheet, "A1", "LAPORAN HASIL STOCK OPNAME")
	f.SetCellValue(sheet, "A2", "Lokasi:")
	f.SetCellValue(sheet, "B2", s.LocationName)
	f.SetCellValue(sheet, "A3", "Petugas:")
	f.SetCellValue(sheet, "B3", s.UserName)
	f.SetCellValue(sheet, "A4", "Waktu Mulai:")
	f.SetCellValue(sheet, "B4", s.StartedAt.Format("2006-01-02 15:04"))
	f.SetCellValue(sheet, "A5", "Status:")
	f.SetCellValue(sheet, "B5", s.Status)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
	})
	f.SetCellStyle(sheet, "A1", "A1", headerStyle)

	// Items Table Header
	f.SetCellValue(sheet, "A7", "LIST BARANG DITEMUKAN")
	tableHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"4F46E5"}, Pattern: 1},
	})
	
	cols := []string{"No", "Kode Barang", "Nama Barang", "Kondisi Saat Scan", "Waktu Scan", "Catatan"}
	for i, col := range cols {
		cell, _ := excelize.CoordinatesToCellName(i+1, 8)
		f.SetCellValue(sheet, cell, col)
	}
	f.SetRowStyle(sheet, 8, 8, tableHeaderStyle)

	// Fetch Found Items
	rows, _ := h.db.Query(ctx, 
		`SELECT i.code, i.name, ai.found_condition, ai.scanned_at, COALESCE(ai.notes, '-')
		 FROM audit_items ai
		 JOIN items i ON ai.item_id = i.id
		 WHERE ai.session_id = $1 ORDER BY ai.scanned_at ASC`, sessionId)
	
	rowIdx := 9
	for rows.Next() {
		var code, name, cond, notes string
		var scannedAt time.Time
		rows.Scan(&code, &name, &cond, &scannedAt, &notes)
		
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIdx), rowIdx-8)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIdx), code)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIdx), name)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIdx), cond)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIdx), scannedAt.Format("15:04:05"))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowIdx), notes)
		rowIdx++
	}
	rows.Close()

	// Missing Items Section
	rowIdx += 2
	f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIdx), "LIST BARANG TIDAK DITEMUKAN (MISSING)")
	rowIdx++
	
	missingCols := []string{"No", "Kode Barang", "Nama Barang", "Status Sistem", "Kondisi Terakhir"}
	for i, col := range missingCols {
		cell, _ := excelize.CoordinatesToCellName(i+1, rowIdx)
		f.SetCellValue(sheet, cell, col)
	}
	f.SetRowStyle(sheet, rowIdx, rowIdx, tableHeaderStyle)
	
	startRow := rowIdx + 1
	rows, _ = h.db.Query(ctx,
		`SELECT i.code, i.name, i.status, i.condition
		 FROM items i
		 WHERE i.location_id = (SELECT location_id FROM audit_sessions WHERE id = $1)
		 AND i.id NOT IN (SELECT item_id FROM audit_items WHERE session_id = $1)`,
		sessionId)
	
	rowIdx++
	missingCount := 0
	for rows.Next() {
		var code, name, status, cond string
		rows.Scan(&code, &name, &status, &cond)
		
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIdx), rowIdx-startRow+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIdx), code)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIdx), name)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIdx), status)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIdx), cond)
		rowIdx++
		missingCount++
	}
	rows.Close()

	if missingCount == 0 {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", startRow), "Semua Barang Ditemukan")
	}

	f.SetActiveSheet(idx)
	fileName := fmt.Sprintf("Audit_%s_%s.xlsx", s.LocationName, time.Now().Format("20060102"))
	
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	f.Write(c.Writer)
}
