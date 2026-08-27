package config

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var DefaultTickers = []string{"AAPL", "MSFT", "GOOGL", "UL", "PYPL", "GIS", "CLX", "HENKY"}

type Config struct {
	Tickers      []string
	FilePath     string
	RankBy       string
	MinROIC      *float64
	MaxDepFactor *float64
	ExportPath   string
	NoCache      bool
	Plain        bool
	Years        int
	Concurrency  int
	View         string // "cards" (default), "table", "grouped"
}

func ParseArgs() (*Config, error) {
	var filePath string
	var rankBy string
	var minROIC float64
	var maxDepFactor float64
	var exportPath string
	var noCache bool
	var plain bool
	var years int
	var concurrency int
	var view string
	var detail bool
	var table bool
	var grouped bool

	fs := flag.NewFlagSet("stock-screener", flag.ContinueOnError)
	fs.StringVar(&filePath, "file", "", "Text file with one ticker per line")
	fs.StringVar(&rankBy, "rank-by", "Score", "Column to rank by (Score, Avg_Dep_Factor_3Y, ROIC_Latest, etc.)")
	fs.Float64Var(&minROIC, "min-roic", -999999, "Filter out companies below this ROIC percent")
	fs.Float64Var(&maxDepFactor, "max-dep-factor", -999999, "Filter out companies above this average depreciation factor")
	fs.StringVar(&exportPath, "export", "", "Export path (.csv, .xlsx, .xls, .md, .markdown)")
	fs.BoolVar(&noCache, "no-cache", false, "Disable local JSON cache for yfinance responses")
	fs.BoolVar(&plain, "plain", false, "Disable colored terminal output")
	fs.IntVar(&years, "years", 3, "Number of historical years to analyze")
	fs.IntVar(&concurrency, "concurrency", 5, "Number of concurrent ticker requests")
	fs.StringVar(&view, "view", "", "Display view: cards (default), table, grouped")
	fs.BoolVar(&detail, "detail", false, "Display detailed cards (default)")
	fs.BoolVar(&detail, "cards", false, "Alias for --detail")
	fs.BoolVar(&table, "table", false, "Display compact summary table")
	fs.BoolVar(&grouped, "grouped", false, "Display themed grouped tables")

	args := os.Args[1:]
	var normalizedArgs []string
	var tickers []string

	boolFlags := map[string]bool{
		"no-cache": true,
		"plain":    true,
		"detail":   true,
		"cards":    true,
		"table":    true,
		"grouped":  true,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--tickers" || arg == "-tickers" {
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				tickers = append(tickers, args[i])
			}
		} else if strings.HasPrefix(arg, "--") {
			flagName := arg[2:]
			if strings.Contains(flagName, "=") {
				normalizedArgs = append(normalizedArgs, "-"+flagName)
			} else {
				normalizedArgs = append(normalizedArgs, "-"+flagName)
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && !boolFlags[flagName] {
					i++
					normalizedArgs = append(normalizedArgs, args[i])
				}
			}
		} else if strings.HasPrefix(arg, "-") {
			flagName := strings.TrimLeft(arg, "-")
			if strings.Contains(flagName, "=") {
				normalizedArgs = append(normalizedArgs, arg)
			} else {
				normalizedArgs = append(normalizedArgs, arg)
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && !boolFlags[flagName] {
					i++
					normalizedArgs = append(normalizedArgs, args[i])
				}
			}
		} else {
			tickers = append(tickers, arg)
		}
	}

	if err := fs.Parse(normalizedArgs); err != nil {
		return nil, err
	}

	if filePath != "" {
		fileTickers, err := loadTickersFromFile(filePath)
		if err != nil {
			return nil, err
		}
		tickers = append(tickers, fileTickers...)
	}

	if len(tickers) == 0 {
		tickers = append(tickers, DefaultTickers...)
	}

	seen := make(map[string]struct{})
	var cleaned []string
	for _, t := range tickers {
		upper := strings.ToUpper(strings.TrimSpace(t))
		if upper != "" {
			if _, exists := seen[upper]; !exists {
				seen[upper] = struct{}{}
				cleaned = append(cleaned, upper)
			}
		}
	}

	finalView := "cards"
	if table || view == "table" {
		finalView = "table"
	} else if grouped || view == "grouped" {
		finalView = "grouped"
	}

	cfg := &Config{
		Tickers:     cleaned,
		FilePath:    filePath,
		RankBy:      rankBy,
		ExportPath:  exportPath,
		NoCache:     noCache,
		Plain:       plain,
		Years:       years,
		Concurrency: concurrency,
		View:        finalView,
	}

	if minROIC != -999999 {
		cfg.MinROIC = &minROIC
	}
	if maxDepFactor != -999999 {
		cfg.MaxDepFactor = &maxDepFactor
	}

	return cfg, nil
}

func loadTickersFromFile(fp string) ([]string, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("ticker file not found: %s", fp)
	}
	defer f.Close()

	var tickers []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			tickers = append(tickers, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tickers, nil
}
