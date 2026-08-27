package ui

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/fatih/color"

	"stock-screener/internal/metrics"
)

// SummaryColumns is a well-proportioned set of the most critical decision metrics that fits standard screens without wrapping.
var SummaryColumns = []ColumnDef{
	{Key: "Ticker", Label: "Ticker", Numeric: false},
	{Key: "CompanyName", Label: "Company", Numeric: false},
	{Key: "Score", Label: "Score", Numeric: true},
	{Key: "IntensityLevel", Label: "Intensity", Numeric: false},
	{Key: "AvgDepFactor3Y", Label: "EBITDA/EBIT 3Y Avg", Numeric: true},
	{Key: "ROICLatest", Label: "ROIC FY", Numeric: true},
	{Key: "ROICAvg3Y", Label: "ROIC 3Y Avg", Numeric: true},
	{Key: "OperatingMargin", Label: "EBIT Mrgn FY", Numeric: true},
	{Key: "RevCAGR3Y", Label: "Rev 3Y CAGR", Numeric: true},
	{Key: "AdjFCFNetIncomeRatio", Label: "AdjFCF/NI 3Y Avg", Numeric: true},
	{Key: "AdjFCF3YEVYield", Label: "3Y Sum FCF/EV", Numeric: true},
	{Key: "EVEBIT", Label: "EV/EBIT", Numeric: true},
	{Key: "NetDebtEBITDA", Label: "Debt/EBITDA", Numeric: true},
}

var (
	ansiRegex  = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	cyanBold   = color.New(color.FgCyan, color.Bold).SprintFunc()
	greenBold  = color.New(color.FgGreen, color.Bold).SprintFunc()
	yellowBold = color.New(color.FgYellow, color.Bold).SprintFunc()
	redBold    = color.New(color.FgRed, color.Bold).SprintFunc()
	whiteBold  = color.New(color.FgWhite, color.Bold).SprintFunc()
)

func visibleLen(s string) int {
	clean := ansiRegex.ReplaceAllString(s, "")
	return len([]rune(clean))
}

func padRight(s string, width int) string {
	vl := visibleLen(s)
	if vl >= width {
		return s
	}
	return s + strings.Repeat(" ", width-vl)
}

func padLeft(s string, width int) string {
	vl := visibleLen(s)
	if vl >= width {
		return s
	}
	return strings.Repeat(" ", width-vl) + s
}

func cardRow(left, right string) string {
	return fmt.Sprintf("│ %s │ %s │", padRight(left, 42), padRight(right, 41))
}

func item(label, val string) string {
	vlLabel := visibleLen(label)
	vlVal := visibleLen(val)
	// total target visible width inside item is 40 chars ("  • " is 4 chars + 36 chars content)
	remWidth := 36 - vlVal - 1
	if remWidth < vlLabel {
		remWidth = vlLabel
	}
	paddedLabel := padRight(label, remWidth)
	return fmt.Sprintf("  • %s %s", paddedLabel, val)
}

func miniBar(score float64, plain bool) string {
	if math.IsNaN(score) {
		return "[░░░░░░░░]"
	}
	filled := int(math.Round(score / 12.5))
	if filled > 8 {
		filled = 8
	}
	if filled < 0 {
		filled = 0
	}
	barStr := strings.Repeat("█", filled) + strings.Repeat("░", 8-filled)
	if plain {
		return "[" + barStr + "]"
	}
	if score >= 75 {
		return "[" + color.New(color.FgGreen).Sprint(barStr) + "]"
	} else if score >= 55 {
		return "[" + color.New(color.FgYellow).Sprint(barStr) + "]"
	}
	return "[" + color.New(color.FgRed).Sprint(barStr) + "]"
}

func formatMarketCap(mc float64) string {
	if math.IsNaN(mc) || mc <= 0 {
		return "N/A"
	}
	if mc >= 1e12 {
		return fmt.Sprintf("$%.2fT", mc/1e12)
	}
	if mc >= 1e9 {
		return fmt.Sprintf("$%.2fB", mc/1e9)
	}
	if mc >= 1e6 {
		return fmt.Sprintf("$%.2fM", mc/1e6)
	}
	return fmt.Sprintf("$%.0f", mc)
}

func formatMultipleX(v float64) string {
	if math.IsNaN(v) || v <= 0 {
		return "-"
	}
	return fmt.Sprintf("%.2fx", v)
}

