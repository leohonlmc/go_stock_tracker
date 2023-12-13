package main

import (
	"math"
	"strconv"
	"sync"
	"time"
)

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
        timeParsed, err := time.Parse("2006-01-02 15:04:05", timestamp)
        if err != nil {
            return 0, 0, 0, err
        }

        if timeParsed.After(latestTime) {
            latestTime = timeParsed
            latestPrice, err = strconv.ParseFloat(entry.Close, 64)
            if err != nil {
                return 0, 0, 0, err
            }
        }

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