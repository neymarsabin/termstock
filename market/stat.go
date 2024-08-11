package market

import "github.com/gocolly/colly/v2"

type MarketData struct {
	Trading          MarketDataType
	HotelsAndTourism MarketDataType
	Nepse            MarketDataType
	Sensitive        MarketDataType
	Float            MarketDataType
	SenFloat         MarketDataType
	Banking          MarketDataType
	Investment       MarketDataType
	MutualFund       MarketDataType
	LifeInsurance    MarketDataType
	Others           MarketDataType
	MicroFinance     MarketDataType
	Finance          MarketDataType
	AsOfUpdate       string
}

type MarketDataType struct {
	Points    string
	Change    string
	MarketCap string
}

func FetchMarketStat() {
	c := colly.NewCollector()

	c.OnHTML("", func(e *colly.HTMLElement) {
		// find and insert values in the MarketData
	})

	c.Visit("https://merolagani.com/LatestMarket.aspx")
}