func formatTrend(v float64, plain bool) string {
	if math.IsNaN(v) {
		return ""
	}
	pct := fmt.Sprintf("%+.1f%%", v*100)
	if plain {
		return pct
	}
	if v > 0.005 {
		return color.New(color.FgGreen).Sprintf("▲ %s", pct)
	} else if v < -0.005 {
		return color.New(color.FgRed).Sprintf("▼ %s", pct)
	}
	return pct
}

func generateInsights(r *metrics.ScreenResult, plain bool) (string, string) {
	var strengths []string
	var cautions []string

	// Strengths rules
	if r.ROICLatest >= 0.25 {
		strengths = append(strengths, fmt.Sprintf("Elite ROIC (%s)", FormatPercent(r.ROICLatest)))
	} else if r.ROICAvg3Y >= 0.15 {
		strengths = append(strengths, fmt.Sprintf("Strong 3Y Avg ROIC (%s)", FormatPercent(r.ROICAvg3Y)))
	}

	if r.AdjFCFNetIncomeRatio >= 0.85 {
		strengths = append(strengths, fmt.Sprintf("High Cash Conversion (%s 3Y Avg)", FormatPercent(r.AdjFCFNetIncomeRatio)))
	}

	if r.AvgDepFactor3Y < 1.25 && !math.IsNaN(r.AvgDepFactor3Y) {
		strengths = append(strengths, fmt.Sprintf("Asset-Light Capital (%.2f 3Y Avg)", r.AvgDepFactor3Y))
	}

	if r.ShareCountCAGR3Y <= -0.01 {
		strengths = append(strengths, fmt.Sprintf("Active Buybacks (%s 3Y CAGR)", FormatPercent(r.ShareCountCAGR3Y)))
	}

	if r.NetDebtEBITDA <= 0 && !math.IsNaN(r.NetDebtEBITDA) {
		strengths = append(strengths, "Net Cash Balance Sheet")
	} else if r.NetDebtEBITDA < 1.0 && !math.IsNaN(r.NetDebtEBITDA) {
		strengths = append(strengths, fmt.Sprintf("Low Leverage (%.2fx Debt/EBITDA)", r.NetDebtEBITDA))
	}

	if r.EBITCAGR3Y >= 0.15 {
		strengths = append(strengths, fmt.Sprintf("Fast EBIT Growth (+%.1f%% 3Y CAGR)", r.EBITCAGR3Y*100))
	} else if r.RevCAGR3Y >= 0.12 {
		strengths = append(strengths, fmt.Sprintf("Strong Revenue (+%.1f%% 3Y CAGR)", r.RevCAGR3Y*100))
	}

	if r.AdjFCF3YEVYield >= 0.07 {
		strengths = append(strengths, fmt.Sprintf("Attractive Owner Yield (%s 3Y Sum/EV)", FormatPercent(r.AdjFCF3YEVYield)))
	}

	// Cautions rules
	if r.AvgDepFactor3Y > 1.7 && !math.IsNaN(r.AvgDepFactor3Y) {
		cautions = append(cautions, fmt.Sprintf("High Reinvestment Burden (%.2f 3Y Avg Dep)", r.AvgDepFactor3Y))
	}

	if r.EVEBIT >= 28.0 && !math.IsNaN(r.EVEBIT) {
		cautions = append(cautions, fmt.Sprintf("Premium Multiple (%.1fx EV/EBIT)", r.EVEBIT))
	}

	if r.NetDebtEBITDA >= 3.0 && !math.IsNaN(r.NetDebtEBITDA) {
		cautions = append(cautions, fmt.Sprintf("Elevated Debt (%.2fx EBITDA)", r.NetDebtEBITDA))
	}

	if r.AdjFCFNetIncomeRatio < 0.60 && !math.IsNaN(r.AdjFCFNetIncomeRatio) {
		cautions = append(cautions, fmt.Sprintf("Weak Cash Conversion (%s 3Y Avg FCF/NI)", FormatPercent(r.AdjFCFNetIncomeRatio)))
	}

	if r.RevCAGR3Y <= 0.02 && !math.IsNaN(r.RevCAGR3Y) {
		cautions = append(cautions, fmt.Sprintf("Sluggish Top-Line Growth (%s 3Y CAGR)", FormatPercent(r.RevCAGR3Y)))
	}

	if r.ROICLatest < 0.10 && !math.IsNaN(r.ROICLatest) {
		cautions = append(cautions, fmt.Sprintf("Low Return on Capital (%s ROIC)", FormatPercent(r.ROICLatest)))
	}

	if r.ShareCountCAGR3Y >= 0.02 {
		cautions = append(cautions, fmt.Sprintf("Share Dilution (+%s 3Y CAGR)", FormatPercent(r.ShareCountCAGR3Y)))
	}

	if r.SBCRevRatio >= 0.06 {
		cautions = append(cautions, fmt.Sprintf("High SBC Burden (%s 3Y Avg of Rev)", FormatPercent(r.SBCRevRatio)))
	}

	if len(strengths) == 0 {
		strengths = append(strengths, "Stable business profile")
	}
	if len(cautions) == 0 {
		cautions = append(cautions, "No major red flags detected across key ratios")
	}

	maxStr := 3
	if len(strengths) > maxStr {
		strengths = strengths[:maxStr]
	}
	maxCaut := 2
	if len(cautions) > maxCaut {
		cautions = cautions[:maxCaut]
	}

	greenTag := "★ Strengths: "
	redTag := "⚠ Watch-Outs: "
	if !plain {
		greenTag = color.New(color.FgGreen, color.Bold).Sprint("★ Strengths: ")
		redTag = color.New(color.FgRed, color.Bold).Sprint("⚠ Watch-Outs: ")
	}

	return greenTag + strings.Join(strengths, " • "), redTag + strings.Join(cautions, " • ")
}

