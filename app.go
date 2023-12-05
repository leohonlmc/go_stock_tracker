package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // Adjust the origin policy as needed
    },
}

var latestStockData []StockData // Global variable to store the latest stock data

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

func wshandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Failed to upgrade to WebSocket:", err)
        return
    }
    defer conn.Close()

    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            break
        }

        log.Printf("Received: %s", p)

        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println("Write error:", err)
            break
        }
    }
}

// main function
func main() {
	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
        wshandler(c.Writer, c.Request)
    })

    stocks := []string{"IBM", "PLTR"}
	c := make(chan StockData, len(stocks))

	var wg sync.WaitGroup
	for _, stock := range stocks {
		wg.Add(1)
		go fetchStockData(c, stock, &wg)
	}
	wg.Wait()

	for i := 0; i < len(stocks); i++ {
		latestStockData = append(latestStockData, <-c)
	}

	r.GET("/startFetching", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Stock data fetching successfully started",
			"data":    latestStockData,
		})
	})

	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			var wg sync.WaitGroup
			for _, stock := range stocks {
				wg.Add(1)
				go fetchStockData(c, stock, &wg)
			}
			wg.Wait()

			tempData := make([]StockData, 0)
			for i := 0; i < len(stocks); i++ {
				tempData = append(tempData, <-c)
			}
			latestStockData = tempData
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	r.Run(":" + port)
}

