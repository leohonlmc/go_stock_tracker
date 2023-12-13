package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

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

        var wsMsg WSMessage
        if err := json.Unmarshal(msg, &wsMsg); err != nil {
            log.Println("Error decoding message:", err)
            continue
        }

        switch wsMsg.Action {
        case "getStock":
            stockTicker := wsMsg.Ticker
            go func() {
                // one minute ticker
                ticker := time.NewTicker(2 * time.Second)
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

        case "getStocks":
            stocks := wsMsg.Tickers

            var wg sync.WaitGroup
            results := make(chan StockData, len(stocks))
            
            for _, ticker := range stocks {
                wg.Add(1)
                go fetchStockData(results, ticker, &wg)
            }

            wg.Wait()
            close(results)

            allStocksData := make([]StockData, 0, len(stocks))
            for stockData := range results {
                allStocksData = append(allStocksData, stockData)
            }

            response, err := json.Marshal(map[string]interface{}{
                "action": "stocksData",
                "data": allStocksData,
            })

            if err != nil {
                log.Println("Error encoding response:", err)
                continue
            }
            
            // send response
            if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
                log.Println("Write error:", err)
                break
            }
        }
    }
}