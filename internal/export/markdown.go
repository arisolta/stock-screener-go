package export

import (
	"fmt"
	"os"
	"strings"

	"stock-screener/internal/metrics"
	"stock-screener/internal/ui"
)

func ExportMarkdown(results []*metrics.ScreenResult, path string) error {
	var headers []string
	for _, col := range ui.AllExportColumns {
		headers = append(headers, col.Header)
	}

	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	var rows [][]string
	for _, r := range results {
		var row []string
		for i, col := range ui.AllExportColumns {
			val := ui.GetFormattedValue(r, col.Key)
			row = append(row, val)
			if len(val) > widths[i] {
				widths[i] = len(val)
			}
		}
		rows = append(rows, row)
	}

	var sb strings.Builder
	// Header
	var headerCells []string
	var sepCells []string
	for i, h := range headers {
		headerCells = append(headerCells, fmt.Sprintf("%-*s", widths[i], h))
		sepCells = append(sepCells, strings.Repeat("-", widths[i]))
	}
	sb.WriteString("| " + strings.Join(headerCells, " | ") + " |\n")
	sb.WriteString("| " + strings.Join(sepCells, " | ") + " |\n")

	// Rows
	for _, row := range rows {
		var cells []string
		for i, val := range row {
			cells = append(cells, fmt.Sprintf("%-*s", widths[i], val))
		}
		sb.WriteString("| " + strings.Join(cells, " | ") + " |\n")
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}
