package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)


func fetchStockData(c chan StockData, stock string, wg *sync.WaitGroup) {
	
	// The counter is decremented by calling wg.Done()
	defer wg.Done()

	// load .env file
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	key := os.Getenv("API_KEY")

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