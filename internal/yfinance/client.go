package yfinance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
)

// Client handles authenticated requests to Yahoo Finance.
type Client struct {
	httpClient *http.Client
	crumb      string
	mu         sync.RWMutex
}

// NewClient initializes a new Yahoo Finance API client with cookie handling.
func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		httpClient: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
	}
}

// ensureAuth obtains a session cookie and crumb if not already present.
func (c *Client) ensureAuth() error {
	c.mu.RLock()
	if c.crumb != "" {
		c.mu.RUnlock()
		return nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.crumb != "" {
		return nil
	}

	// 1. Visit fc.yahoo.com to set base cookies
	req, err := http.NewRequest("GET", "https://fc.yahoo.com", nil)
	if err != nil {
		return fmt.Errorf("failed to create fc.yahoo.com request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := c.httpClient.Do(req)
	if err == nil {
		resp.Body.Close()
	}

	// 2. Fetch the crumb
	crumbReq, err := http.NewRequest("GET", "https://query2.finance.yahoo.com/v1/test/getcrumb", nil)
	if err != nil {
		return fmt.Errorf("failed to create getcrumb request: %w", err)
	}
	crumbReq.Header.Set("User-Agent", userAgent)
	crumbReq.Header.Set("Accept", "*/*")
	crumbReq.Header.Set("Accept-Language", "en-US,en;q=0.9")

	crumbResp, err := c.httpClient.Do(crumbReq)
	if err != nil {
		return fmt.Errorf("failed to execute getcrumb request: %w", err)
	}
	defer crumbResp.Body.Close()

	if crumbResp.StatusCode != http.StatusOK {
		// Fallback to query1 if query2 fails
		q1Req, _ := http.NewRequest("GET", "https://query1.finance.yahoo.com/v1/test/getcrumb", nil)
		q1Req.Header.Set("User-Agent", userAgent)
		q1Resp, errQ1 := c.httpClient.Do(q1Req)
		if errQ1 == nil && q1Resp.StatusCode == http.StatusOK {
			defer q1Resp.Body.Close()
			b, _ := io.ReadAll(q1Resp.Body)
			c.crumb = strings.TrimSpace(string(b))
			return nil
		}
		if errQ1 == nil {
			q1Resp.Body.Close()
		}
		return fmt.Errorf("getcrumb returned HTTP %d", crumbResp.StatusCode)
	}

	body, err := io.ReadAll(crumbResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read crumb response: %w", err)
	}

	crumbStr := strings.TrimSpace(string(body))
	if crumbStr == "" || strings.Contains(crumbStr, "<html") || strings.Contains(crumbStr, "Too Many Requests") {
		return fmt.Errorf("invalid crumb received: %s", crumbStr)
	}

	c.crumb = crumbStr
	return nil
}

// GetFinancials retrieves company metadata and fundamental timeseries for a given ticker.
func (c *Client) GetFinancials(ticker string) (*RawFinancialData, error) {
	if err := c.ensureAuth(); err != nil {
		return nil, fmt.Errorf("auth error for %s: %w", ticker, err)
	}

	c.mu.RLock()
	crumb := c.crumb
	c.mu.RUnlock()

	// 1. Fetch QuoteSummary (metadata, marketCap, enterpriseValue, etc.)
	info, err := c.fetchQuoteSummary(ticker, crumb)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quote summary for %s: %w", ticker, err)
	}

	// 2. Fetch Fundamentals Timeseries (income, balance sheet, cash flow)
	annual, quarterly, err := c.fetchTimeseries(ticker, crumb)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch timeseries for %s: %w", ticker, err)
	}

	return &RawFinancialData{
		Info:      *info,
		Annual:    annual,
		Quarterly: quarterly,
	}, nil
}

func (c *Client) fetchQuoteSummary(ticker string, crumb string) (*CompanyInfo, error) {
	url := fmt.Sprintf("https://query2.finance.yahoo.com/v10/finance/quoteSummary/%s?crumb=%s&modules=financialData,defaultKeyStatistics,assetProfile,summaryDetail,price", ticker, crumb)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "*/*")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("quoteSummary HTTP %d", resp.StatusCode)
	}

	var parsed QuoteSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if len(parsed.QuoteSummary.Result) == 0 {
		return nil, fmt.Errorf("no quote summary returned for %s", ticker)
	}

	res := parsed.QuoteSummary.Result[0]
	info := &CompanyInfo{
		Ticker:            ticker,
		ShortName:         res.Price.ShortName,
		LongName:          res.Price.LongName,
		Currency:          res.Price.Currency,
		FinancialCurrency: res.FinancialData.FinancialCurrency,
		Sector:            res.AssetProfile.Sector,
		Industry:          res.AssetProfile.Industry,
	}

	if res.SummaryDetail.MarketCap != nil {
		info.MarketCap = &res.SummaryDetail.MarketCap.Raw
	}
	if res.DefaultKeyStatistics.EnterpriseValue != nil {
		info.EnterpriseValue = &res.DefaultKeyStatistics.EnterpriseValue.Raw
	}
	if res.FinancialData.EffectiveTaxRate != nil {
		info.EffectiveTaxRate = &res.FinancialData.EffectiveTaxRate.Raw
	}

	if info.ShortName == "" && info.LongName != "" {
		info.ShortName = info.LongName
	} else if info.ShortName == "" {
		info.ShortName = ticker
	}

	return info, nil
}

