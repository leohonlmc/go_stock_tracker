package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

// struct for stock data
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

// main function
func main() {

	r := gin.Default()

    stocks := []string{"IBM", "PLTR"}
    c := make(chan StockData)

	// GET /stocks
	r.GET("/startFetching", func(ctx *gin.Context) {
		// Fetch data for each stock when the program starts
		var wg sync.WaitGroup
		for _, stock := range stocks {
			wg.Add(1)
			go fetchStockData(c, stock, &wg)
		}
	
		// Time interval for fetching data each 5 minutes
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
	
		go func() {
			for {
				select {
				case <-ticker.C:
					var wg sync.WaitGroup
					for _, stock := range stocks {
						wg.Add(1)
						go fetchStockData(c, stock, &wg)
					}
					wg.Wait() // Wait for all goroutines to finish for this ticker cycle
				}
			}
		}()
	
		ctx.JSON(200, gin.H{
			"message": "Stock data fetching started",
		})
	})

    // Continuously process data from the channel
    go func() {
        for data := range c {
            fmt.Println(data) 
        }
    }()

    // start http server
    log.Println("HTTP server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

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