func calculateWidths(cols []ColumnDef, results []*metrics.ScreenResult) []int {
	widths := make([]int, len(cols))
	for i, c := range cols {
		w := len(c.Label)
		for _, r := range results {
			val := GetFormattedValue(r, c.Key)
			if len(val) > w {
				w = len(val)
			}
		}
		widths[i] = w
	}
	return widths
}

func renderSingleTable(cols []ColumnDef, results []*metrics.ScreenResult, plain bool) {
	if len(results) == 0 {
		return
	}

	widths := calculateWidths(cols, results)

	border := func(left, join, right string) string {
		var parts []string
		for _, w := range widths {
			parts = append(parts, strings.Repeat("─", w+2))
		}
		return left + strings.Join(parts, join) + right
	}

	fmt.Println(border("┌", "┬", "┐"))

	// Header row
	var headerCells []string
	for i, c := range cols {
		text := c.Label
		if c.Numeric {
			headerCells = append(headerCells, fmt.Sprintf(" %*s ", widths[i], text))
		} else {
			headerCells = append(headerCells, fmt.Sprintf(" %-*s ", widths[i], text))
		}
	}
	fmt.Println("│" + strings.Join(headerCells, "│") + "│")
	fmt.Println(border("├", "┼", "┤"))

	for _, r := range results {
		var cells []string
		for i, c := range cols {
			val := GetFormattedValue(r, c.Key)
			if c.Numeric {
				cells = append(cells, fmt.Sprintf(" %*s ", widths[i], val))
			} else {
				cells = append(cells, fmt.Sprintf(" %-*s ", widths[i], val))
			}
		}
		rowStr := "│" + strings.Join(cells, "│") + "│"
		if !plain {
			switch r.IntensityLevel {
			case "Low":
				rowStr = color.New(color.FgGreen).Sprint(rowStr)
			case "Moderate":
				rowStr = color.New(color.FgYellow).Sprint(rowStr)
			case "High":
				rowStr = color.New(color.FgRed).Sprint(rowStr)
			}
		}
		fmt.Println(rowStr)
	}

	fmt.Println(border("└", "┴", "┘"))
}