func (c *Client) fetchTimeseries(ticker string, crumb string) (map[string]map[string]float64, map[string]map[string]float64, error) {
	types := []string{
		// Annual statements
		"annualTotalRevenue", "annualOperatingRevenue", "annualGrossProfit",
		"annualOperatingIncome", "annualEBIT", "annualEBITDA",
		"annualReconciledDepreciation", "annualDepreciationAndAmortization", "annualDepreciation",
		"annualNetIncome", "annualNetIncomeCommonStockholders",
		"annualOperatingCashFlow", "annualTotalCashFromOperatingActivities",
		"annualCapitalExpenditure", "annualCapitalExpenditures",
		"annualStockBasedCompensation", "annualStockBasedCompensationAndOther",
		"annualShareBasedCompensation", "annualShareBasedCompensationExpense",
		"annualDilutedAverageShares", "annualBasicAverageShares", "annualAverageDilutionEarnings",
		"annualStockholdersEquity", "annualTotalStockholderEquity",
		"annualTotalDebt", "annualLongTermDebtAndCapitalLeaseObligation", "annualLongTermDebt",
		"annualCashAndCashEquivalents", "annualCashCashEquivalentsAndShortTermInvestments",
		// Quarterly statements
		"quarterlyTotalRevenue", "quarterlyOperatingRevenue",
		"quarterlyOperatingIncome", "quarterlyEBIT", "quarterlyEBITDA",
		"quarterlyReconciledDepreciation", "quarterlyDepreciationAndAmortization",
		"quarterlyOperatingCashFlow", "quarterlyTotalCashFromOperatingActivities",
		"quarterlyCapitalExpenditure", "quarterlyCapitalExpenditures",
		"quarterlyStockBasedCompensation", "quarterlyStockBasedCompensationAndOther",
	}

	url := fmt.Sprintf(
		"https://query2.finance.yahoo.com/ws/fundamentals-timeseries/v1/finance/timeseries/%s?symbol=%s&type=%s&period1=1483142400&period2=1893456000&crumb=%s",
		ticker, ticker, strings.Join(types, ","), crumb,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "*/*")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("fundamentals-timeseries HTTP %d", resp.StatusCode)
	}

	var parsed TimeseriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, nil, fmt.Errorf("decode timeseries error: %w", err)
	}

	annual := make(map[string]map[string]float64)
	quarterly := make(map[string]map[string]float64)

	for _, item := range parsed.Timeseries.Result {
		meta, ok := item["meta"].(map[string]interface{})
		if !ok {
			continue
		}
		typeArr, ok := meta["type"].([]interface{})
		if !ok || len(typeArr) == 0 {
			continue
		}
		metricName, _ := typeArr[0].(string)
		if metricName == "" {
			continue
		}

		entries, ok := item[metricName].([]interface{})
		if !ok {
			continue
		}

		targetMap := annual
		if strings.HasPrefix(metricName, "quarterly") {
			targetMap = quarterly
		}

		if _, exists := targetMap[metricName]; !exists {
			targetMap[metricName] = make(map[string]float64)
		}

		for _, entryRaw := range entries {
			entry, ok := entryRaw.(map[string]interface{})
			if !ok {
				continue
			}
			asOfDate, _ := entry["asOfDate"].(string)
			if asOfDate == "" {
				continue
			}
			repVal, ok := entry["reportedValue"].(map[string]interface{})
			if !ok {
				continue
			}
			rawNum, ok := repVal["raw"].(float64)
			if !ok {
				continue
			}
			targetMap[metricName][asOfDate] = rawNum
		}
	}

	return annual, quarterly, nil
}
