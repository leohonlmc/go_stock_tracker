package main

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/gin-contrib/cors"
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
        return true 
    },
}

// var latestStockData []StockData

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

func processStockData( ticker string) (latestPrice, highestPrice, lowestPrice float64, err error) {
    var latestTime time.Time
    highestPrice = -1
    lowestPrice = math.MaxFloat64

	c := make(chan StockData, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	go fetchStockData(c, ticker, &wg)
	wg.Wait()

	stockData := <-c

    for timestamp, entry := range stockData.TimeSeries {
        // Parse the timestamp
        timeParsed, err := time.Parse("2006-01-02 15:04:05", timestamp)
        if err != nil {
            return 0, 0, 0, err
        }

        // Check if this is the latest timestamp
        if timeParsed.After(latestTime) {
            latestTime = timeParsed
            latestPrice, err = strconv.ParseFloat(entry.Close, 64)
            if err != nil {
                return 0, 0, 0, err
            }
        }

        // Check for highest and lowest prices
        highPrice, err := strconv.ParseFloat(entry.High, 64)
        if err != nil {
            return 0, 0, 0, err
        }
        if highPrice > highestPrice {
            highestPrice = highPrice
        }

        lowPrice, err := strconv.ParseFloat(entry.Low, 64)
        if err != nil {
            return 0, 0, 0, err
        }
        if lowPrice < lowestPrice {
            lowestPrice = lowPrice
        }
    }

    return latestPrice, highestPrice, lowestPrice, nil
}

func wshandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Failed to upgrade to WebSocket:", err)
        return
    }
    defer conn.Close()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            break
        }

        var request map[string]string
        if err := json.Unmarshal(msg, &request); err != nil {
            log.Println("Error decoding message:", err)
            continue
        }

        if request["action"] == "getStock" {
            stockTicker := request["ticker"]
            go func() {
                ticker := time.NewTicker(50 * time.Millisecond) // Set your desired interval
                defer ticker.Stop()

                for range ticker.C {
                    latestPrice, highestPrice, lowestPrice, err := processStockData(stockTicker)
                    if err != nil {
                        log.Println("Error processing stock data:", err)
                        continue
                    }

                    response, err := json.Marshal(map[string]interface{}{
                        "action":       "stockData",
                        "ticker":       stockTicker,
                        "latestPrice":  latestPrice,
                        "highestPrice": highestPrice,
                        "lowestPrice":  lowestPrice,
                    })
                    if err != nil {
                        log.Println("Error encoding response:", err)
                        continue
                    }

                    if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
                        log.Println("Write error:", err)
                        break
                    }
                }
            }()
        }
    }
}

// main function
func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:3000"},
        AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
    }))

	r.GET("/ws", func(c *gin.Context) {
        wshandler(c.Writer, c.Request)
    })

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	r.Run(":" + port)
}


// stocks := []string{"IBM", "PLTR"}
// 	c := make(chan StockData, len(stocks))

// 	var wg sync.WaitGroup
// 	for _, stock := range stocks {
// 		wg.Add(1)
// 		go fetchStockData(c, stock, &wg)
// 	}
// 	wg.Wait()

// 	for i := 0; i < len(stocks); i++ {
// 		latestStockData = append(latestStockData, <-c)
// 	}

// 	r.GET("/highLevelData", func(ctx *gin.Context) {
// 		ticker := ctx.Query("ticker")
// 		latestPrice, highestPrice, lowestPrice, err := processStockData(ticker)

// 		if err != nil {
// 			ctx.JSON(500, gin.H{
// 				"message": "Error processing stock data",
// 				"error":   err,
// 			})
// 			return
// 		}

// 		ctx.JSON(200, gin.H{
// 			"latestPrice":  latestPrice,
// 			"highestPrice": highestPrice,
// 			"lowestPrice":  lowestPrice,
// 			"ticker":       ticker,
// 		})

// 	})

// 	r.GET("/startFetching", func(ctx *gin.Context) {
// 		ctx.JSON(200, gin.H{
// 			"message": "Stock data fetching successfully started",
// 			"data":    latestStockData,
// 		})
// 	})

// 	ticker := time.NewTicker(5 * time.Minute)
// 	go func() {
// 		for range ticker.C {
// 			var wg sync.WaitGroup
// 			for _, stock := range stocks {
// 				wg.Add(1)
// 				go fetchStockData(c, stock, &wg)
// 			}
// 			wg.Wait()

// 			tempData := make([]StockData, 0)
// 			for i := 0; i < len(stocks); i++ {
// 				tempData = append(tempData, <-c)
// 			}
// 			latestStockData = tempData
// 		}
// 	}()


