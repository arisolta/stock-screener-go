package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"stock-screener/internal/cache"
	"stock-screener/internal/config"
	"stock-screener/internal/export"
	"stock-screener/internal/metrics"
	"stock-screener/internal/ui"
	"stock-screener/internal/yfinance"
)

var higherIsBetter = map[string]bool{
	"ROIC_Latest":               true,
	"ROIC_Avg_3Y":               true,
	"ROIC_Trend_3Y":             true,
	"Gross_Margin":              true,
	"Operating_Margin":          true,
	"Operating_Margin_3Y":       true,
	"EBIT_Margin_Trend_3Y":      true,
	"FCF_EBITDA_Ratio":          true,
	"FCF_Margin":                true,
	"Adj_FCF_EBITDA_Ratio":      true,
	"Adj_FCF_Net_Income_Latest": true,
	"Adj_FCF_Net_Income_Ratio":  true,
	"Adj_FCF_Margin":            true,
	"TTM_FCF_Yield":             true,
	"TTM_Adj_FCF_Yield":         true,
	"FCF_Yield":                 true,
	"Adj_FCF_Yield":             true,
	"FCF_3Y_EV_Yield":           true,
	"Adj_FCF_3Y_EV_Yield":       true,
	"Returns_Score":             true,
	"Growth_Score":              true,
	"Cash_Score":                true,
	"Valuation_Score":           true,
	"Score":                     true,
}

func getRawMetricValue(r *metrics.ScreenResult, key string) (float64, bool) {
	normKey := strings.ToLower(strings.ReplaceAll(key, "_", ""))
	switch normKey {
	case "score":
		return r.Score, true
	case "returnsscore", "capitalscore", "qualityscore":
		return r.ReturnsScore, true
	case "growthscore":
		return r.GrowthScore, true
	case "cashscore":
		return r.CashScore, true
	case "valuationscore":
		return r.ValuationScore, true
	case "depfactorlatest", "depfactor":
		return r.DepFactorLatest, true
	case "avgdepfactor3y", "depfactor3y":
		return r.AvgDepFactor3Y, true
	case "depfactorttm":
		return r.DepFactorTTM, true
	case "daebitratio":
		return r.DAEBITRatio, true
	case "revcagr3y", "revcagr":
		return r.RevCAGR3Y, true
	case "ebitcagr3y", "ebitcagr":
		return r.EBITCAGR3Y, true
	case "roiclatest", "roic":
		return r.ROICLatest, true
	case "roicavg3y", "roic3y":
		return r.ROICAvg3Y, true
	case "roictrend3y":
		return r.ROICTrend3Y, true
	case "grossmargin":
		return r.GrossMargin, true
	case "operatingmargin", "ebitmargin":
		return r.OperatingMargin, true
	case "operatingmargin3y", "ebitmargin3y":
		return r.OperatingMargin3Y, true
	case "ebitmargintrend3y":
		return r.EBITMarginTrend3Y, true
	case "fcfebitdaratio":
		return r.FCFEBITDARatio, true
	case "fcfmargin":
		return r.FCFMargin, true
	case "adjfcfebitdaratio":
		return r.AdjFCFEBITDARatio, true
	case "adjfcfnetincomelatest", "adjfcfnilatest":
		return r.AdjFCFNetIncomeLatest, true
	case "adjfcfnetincomeratio", "adjfcfni3y":
		return r.AdjFCFNetIncomeRatio, true
	case "adjfcfmargin":
		return r.AdjFCFMargin, true
	case "sbcrevlatest":
		return r.SBCRevLatest, true
	case "sbcrevratio":
		return r.SBCRevRatio, true
	case "sbcclosingratio", "sbcclosing":
		return r.SBCFCFRatio, true
	case "capexrevlatest":
		return r.CapExRevLatest, true
	case "capexrevratio":
		return r.CapExRevRatio, true
	case "netdebtebitda", "debtebitda":
		return r.NetDebtEBITDA, true
	case "sharecountcagr3y":
		return r.ShareCountCAGR3Y, true
	case "revenuestability":
		return r.RevenueStability, true
	case "evebitda":
		return r.EVEBITDA, true
	case "evebit":
		return r.EVEBIT, true
	case "ttmfcfyield":
		return r.TTMFCFYield, true
	case "ttmadjfcfyield":
		return r.TTMAdjFCFYield, true
	case "fcfyield":
		return r.FCFYield, true
	case "adjfcfyield":
		return r.AdjFCFYield, true
	case "fcf3yevyield":
		return r.FCF3YEVYield, true
	case "adjfcf3yevyield":
		return r.AdjFCF3YEVYield, true
	default:
		return 0, false
	}
}

