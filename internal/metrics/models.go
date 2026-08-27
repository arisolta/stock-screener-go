package metrics

// AnnualRow holds parsed line items and derived metrics for a single fiscal period.
type AnnualRow struct {
	Period          string  `json:"period"`
	Revenue         float64 `json:"revenue"`
	GrossProfit     float64 `json:"gross_profit"`
	EBIT            float64 `json:"ebit"`
	NetIncome       float64 `json:"net_income"`
	EBITDA          float64 `json:"ebitda"`
	Depreciation    float64 `json:"depreciation"`
	FCF             float64 `json:"fcf"`
	SBC             float64 `json:"sbc"`
	AdjFCF          float64 `json:"adj_fcf"`
	CapEx           float64 `json:"capex"`
	AverageShares   float64 `json:"average_shares"`
	GrossMargin     float64 `json:"gross_margin"`
	OperatingMargin float64 `json:"operating_margin"`
	DepFactor       float64 `json:"dep_factor"`
	DAToEBIT        float64 `json:"da_to_ebit"`
	FCFEBITDA       float64 `json:"fcf_ebitda"`
	FCFMargin       float64 `json:"fcf_margin"`
	AdjFCFEBITDA    float64 `json:"adj_fcf_ebitda"`
	AdjFCFNetIncome float64 `json:"adj_fcf_net_income"`
	AdjFCFMargin    float64 `json:"adj_fcf_margin"`
	SBCRevenue      float64 `json:"sbc_revenue"`
	SBCFCF          float64 `json:"sbc_fcf"`
	CapExRevenue    float64 `json:"capex_revenue"`
	NetDebt         float64 `json:"net_debt"`
	ROIC            float64 `json:"roic"`
}

// Scores represents the 4-category calculated scoring breakdown.
type Scores struct {
	ReturnsScore   float64 `json:"returns_score"`
	GrowthScore    float64 `json:"growth_score"`
	CashScore      float64 `json:"cash_score"`
	ValuationScore float64 `json:"valuation_score"`
	Score          float64 `json:"score"`
}

// ScreenResult represents the full output metrics for a single stock.
type ScreenResult struct {
	Ticker                   string  `json:"Ticker"`
	CompanyName              string  `json:"Company Name"`
	Currency                 string  `json:"Currency"`
	Sector                   string  `json:"Sector"`
	Industry                 string  `json:"Industry"`
	MarketCap                float64 `json:"Market_Cap"`
	EnterpriseValue          float64 `json:"Enterprise_Value"`
	IntensityLevel           string  `json:"Intensity_Level"`
	
	DepFactorLatest          float64 `json:"Dep_Factor_Latest"`
	AvgDepFactor3Y           float64 `json:"Avg_Dep_Factor_3Y"`
	DepFactorTTM             float64 `json:"Dep_Factor_TTM"`
	DAEBITRatio              float64 `json:"D&A_EBIT_Ratio"`
	
	RevCAGR3Y                float64 `json:"Rev_CAGR_3Y"`
	EBITCAGR3Y               float64 `json:"EBIT_CAGR_3Y"`
	ROICLatest               float64 `json:"ROIC_Latest"`
	ROICAvg3Y                float64 `json:"ROIC_Avg_3Y"`
	ROICTrend3Y              float64 `json:"ROIC_Trend_3Y"`
	
	GrossMargin              float64 `json:"Gross_Margin"`
	OperatingMargin          float64 `json:"Operating_Margin"`
	OperatingMargin3Y        float64 `json:"Operating_Margin_3Y"`
	EBITMarginTrend3Y        float64 `json:"EBIT_Margin_Trend_3Y"`
	
	FCFEBITDARatio           float64 `json:"FCF_EBITDA_Ratio"`
	FCFMargin                float64 `json:"FCF_Margin"`
	AdjFCFEBITDARatio        float64 `json:"Adj_FCF_EBITDA_Ratio"`
	AdjFCFNetIncomeLatest    float64 `json:"Adj_FCF_Net_Income_Latest"`
	AdjFCFNetIncomeRatio     float64 `json:"Adj_FCF_Net_Income_Ratio"`
	AdjFCFMargin             float64 `json:"Adj_FCF_Margin"`
	
	SBCRevLatest             float64 `json:"SBC_Rev_Latest"`
	SBCRevRatio              float64 `json:"SBC_Rev_Ratio"`
	SBCFCFRatio              float64 `json:"SBC_FCF_Ratio"`
	CapExRevLatest           float64 `json:"CapEx_Rev_Latest"`
	CapExRevRatio            float64 `json:"CapEx_Rev_Ratio"`
	
	NetDebtEBITDA            float64 `json:"Net_Debt_EBITDA"`
	ShareCountCAGR3Y         float64 `json:"Share_Count_CAGR_3Y"`
	RevenueStability         float64 `json:"Revenue_Stability"`
	
	EVEBITDA                 float64 `json:"EV_EBITDA"`
	EVEBIT                   float64 `json:"EV_EBIT"`
	TTMFCFYield              float64 `json:"TTM_FCF_Yield"`
	TTMAdjFCFYield           float64 `json:"TTM_Adj_FCF_Yield"`
	FCFYield                 float64 `json:"FCF_Yield"`
	AdjFCFYield              float64 `json:"Adj_FCF_Yield"`
	FCF3YEVYield             float64 `json:"FCF_3Y_EV_Yield"`
	AdjFCF3YEVYield          float64 `json:"Adj_FCF_3Y_EV_Yield"`
	
	DepFactorY1              float64 `json:"Dep_Factor_Y1"`
	DepFactorY2              float64 `json:"Dep_Factor_Y2"`
	DepFactorY3              float64 `json:"Dep_Factor_Y3"`
	FiscalY1                 string  `json:"Fiscal_Y1"`
	FiscalY2                 string  `json:"Fiscal_Y2"`
	FiscalY3                 string  `json:"Fiscal_Y3"`
	
	ReturnsScore             float64 `json:"Returns_Score"`
	GrowthScore              float64 `json:"Growth_Score"`
	CashScore                float64 `json:"Cash_Score"`
	ValuationScore           float64 `json:"Valuation_Score"`
	Score                    float64 `json:"Score"`
}