// PrintCardView displays all clean metrics for each company organized into 4 1-to-1 matching quadrants.
func PrintCardView(results []*metrics.ScreenResult, plain bool) {
	for idx, r := range results {
		if idx > 0 {
			fmt.Println()
		}

		scoreColor := greenBold
		qualityLabel := "Strong Quality"
		if r.Score < 60 {
			scoreColor = redBold
			qualityLabel = "Mixed / Weak"
		} else if r.Score < 75 {
			scoreColor = yellowBold
			qualityLabel = "Moderate Quality"
		}

		intensityColor := greenBold
		if r.IntensityLevel == "Moderate" {
			intensityColor = yellowBold
		} else if r.IntensityLevel == "High" {
			intensityColor = redBold
		}

		compName := r.CompanyName
		if len([]rune(compName)) > 26 {
			compName = string([]rune(compName)[:25]) + "…"
		}

		secInfo := r.Sector
		if secInfo == "" && r.Industry != "" {
			secInfo = r.Industry
		}
		if secInfo != "" {
			secInfo = " • " + secInfo
		}

		mktCapStr := formatMarketCap(r.MarketCap)

		fmt.Println("┌" + strings.Repeat("─", 88) + "┐")
		
		// Line 1: Header Identity
		line1 := fmt.Sprintf(" %s (%s)%s • MktCap: %s • Cur: %s",
			cyanBold(r.Ticker), whiteBold(compName), secInfo, mktCapStr, r.Currency)
		if plain {
			line1 = fmt.Sprintf(" %s (%s)%s • MktCap: %s • Cur: %s",
				r.Ticker, compName, secInfo, mktCapStr, r.Currency)
		}
		fmt.Printf("│ %s │\n", padRight(line1, 86))

		// Line 2: Composite Score and Intensity Badge
		line2 := fmt.Sprintf(" Overall Score: %s/100 %s %s • Intensity: %s (3Y Avg: %s)",
			scoreColor(fmt.Sprintf("%.1f", r.Score)), miniBar(r.Score, plain), qualityLabel,
			intensityColor(r.IntensityLevel), FormatMultiple(r.AvgDepFactor3Y))
		if plain {
			line2 = fmt.Sprintf(" Overall Score: %.1f/100 %s %s • Intensity: %s (3Y Avg: %s)",
				r.Score, miniBar(r.Score, plain), qualityLabel, r.IntensityLevel, FormatMultiple(r.AvgDepFactor3Y))
		}
		fmt.Printf("│ %s │\n", padRight(line2, 86))

		// 4 Sub-Scores mapping 1:1 to the 4 quadrants below
		fmt.Println("├" + strings.Repeat("─", 88) + "┤")
		subLine1 := fmt.Sprintf(" SUB-SCORES:  Returns (30%%)   %s %-3.0f     Growth (25%%)    %s %-3.0f",
			miniBar(r.ReturnsScore, plain), r.ReturnsScore,
			miniBar(r.GrowthScore, plain), r.GrowthScore)
		subLine2 := fmt.Sprintf("              Cash (20%%)      %s %-3.0f     Valuation (25%%) %s %-3.0f",
			miniBar(r.CashScore, plain), r.CashScore,
			miniBar(r.ValuationScore, plain), r.ValuationScore)
		fmt.Printf("│ %s │\n", padRight(subLine1, 86))
		fmt.Printf("│ %s │\n", padRight(subLine2, 86))

		fmt.Println("├" + strings.Repeat("─", 44) + "┬" + strings.Repeat("─", 43) + "┤")

		// QUADRANT 1 & 2
		sec1 := yellowBold("1. PROFITABILITY & RETURNS (30%)")
		sec2 := yellowBold("2. GROWTH & STABILITY (25%)")
		if plain {
			sec1 = "1. PROFITABILITY & RETURNS (30%)"
			sec2 = "2. GROWTH & STABILITY (25%)"
		}
		fmt.Println(cardRow(sec1, sec2))
		fmt.Println(cardRow(item("EBITDA/EBIT (Dep) FY:", FormatMultiple(r.DepFactorLatest)), item("Gross Margin FY:", FormatPercent(r.GrossMargin))))
		fmt.Println(cardRow(item("EBITDA/EBIT 3Y Avg:", FormatMultiple(r.AvgDepFactor3Y)), item("Operating Margin FY:", FormatPercent(r.OperatingMargin))))
		fmt.Println(cardRow(item("EBITDA/EBIT TTM:", FormatMultiple(r.DepFactorTTM)), item("Operating Margin 3Y Avg:", FormatPercent(r.OperatingMargin3Y))))
		
		depHistory := fmt.Sprintf("%s > %s > %s", FormatMultiple(r.DepFactorY3), FormatMultiple(r.DepFactorY2), FormatMultiple(r.DepFactorY1))
		fmt.Println(cardRow(item("3Y Dep Trajectory:", depHistory), item("EBIT Margin 3Y Trend:", formatTrend(r.EBITMarginTrend3Y, plain))))
		fmt.Println(cardRow(item("ROIC Latest FY:", FormatPercent(r.ROICLatest)), item("Revenue 3Y CAGR:", FormatPercent(r.RevCAGR3Y))))
		fmt.Println(cardRow(item("ROIC 3Y Avg:", FormatPercent(r.ROICAvg3Y)), item("EBIT 3Y CAGR:", FormatPercent(r.EBITCAGR3Y))))
		fmt.Println(cardRow(item("ROIC 3Y Trend:", formatTrend(r.ROICTrend3Y, plain)), item("Shares 3Y CAGR:", FormatPercent(r.ShareCountCAGR3Y))))
		fmt.Println(cardRow("", item("Revenue Stability:", fmt.Sprintf("%s (StdDev)", FormatPercent(r.RevenueStability)))))

		fmt.Println("├" + strings.Repeat("─", 44) + "┼" + strings.Repeat("─", 43) + "┤")

		// QUADRANT 3 & 4
		sec3 := yellowBold("3. CASH FLOW CONVERSION (20%)")
		sec4 := yellowBold("4. VALUATION & CAPITAL STRUCTURE (25%)")
		if plain {
			sec3 = "3. CASH FLOW CONVERSION (20%)"
			sec4 = "4. VALUATION & CAPITAL STRUCTURE (25%)"
		}
		fmt.Println(cardRow(sec3, sec4))
		fmt.Println(cardRow(item("Adj FCF / Net Income FY:", FormatPercent(r.AdjFCFNetIncomeLatest)), item("EV / EBIT Multiple:", formatMultipleX(r.EVEBIT))))
		fmt.Println(cardRow(item("Adj FCF / NI 3Y Avg:", FormatPercent(r.AdjFCFNetIncomeRatio)), item("EV / EBITDA Multiple:", formatMultipleX(r.EVEBITDA))))
		fmt.Println(cardRow(item("Adj FCF Margin 3Y Avg:", FormatPercent(r.AdjFCFMargin)), item("TTM Adj FCF Yield:", FormatPercent(r.TTMAdjFCFYield))))
		fmt.Println(cardRow(item("CapEx / Revenue FY:", FormatPercent(r.CapExRevLatest)), item("3Y Sum AdjFCF / EV:", FormatPercent(r.AdjFCF3YEVYield))))
		fmt.Println(cardRow(item("CapEx / Revenue 3Y Avg:", FormatPercent(r.CapExRevRatio)), item("Net Debt / EBITDA:", formatMultipleX(r.NetDebtEBITDA))))
		fmt.Println(cardRow(item("SBC / Revenue FY:", FormatPercent(r.SBCRevLatest)), ""))
		fmt.Println(cardRow(item("SBC / Revenue 3Y Avg:", FormatPercent(r.SBCRevRatio)), ""))
		fmt.Println(cardRow(item("SBC / FCF 3Y Avg:", FormatPercent(r.SBCFCFRatio)), ""))

		// Bottom Verdict Section (Strengths & Watch-Outs)
		fmt.Println("├" + strings.Repeat("─", 88) + "┤")
		strLine, cautLine := generateInsights(r, plain)
		fmt.Printf("│ %s │\n", padRight(" "+strLine, 86))
		fmt.Printf("│ %s │\n", padRight(" "+cautLine, 86))

		fmt.Println("└" + strings.Repeat("─", 88) + "┘")
	}
}