func main() {
	cfg, err := config.ParseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(2)
	}

	diskCache := cache.New("", 24*time.Hour)
	yfClient := yfinance.NewClient()

	var results []*metrics.ScreenResult
	var errors []string
	var mu sync.Mutex

	// Worker pool for concurrent fetching
	tickerChan := make(chan string, len(cfg.Tickers))
	for _, t := range cfg.Tickers {
		tickerChan <- t
	}
	close(tickerChan)

	concurrency := cfg.Concurrency
	if concurrency <= 0 {
		concurrency = 5
	}
	if concurrency > len(cfg.Tickers) {
		concurrency = len(cfg.Tickers)
	}

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ticker := range tickerChan {
				var rawData *yfinance.RawFinancialData
				var err error

				if !cfg.NoCache {
					if cached, ok := diskCache.Get(ticker); ok {
						rawData = cached
					}
				}

				if rawData == nil {
					rawData, err = yfClient.GetFinancials(ticker)
					if err != nil {
						mu.Lock()
						errors = append(errors, fmt.Sprintf("%s: %v", ticker, err))
						mu.Unlock()
						continue
					}
					if !cfg.NoCache {
						_ = diskCache.Set(ticker, rawData)
					}
				}

				res, err := metrics.CalculateMetrics(rawData, cfg.Years)
				if err != nil {
					mu.Lock()
					errors = append(errors, fmt.Sprintf("%s: %v", ticker, err))
					mu.Unlock()
					continue
				}

				mu.Lock()
				results = append(results, res)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if len(results) == 0 {
		ui.PrintSummary(results, errors)
		os.Exit(1)
	}

	// Apply filtering
	var filtered []*metrics.ScreenResult
	for _, r := range results {
		if cfg.MinROIC != nil {
			threshold := *cfg.MinROIC / 100.0
			if math.IsNaN(r.ROICLatest) || r.ROICLatest < threshold {
				continue
			}
		}
		if cfg.MaxDepFactor != nil {
			if math.IsNaN(r.AvgDepFactor3Y) || r.AvgDepFactor3Y > *cfg.MaxDepFactor {
				continue
			}
		}
		filtered = append(filtered, r)
	}
	results = filtered

	// Sort results
	rankKey := cfg.RankBy
	higherBetter := higherIsBetter[rankKey]
	sort.SliceStable(results, func(i, j int) bool {
		valI, okI := getRawMetricValue(results[i], rankKey)
		valJ, okJ := getRawMetricValue(results[j], rankKey)
		if !okI || !okJ {
			return results[i].Score > results[j].Score
		}

		if math.IsNaN(valI) && math.IsNaN(valJ) {
			return results[i].Score > results[j].Score
		}
		if math.IsNaN(valI) {
			return false
		}
		if math.IsNaN(valJ) {
			return true
		}

		if higherBetter {
			return valI > valJ
		}
		return valI < valJ
	})

	// Render view based on config
	switch cfg.View {
	case "cards":
		ui.PrintCardView(results, cfg.Plain)
	case "grouped":
		ui.PrintGroupedTables(results, cfg.Plain)
	default:
		ui.PrintTable(results, cfg.Plain)
	}

	ui.PrintSummary(results, errors)

	if cfg.ExportPath != "" {
		if err := export.ExportResults(results, cfg.ExportPath); err != nil {
			fmt.Fprintf(os.Stderr, "Export failed: %v\n", err)
			os.Exit(3)
		}
		fmt.Printf("\nExported results to %s\n", cfg.ExportPath)
	}
}
