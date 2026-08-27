# Capital Intensity & Quality Stock Screener (Go)

A fast, concurrent CLI stock screener written in pure Go inspired by Michael Mauboussin's capital intensity and return on invested capital (ROIC) framework.

It ranks public companies using depreciation factors, true owner cash conversion (SBC-adjusted FCF), returns on invested capital, growth stability, balance sheet leverage, and valuation multiples pulled directly from financial statements.

---

## Key Features

- **Blazing Fast Concurrency**: Multi-worker goroutine pool to screen hundreds of tickers in seconds with zero external runtime dependencies.
- **Harmonious 4-Quadrant Design**: 4 intuitive sub-scores that map 1-to-1 directly onto the 4 visual quadrants of each company card.
- **Mauboussin Capital Intensity Analysis**: Separates asset-light compounders (`EBITDA/EBIT < 1.3`) from high-reinvestment businesses (`EBITDA/EBIT > 1.7`).
- **Owner Earnings Focus**: Treats Stock-Based Compensation (SBC) as a true economic cost to calculate adjusted FCF margins and conversion.
- **Beautiful Terminal UI**: Pixel-perfect 88-column stock cards with visual sub-score progress bars (`[██████░░]`), trend glyphs (`▲`/`▼`), and executive strengths/watch-outs analysis.
- **Flexible Exporting**: Export all calculated metrics to `.csv`, styled `.xlsx` (Excel), or `.md` (Markdown).
- **Disk Caching**: Built-in 24-hour TTL caching in `.cache/capital_screener/` to prevent redundant API calls.

---

## Build & Install

Ensure you have [Go 1.21+](https://go.dev/) installed:

```bash
# Build the local binary
go build -o stock-screener cmd/stock-screener/main.go

# Or install it globally to your $GOPATH/bin
go install ./cmd/stock-screener
```

---

## Usage & Display Modes

### 1. Detailed Stock Cards (Default View)
Displays structured 2-column cards containing all key financial metrics, visual sub-scores, and bulleted takeaways organized into 4 distinct quadrants:

```bash
./stock-screener --tickers AAPL MSFT GOOGL
```

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│  AAPL (Apple Inc.) • Technology • MktCap: $3.45T • Cur: USD                            │
│  Overall Score: 78.4/100 [████████] Strong Quality • Intensity: Low (3Y Avg: 1.09)     │
├────────────────────────────────────────────────────────────────────────────────────────┤
│  SUB-SCORES:  Returns (30%)   [███████░] 92      Growth (25%)    [████░░░░] 48         │
│               Cash (20%)      [██████░░] 79      Valuation (25%) [████░░░░] 49         │
├────────────────────────────────────────────┬───────────────────────────────────────────┤
│ 1. PROFITABILITY & RETURNS (30%)           │ 2. GROWTH & STABILITY (25%)               │
│   • EBITDA/EBIT (Dep) FY:           1.09   │   • Gross Margin FY:               46.9%  │
│   • EBITDA/EBIT 3Y Avg:             1.09   │   • Operating Margin FY:           32.0%  │
│   • EBITDA/EBIT TTM:                1.08   │   • Operating Margin 3Y Avg:       31.1%  │
│   • 3Y Dep Trajectory: 1.10 > 1.09 > 1.09  │   • EBIT Margin 3Y Trend:        ▲ +1.7%  │
│   • ROIC Latest FY:                77.0%   │   • Revenue 3Y CAGR:                1.8%  │
│   • ROIC 3Y Avg:                   71.0%   │   • EBIT 3Y CAGR:                   3.7%  │
│   • ROIC 3Y Trend:              ▲ +17.9%   │   • Shares 3Y CAGR:                -2.8%  │
│                                            │   • Revenue Stability:     3.8% (StdDev)  │
├────────────────────────────────────────────┼───────────────────────────────────────────┤
│ 3. CASH FLOW CONVERSION (20%)              │ 4. VALUATION & CAPITAL STRUCTURE (25%)    │
│   • Adj FCF / Net Income FY:       76.7%   │   • EV / EBIT Multiple:           29.35x  │
│   • Adj FCF / NI 3Y Avg:           90.6%   │   • EV / EBITDA Multiple:         27.06x  │
│   • Adj FCF Margin 3Y Avg:         22.9%   │   • TTM Adj FCF Yield:              2.7%  │
│   • CapEx / Revenue FY:             3.1%   │   • 3Y Sum AdjFCF / EV:             6.0%  │
│   • CapEx / Revenue 3Y Avg:         2.8%   │   • Net Debt / EBITDA:             0.37x  │
│   • SBC / Revenue FY:               3.1%   │                                           │
│   • SBC / Revenue 3Y Avg:           3.0%   │                                           │
│   • SBC / FCF 3Y Avg:              11.5%   │                                           │
├────────────────────────────────────────────────────────────────────────────────────────┤
│  ★ Strengths: Elite ROIC (77.0%) • High Cash Conversion (90.6% 3Y Avg) • Asset-Light Capital (1.09 3Y Avg) │
│  ⚠ Watch-Outs: Premium Multiple (29.3x EV/EBIT) • Sluggish Top-Line Growth (1.8% 3Y CAGR) │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

### 2. Compact Comparison Table (`--table`)
For screening multiple stocks side-by-side on a single terminal screen without wrapping:
```bash
./stock-screener --table
./stock-screener --tickers AAPL MSFT GOOGL PYPL --table
```

### 3. Thematic Grouped Tables (`--grouped`)
Partitions all metrics across 4 focused tables:
1. Composite Scores & Capital Intensity
2. Returns, Margins & Growth
3. Cash Flow Conversion & Reinvestment
4. Valuation & Balance Sheet Risk
```bash
./stock-screener --grouped
```

---

## Filtering, Sorting & Universe Files

```bash
# Screen from a universe file with 10 concurrent workers
./stock-screener --file sp500.txt --concurrency 10

# Filter companies with ROIC >= 15% and Depreciation Factor <= 1.4
./stock-screener --min-roic 15 --max-dep-factor 1.4

# Rank output by a specific metric (Score, ROIC_Latest, Returns_Score, Valuation_Score, etc.)
./stock-screener --rank-by Returns_Score

# Bypass local cache
./stock-screener --no-cache
```

---

## Exporting Results

Pass `--export <filename>` to save all calculated metrics to disk:

```bash
./stock-screener --file sp500.txt --export results.xlsx
./stock-screener --tickers AAPL MSFT GOOGL --export screener.csv
./stock-screener --tickers AAPL MSFT GOOGL --export report.md
```

---

## 4-Quadrant Scoring Formula

The composite **Score** (0–100) maps 1-to-1 directly onto the 4 card quadrants:

$$\text{Score} = (\text{Returns Score} \times 30\%) + (\text{Growth Score} \times 25\%) + (\text{Cash Score} \times 20\%) + (\text{Valuation Score} \times 25\%)$$

See [scoring.md](scoring.md) for full mathematical definitions and piecewise ramp functions.
