package export

import (
	"encoding/csv"
	"os"

	"stock-screener/internal/metrics"
	"stock-screener/internal/ui"
)

func ExportCSV(results []*metrics.ScreenResult, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Write header
	var headers []string
	for _, col := range ui.AllExportColumns {
		headers = append(headers, col.Header)
	}
	if err := w.Write(headers); err != nil {
		return err
	}

	// Write rows
	for _, r := range results {
		var row []string
		for _, col := range ui.AllExportColumns {
			row = append(row, ui.GetFormattedValue(r, col.Key))
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return nil
}
