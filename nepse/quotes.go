package nepse

import (
// "fmt"
// "strings"
//
// "github.com/gocolly/colly/v2"
)

type Quote struct {
	PercentageChange string
	Positive         bool
	Price            string
	MarketCap        string
}

type QuotesBySymbol map[string]Quote

func ScrapeBySymbol(symbol string) *Quote {
	// c := colly.NewCollector()
	// var marketPrice, percentageChange, marketCap string
	// var positive bool
	//
	// c.OnHTML("table tbody tr", func(e *colly.HTMLElement) {
	// 	if e.ChildText("th:nth-child(1)") == "Market Price" {
	// 		marketPrice = e.ChildText("td:nth-child(2)")
	// 	}
	// 	if e.ChildText("th:nth-child(1)") == "% Change" {
	// 		percentageChange = e.ChildText("td:nth-child(2)")
	// 		positive = !strings.Contains(percentageChange, "-")
	// 	}
	//
	// 	if e.ChildText("th:nth-child(1)") == "Market Capitalization" {
	// 		marketCap = e.ChildText("td:nth-child(2)")
	// 	}
	// })
	//
	// url := fmt.Sprintf("https://merolagani.com/CompanyDetail.aspx?symbol=%v", symbol)
	// c.Visit(url)

	quoteValue := &Quote{
		PercentageChange: "+ 0.5%",
		Price:            "690",
		Positive:         true,
		MarketCap:        "123,234452",
	}

	return quoteValue
}
