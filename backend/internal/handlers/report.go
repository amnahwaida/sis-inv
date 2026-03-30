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

	// Set header
	headers := []string{"No", "Kode Aset", "Nama Barang", "Kategori", "Lokasi", "Kondisi", "Status", "Tgl Pembelian", "Harga"}
	for i, head := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, head)
	}

	// Stylize Header
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

	// Set headers for file download
	fileName := fmt.Sprintf("SIS-INV_Laporan-Barang_%s.xlsx", time.Now().Format("20060102_1504"))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to generate Excel file"))
	}
}
