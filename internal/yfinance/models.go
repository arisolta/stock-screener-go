package yfinance

// RawValue represents a raw numerical value from Yahoo Finance JSON.
type RawValue struct {
	Raw float64 `json:"raw"`
	Fmt string  `json:"fmt,omitempty"`
}

// QuoteSummaryResponse represents the top-level response from Yahoo Finance quoteSummary.
type QuoteSummaryResponse struct {
	QuoteSummary struct {
		Result []QuoteSummaryResult `json:"result"`
		Error  interface{}          `json:"error"`
	} `json:"quoteSummary"`
}

type QuoteSummaryResult struct {
	Price struct {
		ShortName string `json:"shortName"`
		LongName  string `json:"longName"`
		Currency  string `json:"currency"`
	} `json:"price"`
	FinancialData struct {
		FinancialCurrency string    `json:"financialCurrency"`
		EffectiveTaxRate  *RawValue `json:"effectiveTaxRate"`
	} `json:"financialData"`
	DefaultKeyStatistics struct {
		EnterpriseValue *RawValue `json:"enterpriseValue"`
	} `json:"defaultKeyStatistics"`
	SummaryDetail struct {
		MarketCap *RawValue `json:"marketCap"`
	} `json:"summaryDetail"`
	AssetProfile struct {
		Sector   string `json:"sector"`
		Industry string `json:"industry"`
	} `json:"assetProfile"`
}

// TimeseriesResponse represents the response from Yahoo Finance fundamentals-timeseries endpoint.
type TimeseriesResponse struct {
	Timeseries struct {
		Result []map[string]interface{} `json:"result"`
		Error  interface{}              `json:"error"`
	} `json:"timeseries"`
}

// CompanyInfo contains aggregated high-level metadata.
type CompanyInfo struct {
	Ticker            string   `json:"ticker"`
	ShortName         string   `json:"short_name"`
	LongName          string   `json:"long_name"`
	Currency          string   `json:"currency"`
	FinancialCurrency string   `json:"financial_currency"`
	MarketCap         *float64 `json:"market_cap"`
	EnterpriseValue   *float64 `json:"enterprise_value"`
	EffectiveTaxRate  *float64 `json:"effective_tax_rate"`
	Sector            string   `json:"sector"`
	Industry          string   `json:"industry"`
}

// RawFinancialData is the complete extracted payload for a ticker.
type RawFinancialData struct {
	Info      CompanyInfo                   `json:"info"`
	Annual    map[string]map[string]float64 `json:"annual"`    // metricName -> period (YYYY-MM-DD) -> value
	Quarterly map[string]map[string]float64 `json:"quarterly"` // metricName -> period (YYYY-MM-DD) -> value
}
