package nepse

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

type Quote struct {
	PercentageChange string
	Positive         bool
	Price            string
}

type QuotesBySymbol map[string]Quote

func Scrape() *Quote {
	// Instantiate default collector
	c := colly.NewCollector()
	var marketPrice, percentageChange string
	var positive bool

	// Extract comment
	c.OnHTML("table tbody tr", func(e *colly.HTMLElement) {
		if e.ChildText("th:nth-child(1)") == "Market Price" {
			marketPrice = e.ChildText("td:nth-child(2)")
		}
		if e.ChildText("th:nth-child(1)") == "% Change" {
			percentageChange = e.ChildText("td:nth-child(2)")
			positive = !strings.Contains(percentageChange, "-")
		}
	})

	c.Visit("https://merolagani.com/CompanyDetail.aspx?symbol=NABIL")

	quoteValue := &Quote{
		PercentageChange: percentageChange,
		Price:            marketPrice,
		Positive:         positive,
	}

	return quoteValue
}