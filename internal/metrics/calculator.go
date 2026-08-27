package metrics

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"stock-screener/internal/yfinance"
)

func SafeDivide(n, d float64) float64 {
	if math.IsNaN(n) || math.IsNaN(d) || d == 0 {
		return math.NaN()
	}
	return n / d
}

func SafePositiveDivide(n, d float64) float64 {
	if math.IsNaN(d) || d <= 0 {
		return math.NaN()
	}
	return SafeDivide(n, d)
}

func NanMean(vals []float64) float64 {
	var sum float64
	var count int
	for _, v := range vals {
		if !math.IsNaN(v) {
			sum += v
			count++
		}
	}
	if count == 0 {
		return math.NaN()
	}
	return sum / float64(count)
}

func NanStd(vals []float64) float64 {
	var valid []float64
	for _, v := range vals {
		if !math.IsNaN(v) {
			valid = append(valid, v)
		}
	}
	if len(valid) == 0 {
		return math.NaN()
	}
	mean := NanMean(valid)
	var sumSq float64
	for _, v := range valid {
		diff := v - mean
		sumSq += diff * diff
	}
	return math.Sqrt(sumSq / float64(len(valid)))
}

func CAGR(start, end float64, periods int) float64 {
	if periods <= 0 || math.IsNaN(start) || math.IsNaN(end) || start <= 0 || end <= 0 {
		return math.NaN()
	}
	return math.Pow(end/start, 1.0/float64(periods)) - 1.0
}

func IntensityLevel(avgDepFactor float64) string {
	if math.IsNaN(avgDepFactor) {
		return "Unknown"
	}
	if avgDepFactor < 1.4 {
		return "Low"
	}
	if avgDepFactor <= 1.7 {
		return "Moderate"
	}
	return "High"
}

func getMetricValue(statement map[string]map[string]float64, keys []string, period string) float64 {
	for _, k := range keys {
		if m, ok := statement[k]; ok {
			if v, ok2 := m[period]; ok2 {
				return v
			}
		}
	}
	return math.NaN()
}

func getAllPeriods(statement map[string]map[string]float64) []string {
	periodSet := make(map[string]struct{})
	for _, m := range statement {
		for p := range m {
			periodSet[p] = struct{}{}
		}
	}
	periods := make([]string, 0, len(periodSet))
	for p := range periodSet {
		periods = append(periods, p)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(periods)))
	return periods
}

type TTMData struct {
	Revenue float64
	EBIT    float64
	EBITDA  float64
	FCF     float64
	SBC     float64
	AdjFCF  float64
	CapEx   float64
	Valid   bool
}

func calculateTTM(quarterly map[string]map[string]float64) TTMData {
	periods := getAllPeriods(quarterly)
	if len(periods) < 4 {
		return TTMData{
			Revenue: math.NaN(), EBIT: math.NaN(), EBITDA: math.NaN(),
			FCF: math.NaN(), SBC: math.NaN(), AdjFCF: math.NaN(), CapEx: math.NaN(),
			Valid: false,
		}
	}
	usePeriods := periods[:4]

	sumLine := func(keys []string) float64 {
		var sum float64
		var count int
		for _, p := range usePeriods {
			v := getMetricValue(quarterly, keys, p)
			if !math.IsNaN(v) {
				sum += v
				count++
			}
		}
		if count == 0 {
			return math.NaN()
		}
		return sum
	}

	rev := sumLine([]string{"quarterlyTotalRevenue", "quarterlyOperatingRevenue"})
	ebit := sumLine([]string{"quarterlyEBIT", "quarterlyOperatingIncome"})
	ebitda := sumLine([]string{"quarterlyEBITDA"})
	if math.IsNaN(ebitda) && !math.IsNaN(ebit) {
		dep := sumLine([]string{"quarterlyReconciledDepreciation", "quarterlyDepreciationAndAmortization"})
		if !math.IsNaN(dep) {
			ebitda = ebit + dep
		}
	}
	ocf := sumLine([]string{"quarterlyOperatingCashFlow", "quarterlyTotalCashFromOperatingActivities"})
	capex := sumLine([]string{"quarterlyCapitalExpenditure", "quarterlyCapitalExpenditures"})
	sbc := sumLine([]string{"quarterlyStockBasedCompensation", "quarterlyStockBasedCompensationAndOther"})

	fcf := math.NaN()
	if !math.IsNaN(ocf) {
		capexVal := 0.0
		if !math.IsNaN(capex) {
			capexVal = capex
		}
		fcf = ocf + capexVal
	}

	adjFCF := fcf
	if !math.IsNaN(fcf) && !math.IsNaN(sbc) {
		adjFCF = fcf - sbc
	}

	return TTMData{
		Revenue: rev,
		EBIT:    ebit,
		EBITDA:  ebitda,
		FCF:     fcf,
		SBC:     sbc,
		AdjFCF:  adjFCF,
		CapEx:   capex,
		Valid:   true,
	}
}