// PrintGroupedTables displays metrics partitioned across clean thematic tables.
func PrintGroupedTables(results []*metrics.ScreenResult, plain bool) {
	fmt.Println(yellowBold("═══ 1. COMPOSITE SCORES & CAPITAL INTENSITY ═══"))
	cols1 := []ColumnDef{
		{Key: "Ticker", Label: "Ticker"},
		{Key: "CompanyName", Label: "Company"},
		{Key: "Score", Label: "Total Score", Numeric: true},
		{Key: "ReturnsScore", Label: "Returns Score", Numeric: true},
		{Key: "GrowthScore", Label: "Growth Score", Numeric: true},
		{Key: "CashScore", Label: "Cash Score", Numeric: true},
		{Key: "ValuationScore", Label: "Val Score", Numeric: true},
		{Key: "AvgDepFactor3Y", Label: "EBITDA/EBIT 3Y Avg", Numeric: true},
		{Key: "DepFactorLatest", Label: "EBITDA/EBIT FY", Numeric: true},
		{Key: "IntensityLevel", Label: "Intensity"},
	}
	renderSingleTable(cols1, results, plain)

	fmt.Printf("\n%s\n", yellowBold("═══ 2. RETURNS, MARGINS & GROWTH ═══"))
	cols2 := []ColumnDef{
		{Key: "Ticker", Label: "Ticker"},
		{Key: "ROICLatest", Label: "ROIC FY", Numeric: true},
		{Key: "ROICAvg3Y", Label: "ROIC 3Y Avg", Numeric: true},
		{Key: "ROICTrend3Y", Label: "ROIC 3Y Trend", Numeric: true},
		{Key: "GrossMargin", Label: "Gross Mrgn", Numeric: true},
		{Key: "OperatingMargin", Label: "EBIT Mrgn FY", Numeric: true},
		{Key: "OperatingMargin3Y", Label: "EBIT Mrgn 3Y Avg", Numeric: true},
		{Key: "EBITMarginTrend3Y", Label: "EBIT 3Y Trend", Numeric: true},
		{Key: "RevCAGR3Y", Label: "Rev 3Y CAGR", Numeric: true},
		{Key: "EBITCAGR3Y", Label: "EBIT 3Y CAGR", Numeric: true},
		{Key: "ShareCountCAGR3Y", Label: "Shares 3Y CAGR", Numeric: true},
	}
	renderSingleTable(cols2, results, plain)

	fmt.Printf("\n%s\n", yellowBold("═══ 3. CASH FLOW CONVERSION & REINVESTMENT ═══"))
	cols3 := []ColumnDef{
		{Key: "Ticker", Label: "Ticker"},
		{Key: "AdjFCFNetIncomeLatest", Label: "AdjFCF/NI FY", Numeric: true},
		{Key: "AdjFCFNetIncomeRatio", Label: "AdjFCF/NI 3Y Avg", Numeric: true},
		{Key: "AdjFCFMargin", Label: "AdjFCF Mrgn 3Y Avg", Numeric: true},
		{Key: "CapExRevLatest", Label: "CapEx/Rev FY", Numeric: true},
		{Key: "CapExRevRatio", Label: "CapEx/Rev 3Y Avg", Numeric: true},
		{Key: "SBCRevRatio", Label: "SBC/Rev 3Y Avg", Numeric: true},
		{Key: "SBCFCFRatio", Label: "SBC/FCF 3Y Avg", Numeric: true},
	}
	renderSingleTable(cols3, results, plain)

	fmt.Printf("\n%s\n", yellowBold("═══ 4. VALUATION & BALANCE SHEET ═══"))
	cols4 := []ColumnDef{
		{Key: "Ticker", Label: "Ticker"},
		{Key: "EVEBIT", Label: "EV/EBIT", Numeric: true},
		{Key: "EVEBITDA", Label: "EV/EBITDA", Numeric: true},
		{Key: "TTMAdjFCFYield", Label: "TTM AdjFCF Yld", Numeric: true},
		{Key: "AdjFCF3YEVYield", Label: "3Y Sum AdjFCF/EV", Numeric: true},
		{Key: "AdjFCFYield", Label: "AdjFCF Yld", Numeric: true},
		{Key: "NetDebtEBITDA", Label: "Debt/EBITDA", Numeric: true},
		{Key: "Currency", Label: "Currency"},
	}
	renderSingleTable(cols4, results, plain)
}

