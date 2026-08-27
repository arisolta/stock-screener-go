package ui

import (
	"fmt"
	"math"

	"stock-screener/internal/metrics"
)

type ColumnDef struct {
	Key     string
	Label   string
	Header  string
	Numeric bool
}

// TerminalColumns contains all calculated metrics and scores for complete display.
var TerminalColumns = []ColumnDef{
	{Key: "Ticker", Label: "Ticker", Header: "Ticker", Numeric: false},
	{Key: "CompanyName", Label: "Company", Header: "Company", Numeric: false},
	{Key: "Currency", Label: "Cur", Header: "Currency", Numeric: false},
	{Key: "DepFactorLatest", Label: "EBITDA/EBIT FY", Header: "EBITDA/EBIT (Dep Factor) FY", Numeric: true},
	{Key: "AvgDepFactor3Y", Label: "EBITDA/EBIT 3Y Avg", Header: "EBITDA/EBIT 3Y Avg", Numeric: true},
	{Key: "DepFactorTTM", Label: "EBITDA/EBIT TTM", Header: "EBITDA/EBIT TTM", Numeric: true},
	{Key: "ROICLatest", Label: "ROIC FY", Header: "ROIC FY", Numeric: true},
	{Key: "ROICAvg3Y", Label: "ROIC 3Y Avg", Header: "ROIC 3Y Avg", Numeric: true},
	{Key: "ROICTrend3Y", Label: "ROIC 3Y Trend", Header: "ROIC 3Y Trend", Numeric: true},
	{Key: "GrossMargin", Label: "Gross Mrgn", Header: "Gross Margin FY", Numeric: true},
	{Key: "OperatingMargin", Label: "EBIT Mrgn FY", Header: "Operating Margin FY", Numeric: true},
	{Key: "OperatingMargin3Y", Label: "EBIT Mrgn 3Y Avg", Header: "Operating Margin 3Y Avg", Numeric: true},
	{Key: "EBITMarginTrend3Y", Label: "EBIT 3Y Trend", Header: "EBIT Margin 3Y Trend", Numeric: true},
	{Key: "RevCAGR3Y", Label: "Rev 3Y CAGR", Header: "Revenue 3Y CAGR", Numeric: true},
	{Key: "EBITCAGR3Y", Label: "EBIT 3Y CAGR", Header: "EBIT 3Y CAGR", Numeric: true},
	{Key: "ShareCountCAGR3Y", Label: "Shares 3Y CAGR", Header: "Shares 3Y CAGR", Numeric: true},
	{Key: "RevenueStability", Label: "Rev Stability", Header: "Revenue Stability (StdDev)", Numeric: true},
	{Key: "AdjFCFNetIncomeLatest", Label: "AdjFCF/NI FY", Header: "Adj FCF / Net Income FY", Numeric: true},
	{Key: "AdjFCFNetIncomeRatio", Label: "AdjFCF/NI 3Y Avg", Header: "Adj FCF / Net Income 3Y Avg", Numeric: true},
	{Key: "AdjFCFEBITDARatio", Label: "AdjFCF/EBITDA 3Y Avg", Header: "Adj FCF / EBITDA 3Y Avg", Numeric: true},
	{Key: "AdjFCFMargin", Label: "AdjFCF Mrgn 3Y Avg", Header: "Adj FCF Margin 3Y Avg", Numeric: true},
	{Key: "FCFEBITDARatio", Label: "FCF/EBITDA 3Y Avg", Header: "FCF / EBITDA 3Y Avg", Numeric: true},
	{Key: "FCFMargin", Label: "FCF Mrgn 3Y Avg", Header: "FCF Margin 3Y Avg", Numeric: true},
	{Key: "CapExRevLatest", Label: "CapEx/Rev FY", Header: "CapEx / Revenue FY", Numeric: true},
	{Key: "CapExRevRatio", Label: "CapEx/Rev 3Y Avg", Header: "CapEx / Revenue 3Y Avg", Numeric: true},
	{Key: "SBCRevLatest", Label: "SBC/Rev FY", Header: "SBC / Revenue FY", Numeric: true},
	{Key: "SBCRevRatio", Label: "SBC/Rev 3Y Avg", Header: "SBC / Revenue 3Y Avg", Numeric: true},
	{Key: "SBCFCFRatio", Label: "SBC/FCF 3Y Avg", Header: "SBC / FCF 3Y Avg", Numeric: true},
	{Key: "DAEBITRatio", Label: "D&A/EBIT 3Y Avg", Header: "D&A / EBIT 3Y Avg", Numeric: true},
	{Key: "EVEBIT", Label: "EV/EBIT", Header: "EV / EBIT Multiple", Numeric: true},
	{Key: "EVEBITDA", Label: "EV/EBITDA", Header: "EV / EBITDA Multiple", Numeric: true},
	{Key: "TTMFCFYield", Label: "TTM FCF Yld", Header: "TTM FCF Yield", Numeric: true},
	{Key: "TTMAdjFCFYield", Label: "TTM AdjFCF Yld", Header: "TTM Adj FCF Yield", Numeric: true},
	{Key: "FCFYield", Label: "FCF Yld", Header: "FCF Yield", Numeric: true},
	{Key: "AdjFCFYield", Label: "AdjFCF Yld", Header: "Adj FCF Yield", Numeric: true},
	{Key: "FCF3YEVYield", Label: "3Y Sum FCF/EV", Header: "3Y Sum FCF / EV Yield", Numeric: true},
	{Key: "AdjFCF3YEVYield", Label: "3Y Sum AdjFCF/EV", Header: "3Y Sum Adj FCF / EV Yield", Numeric: true},
	{Key: "ReturnsScore", Label: "Returns Score", Header: "Returns Score", Numeric: true},
	{Key: "GrowthScore", Label: "Growth Score", Header: "Growth Score", Numeric: true},
	{Key: "CashScore", Label: "Cash Score", Header: "Cash Score", Numeric: true},
	{Key: "ValuationScore", Label: "Val Score", Header: "Valuation Score", Numeric: true},
	{Key: "Score", Label: "Score", Header: "Composite Score", Numeric: true},
	{Key: "IntensityLevel", Label: "Intensity", Header: "Capital Intensity", Numeric: false},
}