func CalculateMetrics(raw *yfinance.RawFinancialData, years int) (*ScreenResult, error) {
	if years <= 0 {
		years = 3
	}

	allPeriods := getAllPeriods(raw.Annual)
	if len(allPeriods) < 2 {
		return nil, fmt.Errorf("not enough annual financial statement history (%d periods found)", len(allPeriods))
	}

	taxRate := 0.21
	if raw.Info.EffectiveTaxRate != nil {
		tr := *raw.Info.EffectiveTaxRate
		if !math.IsNaN(tr) && tr >= 0 && tr <= 0.5 {
			taxRate = tr
		}
	}

	var annualRows []AnnualRow
	for _, p := range allPeriods {
		rev := getMetricValue(raw.Annual, []string{"annualTotalRevenue", "annualOperatingRevenue"}, p)
		gp := getMetricValue(raw.Annual, []string{"annualGrossProfit"}, p)
		ebit := getMetricValue(raw.Annual, []string{"annualEBIT", "annualOperatingIncome"}, p)
		ni := getMetricValue(raw.Annual, []string{"annualNetIncome", "annualNetIncomeCommonStockholders"}, p)
		ebitda := getMetricValue(raw.Annual, []string{"annualEBITDA"}, p)
		dep := getMetricValue(raw.Annual, []string{"annualReconciledDepreciation", "annualDepreciationAndAmortization", "annualDepreciation"}, p)
		if math.IsNaN(ebitda) && !math.IsNaN(ebit) && !math.IsNaN(dep) {
			ebitda = ebit + dep
		}

		ocf := getMetricValue(raw.Annual, []string{"annualOperatingCashFlow", "annualTotalCashFromOperatingActivities"}, p)
		capex := getMetricValue(raw.Annual, []string{"annualCapitalExpenditure", "annualCapitalExpenditures"}, p)
		sbc := getMetricValue(raw.Annual, []string{"annualStockBasedCompensation", "annualStockBasedCompensationAndOther", "annualShareBasedCompensation", "annualShareBasedCompensationExpense"}, p)

		fcf := math.NaN()
		if !math.IsNaN(ocf) {
			capexVal := 0.0
			if !math.IsNaN(capex) {
				capexVal = capex
			}
			fcf = ocf + capexVal
		}

		adjFCF := fcf
		if !math.IsNaN(fcf) && !math.IsNaN(sbc) {
			adjFCF = fcf - sbc
		}

		shares := getMetricValue(raw.Annual, []string{"annualDilutedAverageShares", "annualBasicAverageShares", "annualAverageDilutionEarnings"}, p)
		eq := getMetricValue(raw.Annual, []string{"annualStockholdersEquity", "annualTotalStockholderEquity"}, p)
		debt := getMetricValue(raw.Annual, []string{"annualTotalDebt", "annualLongTermDebtAndCapitalLeaseObligation", "annualLongTermDebt"}, p)
		cash := getMetricValue(raw.Annual, []string{"annualCashAndCashEquivalents", "annualCashCashEquivalentsAndShortTermInvestments"}, p)

		// Invested Capital = Equity + Debt - Cash
		investedCap := math.NaN()
		if !math.IsNaN(eq) || !math.IsNaN(debt) || !math.IsNaN(cash) {
			eqVal := 0.0
			debtVal := 0.0
			cashVal := 0.0
			if !math.IsNaN(eq) {
				eqVal = eq
			}
			if !math.IsNaN(debt) {
				debtVal = debt
			}
			if !math.IsNaN(cash) {
				cashVal = cash
			}
			investedCap = eqVal + debtVal - cashVal
		}

		roic := math.NaN()
		if !math.IsNaN(ebit) && !math.IsNaN(investedCap) && investedCap > 0 {
			nopat := ebit * (1.0 - taxRate)
			roic = SafeDivide(nopat, investedCap)
		}

		netDebt := math.NaN()
		if !math.IsNaN(debt) || !math.IsNaN(cash) {
			debtVal := 0.0
			cashVal := 0.0
			if !math.IsNaN(debt) {
				debtVal = debt
			}
			if !math.IsNaN(cash) {
				cashVal = cash
			}
			netDebt = debtVal - cashVal
		}

		depreciationVal := math.NaN()
		if !math.IsNaN(ebitda) && !math.IsNaN(ebit) {
			depreciationVal = ebitda - ebit
		}

		annualRows = append(annualRows, AnnualRow{
			Period:          p,
			Revenue:         rev,
			GrossProfit:     gp,
			EBIT:            ebit,
			NetIncome:       ni,
			EBITDA:          ebitda,
			Depreciation:    depreciationVal,
			FCF:             fcf,
			SBC:             sbc,
			AdjFCF:          adjFCF,
			CapEx:           math.Abs(capex),
			AverageShares:   shares,
			GrossMargin:     SafePositiveDivide(gp, rev),
			OperatingMargin: SafePositiveDivide(ebit, rev),
			DepFactor:       SafePositiveDivide(ebitda, ebit),
			DAToEBIT:        SafePositiveDivide(depreciationVal, ebit),
			FCFEBITDA:       SafePositiveDivide(fcf, ebitda),
			FCFMargin:       SafePositiveDivide(fcf, rev),
			AdjFCFEBITDA:    SafePositiveDivide(adjFCF, ebitda),
			AdjFCFNetIncome: SafePositiveDivide(adjFCF, ni),
			AdjFCFMargin:    SafePositiveDivide(adjFCF, rev),
			SBCRevenue:      SafePositiveDivide(sbc, rev),
			SBCFCF:          SafePositiveDivide(sbc, fcf),
			CapExRevenue:    SafePositiveDivide(math.Abs(capex), rev),
			NetDebt:         netDebt,
			ROIC:            roic,
		})
	}

	latest := annualRows[0]
	numAvg := years
	if len(annualRows) < numAvg {
		numAvg = len(annualRows)
	}
	avgRows := annualRows[:numAvg]

	growthIntervals := years
	if len(annualRows)-1 < growthIntervals {
		growthIntervals = len(annualRows) - 1
	}
	oldestGrowth := annualRows[growthIntervals]

	revCAGR := CAGR(oldestGrowth.Revenue, latest.Revenue, growthIntervals)
	ebitCAGR := CAGR(oldestGrowth.EBIT, latest.EBIT, growthIntervals)
	sharesCAGR := CAGR(oldestGrowth.AverageShares, latest.AverageShares, growthIntervals)

	roicTrend := math.NaN()
	if !math.IsNaN(latest.ROIC) && !math.IsNaN(oldestGrowth.ROIC) {
		roicTrend = latest.ROIC - oldestGrowth.ROIC
	}

	ebitMarginTrend := math.NaN()
	if !math.IsNaN(latest.OperatingMargin) && !math.IsNaN(oldestGrowth.OperatingMargin) {
		ebitMarginTrend = latest.OperatingMargin - oldestGrowth.OperatingMargin
	}

	// Revenue stability
	chronological := make([]AnnualRow, growthIntervals+1)
	for i := 0; i <= growthIntervals; i++ {
		chronological[i] = annualRows[growthIntervals-i]
	}
	var revGrowths []float64
	for i := 1; i < len(chronological); i++ {
		if chronological[i-1].Revenue > 0 && !math.IsNaN(chronological[i].Revenue) {
			revGrowths = append(revGrowths, (chronological[i].Revenue/chronological[i-1].Revenue)-1.0)
		}
	}
	revenueStability := NanStd(revGrowths)

	ttm := calculateTTM(raw.Quarterly)

	marketCap := math.NaN()
	if raw.Info.MarketCap != nil {
		marketCap = *raw.Info.MarketCap
	}
	enterpriseValue := math.NaN()
	if raw.Info.EnterpriseValue != nil {
		enterpriseValue = *raw.Info.EnterpriseValue
	}

	currentEBITDA := latest.EBITDA
	if !math.IsNaN(ttm.EBITDA) {
		currentEBITDA = ttm.EBITDA
	}

	currentEBIT := latest.EBIT
	if !math.IsNaN(ttm.EBIT) {
		currentEBIT = ttm.EBIT
	}

	currentFCF := latest.FCF
	if !math.IsNaN(ttm.FCF) {
		currentFCF = ttm.FCF
	}

	currentAdjFCF := latest.AdjFCF
	if !math.IsNaN(ttm.AdjFCF) {
		currentAdjFCF = ttm.AdjFCF
	}

	var valid3yFCF []float64
	var valid3yAdjFCF []float64
	for _, r := range avgRows {
		if !math.IsNaN(r.FCF) {
			valid3yFCF = append(valid3yFCF, r.FCF)
		}
		if !math.IsNaN(r.AdjFCF) {
			valid3yAdjFCF = append(valid3yAdjFCF, r.AdjFCF)
		}
	}

	fcf3ySum := math.NaN()
	if len(valid3yFCF) > 0 {
		var s float64
		for _, v := range valid3yFCF {
			s += v
		}
		fcf3ySum = s
	}

	adjFCF3ySum := math.NaN()
	if len(valid3yAdjFCF) > 0 {
		var s float64
		for _, v := range valid3yAdjFCF {
			s += v
		}
		adjFCF3ySum = s
	}

	extractSlice := func(fn func(r AnnualRow) float64) []float64 {
		s := make([]float64, len(avgRows))
		for i, r := range avgRows {
			s[i] = fn(r)
		}
		return s
	}

	currency := raw.Info.FinancialCurrency
	if currency == "" {
		currency = raw.Info.Currency
	}
	if currency == "" {
		currency = "USD"
	}

	res := &ScreenResult{
		Ticker:                raw.Info.Ticker,
		CompanyName:           raw.Info.ShortName,
		Currency:              currency,
		Sector:                raw.Info.Sector,
		Industry:              raw.Info.Industry,
		MarketCap:             marketCap,
		EnterpriseValue:       enterpriseValue,
		DepFactorLatest:       latest.DepFactor,
		AvgDepFactor3Y:        NanMean(extractSlice(func(r AnnualRow) float64 { return r.DepFactor })),
		DepFactorTTM:          SafePositiveDivide(ttm.EBITDA, ttm.EBIT),
		DAEBITRatio:           NanMean(extractSlice(func(r AnnualRow) float64 { return r.DAToEBIT })),
		RevCAGR3Y:             revCAGR,
		EBITCAGR3Y:            ebitCAGR,
		ROICLatest:            latest.ROIC,
		ROICAvg3Y:             NanMean(extractSlice(func(r AnnualRow) float64 { return r.ROIC })),
		ROICTrend3Y:           roicTrend,
		GrossMargin:           latest.GrossMargin,
		OperatingMargin:       latest.OperatingMargin,
		OperatingMargin3Y:     NanMean(extractSlice(func(r AnnualRow) float64 { return r.OperatingMargin })),
		EBITMarginTrend3Y:     ebitMarginTrend,
		FCFEBITDARatio:        NanMean(extractSlice(func(r AnnualRow) float64 { return r.FCFEBITDA })),
		FCFMargin:             NanMean(extractSlice(func(r AnnualRow) float64 { return r.FCFMargin })),
		AdjFCFEBITDARatio:     NanMean(extractSlice(func(r AnnualRow) float64 { return r.AdjFCFEBITDA })),
		AdjFCFNetIncomeLatest: latest.AdjFCFNetIncome,
		AdjFCFNetIncomeRatio:  NanMean(extractSlice(func(r AnnualRow) float64 { return r.AdjFCFNetIncome })),
		AdjFCFMargin:          NanMean(extractSlice(func(r AnnualRow) float64 { return r.AdjFCFMargin })),
		SBCRevLatest:          latest.SBCRevenue,
		SBCRevRatio:           NanMean(extractSlice(func(r AnnualRow) float64 { return r.SBCRevenue })),
		SBCFCFRatio:           NanMean(extractSlice(func(r AnnualRow) float64 { return r.SBCFCF })),
		CapExRevLatest:        latest.CapExRevenue,
		CapExRevRatio:         NanMean(extractSlice(func(r AnnualRow) float64 { return r.CapExRevenue })),
		NetDebtEBITDA:         SafePositiveDivide(latest.NetDebt, currentEBITDA),
		ShareCountCAGR3Y:      sharesCAGR,
		RevenueStability:      revenueStability,
		EVEBITDA:              SafePositiveDivide(enterpriseValue, currentEBITDA),
		EVEBIT:                SafePositiveDivide(enterpriseValue, currentEBIT),
		TTMFCFYield:           SafePositiveDivide(ttm.FCF, marketCap),
		TTMAdjFCFYield:        SafePositiveDivide(ttm.AdjFCF, marketCap),
		FCFYield:              SafePositiveDivide(currentFCF, marketCap),
		AdjFCFYield:           SafePositiveDivide(currentAdjFCF, marketCap),
		FCF3YEVYield:          SafePositiveDivide(fcf3ySum, enterpriseValue),
		AdjFCF3YEVYield:       SafePositiveDivide(adjFCF3ySum, enterpriseValue),
	}

	if len(avgRows) > 0 {
		res.DepFactorY1 = avgRows[0].DepFactor
		res.FiscalY1 = strings.Split(avgRows[0].Period, "T")[0]
	}
	if len(avgRows) > 1 {
		res.DepFactorY2 = avgRows[1].DepFactor
		res.FiscalY2 = strings.Split(avgRows[1].Period, "T")[0]
	}
	if len(avgRows) > 2 {
		res.DepFactorY3 = avgRows[2].DepFactor
		res.FiscalY3 = strings.Split(avgRows[2].Period, "T")[0]
	}

	res.IntensityLevel = IntensityLevel(res.AvgDepFactor3Y)

	scores := CalculateScores(res)
	res.ReturnsScore = scores.ReturnsScore
	res.GrowthScore = scores.GrowthScore
	res.CashScore = scores.CashScore
	res.ValuationScore = scores.ValuationScore
	res.Score = scores.Score

	return res, nil
}