func PrintTable(results []*metrics.ScreenResult, plain bool) {
	renderSingleTable(SummaryColumns, results, plain)
}

func PrintSummary(results []*metrics.ScreenResult, errors []string) {
	if len(results) == 0 {
		fmt.Println("No valid companies returned.")
		if len(errors) > 0 {
			fmt.Println("\nSkipped / warnings:")
			for _, err := range errors {
				fmt.Printf("  - %s\n", err)
			}
		}
		return
	}

	low, mod, high := 0, 0, 0
	for _, r := range results {
		switch r.IntensityLevel {
		case "Low":
			low++
		case "Moderate":
			mod++
		case "High":
			high++
		}
	}

	fmt.Printf("\nSummary: %d stocks screened | %d Low Intensity (<1.4) | %d Moderate (1.4-1.7) | %d High (>1.7)\n",
		len(results), low, mod, high)
	fmt.Println("Options: Use --table for compact table; --grouped for themed tables; --export for CSV/Excel/MD.")
	fmt.Println("How to read: EBITDA/EBIT is lower-is-better capital intensity; AdjFCF/NI is cash conversion.")
	fmt.Println("Score is 0-100: 80+ strong, 70-80 good, 60-70 mixed, <60 weak.")

	if len(errors) > 0 {
		fmt.Println("\nSkipped / warnings:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
	}
}
