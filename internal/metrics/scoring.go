package metrics

import (
	"math"
)

var ScoreWeights = map[string]float64{
	"Returns_Score":   0.30,
	"Growth_Score":    0.25,
	"Cash_Score":      0.20,
	"Valuation_Score": 0.25,
}

func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func ScoreHigher(value, maxPoints, fullAt float64, zeroAt float64, exceptionalAt *float64) float64 {
	if math.IsNaN(value) {
		return 0.0
	}
	expAt := fullAt + math.Abs(fullAt-zeroAt)
	if exceptionalAt != nil {
		expAt = *exceptionalAt
	}
	if fullAt == zeroAt || expAt <= fullAt {
		if value >= fullAt {
			return maxPoints
		}
		return 0.0
	}
	if value <= zeroAt {
		return 0.0
	}
	goodPoints := maxPoints * 0.85
	if value <= fullAt {
		scaled := (value - zeroAt) / (fullAt - zeroAt)
		return clamp(scaled, 0.0, 1.0) * goodPoints
	}
	if value >= expAt {
		return maxPoints
	}
	scaled := (value - fullAt) / (expAt - fullAt)
	return goodPoints + clamp(scaled, 0.0, 1.0)*(maxPoints-goodPoints)
}

func ScoreLower(value, maxPoints, bestAt float64, zeroAt float64, exceptionalAt *float64) float64 {
	if math.IsNaN(value) {
		return 0.0
	}
	expAt := bestAt
	if exceptionalAt != nil {
		expAt = *exceptionalAt
	}
	if zeroAt == bestAt || expAt >= bestAt {
		if value <= bestAt {
			return maxPoints
		}
		if value >= zeroAt {
			return 0.0
		}
		scaled := (zeroAt - value) / (zeroAt - bestAt)
		return clamp(scaled, 0.0, 1.0) * maxPoints
	}
	if value >= zeroAt {
		return 0.0
	}
	goodPoints := maxPoints * 0.85
	if value >= bestAt {
		scaled := (zeroAt - value) / (zeroAt - bestAt)
		return clamp(scaled, 0.0, 1.0) * goodPoints
	}
	if value <= expAt {
		return maxPoints
	}
	scaled := (bestAt - value) / (bestAt - expAt)
	return goodPoints + clamp(scaled, 0.0, 1.0)*(maxPoints-goodPoints)
}

func ScoreRange(value, maxPoints, lowFull, highZero float64) float64 {
	exp := 0.0
	return ScoreLower(value, maxPoints, lowFull, highZero, &exp)
}

func NormalizeBucketScore(rawScore, maxPoints float64) float64 {
	if maxPoints <= 0 {
		return 0.0
	}
	return math.Round((rawScore/maxPoints)*1000) / 10
}

func ptr(v float64) *float64 {
	return &v
}

func CalculateScores(r *ScreenResult) Scores {
	// 1. Returns & Capital Intensity (30 pts max)
	returnsRaw := ScoreHigher(r.ROICLatest, 15, 0.25, 0.0, ptr(0.50)) +
		ScoreHigher(r.ROICTrend3Y, 4, 0.05, -0.05, ptr(0.12)) +
		ScoreLower(r.AvgDepFactor3Y, 6, 1.05, 2.0, ptr(1.0)) +
		ScoreHigher(r.OperatingMargin, 3, 0.30, 0.05, ptr(0.50)) +
		ScoreHigher(r.GrossMargin, 2, 0.60, 0.20, ptr(0.80))

	// 2. Growth & Stability (25 pts max)
	growthRaw := ScoreHigher(r.RevCAGR3Y, 8, 0.12, -0.05, ptr(0.25)) +
		ScoreHigher(r.EBITCAGR3Y, 11, 0.15, -0.10, ptr(0.35)) +
		ScoreHigher(r.EBITMarginTrend3Y, 3, 0.05, -0.05, ptr(0.12)) +
		ScoreRange(r.RevenueStability, 3, 0.03, 0.20)

	// 3. Cash Flow Conversion & Reinvestment (20 pts max)
	cashRaw := ScoreHigher(r.AdjFCFNetIncomeRatio, 8, 1.0, 0.0, ptr(1.5)) +
		ScoreHigher(r.AdjFCFMargin, 6, 0.20, 0.0, ptr(0.35)) +
		ScoreLower(r.CapExRevRatio, 4, 0.02, 0.15, ptr(0.0)) +
		ScoreLower(r.SBCRevRatio, 2, 0.0, 0.08, ptr(0.0))

	// 4. Valuation & Balance Sheet (25 pts max)
	valuationRaw := ScoreHigher(r.AdjFCFYield, 9, 0.06, 0.0, ptr(0.12)) +
		ScoreLower(r.EVEBIT, 8, 10.0, 40.0, ptr(5.0)) +
		ScoreLower(r.NetDebtEBITDA, 5, 0.0, 4.0, ptr(-1.0)) +
		ScoreLower(r.ShareCountCAGR3Y, 3, -0.03, 0.05, ptr(-0.06))

	s := Scores{
		ReturnsScore:   NormalizeBucketScore(returnsRaw, 30),
		GrowthScore:    NormalizeBucketScore(growthRaw, 25),
		CashScore:      NormalizeBucketScore(cashRaw, 20),
		ValuationScore: NormalizeBucketScore(valuationRaw, 25),
	}
	s.Score = math.Round((s.ReturnsScore*ScoreWeights["Returns_Score"]+
		s.GrowthScore*ScoreWeights["Growth_Score"]+
		s.CashScore*ScoreWeights["Cash_Score"]+
		s.ValuationScore*ScoreWeights["Valuation_Score"])*10) / 10

	return s
}