// AllExportColumns is an alias for export modules
var AllExportColumns = TerminalColumns

func FormatPercent(v float64) string {
	if math.IsNaN(v) {
		return ""
	}
	return fmt.Sprintf("%.1f%%", v*100)
}

func FormatMultiple(v float64) string {
	if math.IsNaN(v) {
		return ""
	}
	return fmt.Sprintf("%.2f", v)
}

func FormatScore(v float64) string {
	if math.IsNaN(v) {
		return ""
	}
	return fmt.Sprintf("%.1f", v)
}

// GetFormattedValue returns the formatted string and raw value for a given field key on ScreenResult.
func GetFormattedValue(r *metrics.ScreenResult, key string) string {
	switch key {
	case "Ticker":
		return r.Ticker
	case "CompanyName":
		return r.CompanyName
	case "Currency":
		return r.Currency
	case "IntensityLevel":
		return r.IntensityLevel
	case "DepFactorLatest":
		return FormatMultiple(r.DepFactorLatest)
	case "AvgDepFactor3Y":
		return FormatMultiple(r.AvgDepFactor3Y)
	case "DepFactorTTM":
		return FormatMultiple(r.DepFactorTTM)
	case "DAEBITRatio":
		return FormatPercent(r.DAEBITRatio)
	case "RevCAGR3Y":
		return FormatPercent(r.RevCAGR3Y)
	case "EBITCAGR3Y":
		return FormatPercent(r.EBITCAGR3Y)
	case "ROICLatest":
		return FormatPercent(r.ROICLatest)
	case "ROICAvg3Y":
		return FormatPercent(r.ROICAvg3Y)
	case "ROICTrend3Y":
		return FormatPercent(r.ROICTrend3Y)
	case "GrossMargin":
		return FormatPercent(r.GrossMargin)
	case "OperatingMargin":
		return FormatPercent(r.OperatingMargin)
	case "OperatingMargin3Y":
		return FormatPercent(r.OperatingMargin3Y)
	case "EBITMarginTrend3Y":
		return FormatPercent(r.EBITMarginTrend3Y)
	case "FCFEBITDARatio":
		return FormatPercent(r.FCFEBITDARatio)
	case "FCFMargin":
		return FormatPercent(r.FCFMargin)
	case "AdjFCFEBITDARatio":
		return FormatPercent(r.AdjFCFEBITDARatio)
	case "AdjFCFNetIncomeLatest":
		return FormatPercent(r.AdjFCFNetIncomeLatest)
	case "AdjFCFNetIncomeRatio":
		return FormatPercent(r.AdjFCFNetIncomeRatio)
	case "AdjFCFMargin":
		return FormatPercent(r.AdjFCFMargin)
	case "SBCRevLatest":
		return FormatPercent(r.SBCRevLatest)
	case "SBCRevRatio":
		return FormatPercent(r.SBCRevRatio)
	case "SBCFCFRatio":
		return FormatPercent(r.SBCFCFRatio)
	case "CapExRevLatest":
		return FormatPercent(r.CapExRevLatest)
	case "CapExRevRatio":
		return FormatPercent(r.CapExRevRatio)
	case "NetDebtEBITDA":
		return FormatMultiple(r.NetDebtEBITDA)
	case "ShareCountCAGR3Y":
		return FormatPercent(r.ShareCountCAGR3Y)
	case "RevenueStability":
		return FormatPercent(r.RevenueStability)
	case "EVEBITDA":
		return FormatMultiple(r.EVEBITDA)
	case "EVEBIT":
		return FormatMultiple(r.EVEBIT)
	case "TTMFCFYield":
		return FormatPercent(r.TTMFCFYield)
	case "TTMAdjFCFYield":
		return FormatPercent(r.TTMAdjFCFYield)
	case "FCFYield":
		return FormatPercent(r.FCFYield)
	case "AdjFCFYield":
		return FormatPercent(r.AdjFCFYield)
	case "FCF3YEVYield":
		return FormatPercent(r.FCF3YEVYield)
	case "AdjFCF3YEVYield":
		return FormatPercent(r.AdjFCF3YEVYield)
	case "DepFactorY1":
		return FormatMultiple(r.DepFactorY1)
	case "DepFactorY2":
		return FormatMultiple(r.DepFactorY2)
	case "DepFactorY3":
		return FormatMultiple(r.DepFactorY3)
	case "FiscalY1":
		return r.FiscalY1
	case "FiscalY2":
		return r.FiscalY2
	case "FiscalY3":
		return r.FiscalY3
	case "ReturnsScore":
		return FormatScore(r.ReturnsScore)
	case "GrowthScore":
		return FormatScore(r.GrowthScore)
	case "CashScore":
		return FormatScore(r.CashScore)
	case "ValuationScore":
		return FormatScore(r.ValuationScore)
	case "Score":
		return FormatScore(r.Score)
	default:
		return ""
	}
}
