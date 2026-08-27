package export

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"stock-screener/internal/metrics"
	"stock-screener/internal/ui"
)

func ExportExcel(results []*metrics.ScreenResult, path string) error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Screener Results"
	defaultSheet := f.GetSheetName(0)
	f.SetSheetName(defaultSheet, sheetName)

	// Create header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#2A4B7C"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return err
	}

	// Write headers
	for i, col := range ui.AllExportColumns {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return err
		}
		f.SetCellValue(sheetName, cell, col.Header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Write rows
	for rowIdx, r := range results {
		for colIdx, col := range ui.AllExportColumns {
			cell, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			if err != nil {
				return err
			}
			val := ui.GetFormattedValue(r, col.Key)
			f.SetCellValue(sheetName, cell, val)
		}
	}

	// Auto-fit column widths
	for i, col := range ui.AllExportColumns {
		maxLen := len(col.Header)
		for _, r := range results {
			val := ui.GetFormattedValue(r, col.Key)
			if len(val) > maxLen {
				maxLen = len(val)
			}
		}
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheetName, colName, colName, float64(maxLen+3))
	}

	return f.SaveAs(path)
}

func ExportResults(results []*metrics.ScreenResult, exportPath string) error {
	ext := ""
	for i := len(exportPath) - 1; i >= 0; i-- {
		if exportPath[i] == '.' {
			ext = exportPath[i:]
			break
		}
	}

	switch ext {
	case ".csv":
		return ExportCSV(results, exportPath)
	case ".xlsx", ".xls":
		return ExportExcel(results, exportPath)
	case ".md", ".markdown":
		return ExportMarkdown(results, exportPath)
	default:
		return fmt.Errorf("unsupported export format: %s (use .csv, .xlsx, .xls, .md, or .markdown)", ext)
	}
}
