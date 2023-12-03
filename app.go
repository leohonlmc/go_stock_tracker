package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type StockData struct {
	MetaData   MetaData              `json:"Meta Data"`
	TimeSeries map[string]StockEntry `json:"Time Series (5min)"`
}

type MetaData struct {
	Information string `json:"1. Information"`
	Symbol      string `json:"2. Symbol"`
	LastRefresh string `json:"3. Last Refreshed"`
	Interval    string `json:"4. Interval"`
	OutputSize  string `json:"5. Output Size"`
	TimeZone    string `json:"6. Time Zone"`
}

type StockEntry struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

func main() {
	stocks := []string{"IBM", "PLTR"}
	c := make(chan StockData)

	var wg sync.WaitGroup

	for _, stock := range stocks {
		wg.Add(1)
		go fetchStockData(c, stock, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for data := range c {
		fmt.Println(data) 
	}
}

func fetchStockData(c chan StockData, stock string, wg *sync.WaitGroup) {
	
	defer wg.Done()

	// load .env file
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	key := os.Getenv("API_KEY")
	fmt.Println(key)

	resp, err := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=" + stock + "&interval=5min&apikey=" + key)
	if err != nil {
		log.Println(err) 
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var data StockData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println(err)
		return
	}

	c <- data
}