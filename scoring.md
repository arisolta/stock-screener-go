# Scoring Methodology

This screener uses a 100-point composite heuristic score to evaluate companies across 4 intuitive, 1-to-1 matching quadrants: **Returns**, **Growth**, **Cash Flow**, and **Valuation & Balance Sheet**.

---

## 4-Quadrant Scoring Architecture

The overall **Score** is a weighted average of four 0–100 normalized sub-scores mapping directly to the four quadrants of the stock card:

$$\text{Score} = (\text{Returns\_Score} \times 30\%) + (\text{Growth\_Score} \times 25\%) + (\text{Cash\_Score} \times 20\%) + (\text{Valuation\_Score} \times 25\%)$$

---

## Sub-Score Point Budgets & Thresholds

### 1. `Returns_Score` (30% Weight) — Quadrant 1: Profitability & Returns
Measures return on capital, margin levels, and asset-light depreciation economics.

| Metric | Max Pts | Better | Exceptional (100%) | Good (~85%) | Zero (0%) |
|:---|:---:|:---:|:---:|:---:|:---:|
| `ROIC_Latest` ($\frac{\text{NOPAT}}{\text{Invested Capital}}$ FY) | 15 | Higher | $\ge 50.0\%$ | $\ge 25.0\%$ | $\le 0.0\%$ |
| `ROIC_Trend_3Y` ($\Delta \text{ROIC}$ 3Y Trend) | 4 | Higher | $\ge +12.0\%$ | $\ge +5.0\%$ | $\le -5.0\%$ |
| `Avg_Dep_Factor_3Y` ($\frac{\text{EBITDA}}{\text{EBIT}}$ 3Y Avg) | 6 | Lower | $\le 1.00$ | $\le 1.05$ | $\ge 2.00$ |
| `Operating_Margin` ($\frac{\text{EBIT}}{\text{Revenue}}$ FY) | 3 | Higher | $\ge 50.0\%$ | $\ge 30.0\%$ | $\le 5.0\%$ |
| `Gross_Margin` ($\frac{\text{Gross Profit}}{\text{Revenue}}$ FY) | 2 | Higher | $\ge 80.0\%$ | $\ge 60.0\%$ | $\le 20.0\%$ |

---

### 2. `Growth_Score` (25% Weight) — Quadrant 2: Growth & Stability
Measures multi-year compounding trajectory, earnings acceleration, and top-line consistency.

| Metric | Max Pts | Better | Exceptional (100%) | Good (~85%) | Zero (0%) |
|:---|:---:|:---:|:---:|:---:|:---:|
| `Rev_CAGR_3Y` (3Y CAGR) | 8 | Higher | $\ge 25.0\%$ | $\ge 12.0\%$ | $\le -5.0\%$ |
| `EBIT_CAGR_3Y` (3Y CAGR) | 11 | Higher | $\ge 35.0\%$ | $\ge 15.0\%$ | $\le -10.0\%$ |
| `EBIT_Margin_Trend_3Y` (3Y Trend) | 3 | Higher | $\ge +12.0\%$ | $\ge +5.0\%$ | $\le -5.0\%$ |
| `Revenue_Stability` ($\sigma$ of YoY growth) | 3 | Lower | $\le 0.0\%$ | $\le 3.0\%$ | $\ge 20.0\%$ |

---

### 3. `Cash_Score` (20% Weight) — Quadrant 3: Cash Flow Conversion & Reinvestment
Measures how reliably accounting profits convert to real owner cash flow, CapEx discipline, and dilution overhead.

| Metric | Max Pts | Better | Exceptional (100%) | Good (~85%) | Zero (0%) |
|:---|:---:|:---:|:---:|:---:|:---:|
| `Adj_FCF_Net_Income_Ratio` (3Y Avg) | 8 | Higher | $\ge 150.0\%$ | $\ge 100.0\%$ | $\le 0.0\%$ |
| `Adj_FCF_Margin` (3Y Avg) | 6 | Higher | $\ge 35.0\%$ | $\ge 20.0\%$ | $\le 0.0\%$ |
| `CapEx_Rev_Ratio` ($\frac{\text{CapEx}}{\text{Revenue}}$ 3Y Avg) | 4 | Lower | $\le 0.0\%$ | $\le 2.0\%$ | $\ge 15.0\%$ |
| `SBC_Rev_Ratio` (3Y Avg) | 2 | Lower | $\le 0.0\%$ | $\le 0.0\%$ | $\ge 8.0\%$ |

---

### 4. `Valuation_Score` (25% Weight) — Quadrant 4: Valuation & Capital Structure
Measures the entry price per unit of owner earnings and balance sheet safety.

| Metric | Max Pts | Better | Exceptional (100%) | Good (~85%) | Zero (0%) |
|:---|:---:|:---:|:---:|:---:|:---:|
| `Adj_FCF_Yield` ($\frac{\text{Adj FCF}}{\text{Market Cap}}$) | 9 | Higher | $\ge 12.0\%$ | $\ge 6.0\%$ | $\le 0.0\%$ |
| `EV_EBIT` ($\frac{\text{Enterprise Value}}{\text{EBIT}}$) | 8 | Lower | $\le 5.0\text{x}$ | $\le 10.0\text{x}$ | $\ge 40.0\text{x}$ |
| `Net_Debt_EBITDA` ($\frac{\text{Net Debt}}{\text{EBITDA}}$) | 5 | Lower | $\le -1.0\text{x}$ | $\le 0.0\text{x}$ | $\ge 4.0\text{x}$ |
| `Share_Count_CAGR_3Y` (3Y CAGR) | 3 | Lower | $\le -6.0\%$ | $\le -3.0\%$ | $\ge +5.0\%$ |

---

## 3-Year Data Window & Mathematical Types

1. **3-Year Arithmetic Averages (`3Y Avg`)**: Unweighted mean of annual figures across the latest 3 fiscal years ($Y_1, Y_2, Y_3$).
2. **3-Year Cumulative Sums (`3Y Sum`)**: Total 3-year cumulative cash generated divided by Enterprise Value ($\frac{\sum \text{Adj FCF}}{\text{EV}}$).
3. **3-Year Compounded Rates (`3Y CAGR`)**: Compounded annualized rate across 4 annual observations ($(\frac{Y_3}{Y_0})^{1/3} - 1$).
4. **3-Year Historical Trends (`3Y Trend`)**: Percentage-point delta from oldest to latest fiscal year ($Y_3 - Y_0$).

---

## Score Interpretation Bands

| Score | Rating | Meaning |
|:---:|:---:|---|
| **80+** | 🟢 **Strong / Elite** | Exceptional economics and capital efficiency. Compelling candidate for deep analysis. |
| **70 – 80** | 🟢 **Good Quality** | Attractive quality profile with healthy returns; verify individual watch-outs. |
| **60 – 70** | 🟡 **Moderate / Mixed** | Mixed signals; capital intensity or valuation may be moderating returns. |
| **50 – 60** | 🔴 **Average** | Average economics; requires a special situation or deep discount. |
| **< 50** | 🔴 **Weak** | Fails key quality, reinvestment, or cash conversion thresholds. |
