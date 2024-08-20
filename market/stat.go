package market

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type CategoryType struct {
	Date          string  `json:"d"`
	Name          string  `json:"n"`
	PercentChange float64 `json:"pc"`
	Points        float64 `json:"v"`
	Si            float64 `json:"si"`
	Timestamp     float64 `json:"t"`
}

type Response struct {
	R struct {
		Categories map[string]CategoryType `json:"Indices"`
	} `json:"R"`
}

type MarketData map[string]CategoryType

func FetchMarketStat() MarketData {
	var responseFromServer Response
	var connectionToken = os.Getenv("NEPSE_CONNECTION_TOKEN")
	url := "https://merolagani.com/signalr/send?transport=serverSentEvents&connectionToken=" + connectionToken
	data := `data=%7B%22H%22%3A%22stocktickerhub%22%2C%22M%22%3A%22GetAllStocks%22%2C%22A%22%3A%5B%5D%2C%22I%22%3A0%7D`
	body := []byte(data)

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	r.Header.Add("Accept", "application/json")
	r.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

	if err != nil {
		log.Fatal("Error while creating request: ", err)
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Fatal("Error while sending request: ", err)
	}

	defer res.Body.Close()

	derr := json.NewDecoder(res.Body).Decode(&responseFromServer)
	if derr != nil {
		log.Fatal("Error while parsing response json : ", derr)
	}

	return responseFromServer.R.Categories
}
