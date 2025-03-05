package analyzer

import (
	"fmt"
	"math"
	"time"

	"crypto_trading/src/handlers"
	"crypto_trading/src/alerting"
)

var(
	MIN_PROFIT float64 = 0.0158 // около 0.58% это комса
	MIN_MONEY_DEAL float64 = 0.3
	START_CURRENCIES = map[string]float64{
		"USDT": 1000,
		"USDC": 1000,
		"BTC": 0.01,
		"ETH": 0.5,
	}
	lastAlertTimes = make(map[string]time.Time)
)

type TradingPairs []handlers.OrderBookData

func StartAnalyzies() {
	go launchInfiniteTrying()
	select {}
}

func launchInfiniteTrying() {
	lastNotifier := time.Now()

	for {
		start := time.Now().UnixNano()
		tryDifferentStartCurrencies()
		//elapsed := time.Now().UnixNano() - start
		//fmt.Println(fmt.Sprintf("Время выполнения анализа: %d", elapsed))
		now := time.Now()
		if now.Sub(lastNotifier) > time.Hour {
			elapsed := now.UnixNano() - start
			go alerting.SendMessage(fmt.Sprintf("Время выполнения анализа: %d", elapsed))
			// go alerting.PrintPrices()
			lastNotifier = now
		}
	}
}

func tryDifferentStartCurrencies() {
	for currency, balance := range START_CURRENCIES {
		tradingPathsSearching(currency, balance)
	}
}

func tradingPathsSearching(currency string, balance float64) {
	handlers.PRICES.Lock()
	defer handlers.PRICES.Unlock()
	for secondCurrency, secondCurrencyOrderBookData := range handlers.PRICES.Cache[currency] {
		_, exists := handlers.PRICES.Cache[secondCurrency]
		if !exists {
			continue
		}
		for thirdCurrency, thirdCurrencyOrderBookData := range handlers.PRICES.Cache[secondCurrency] {
			lastOrderBookData, exists := handlers.PRICES.Cache[thirdCurrency]
			if !exists {
				continue
			}
			_, exists = lastOrderBookData[currency]
			if !exists {
				continue
			}
			tradingPairs := TradingPairs{secondCurrencyOrderBookData, thirdCurrencyOrderBookData, lastOrderBookData[currency]}
			if !CheckZeros(tradingPairs) {
				continue
			}
			startSize := CalculateMarketSize(tradingPairs, balance) 
			if startSize == 0 {
				continue
			}
			profit := AnalizeOfferProfit(tradingPairs, startSize)
			pairsNames := fmt.Sprintf("%s/%s/%s", currency, secondCurrency, thirdCurrency)
			BuyDecision(tradingPairs, profit, pairsNames, startSize*tradingPairs[0].Bid.Price, MIN_MONEY_DEAL * balance)
		}
	}
}

func CheckZeros(tradingPairs TradingPairs) bool {
	firstPair, secondPair, thirdPair := tradingPairs[0], tradingPairs[1], tradingPairs[2]

	if firstPair.Bid.Price == 0 ||
	   secondPair.Bid.Price == 0 ||
	   thirdPair.Ask.Price == 0 {
		return false
	}
	return true
}

func CalculateMarketSize(tradingPairs TradingPairs, balance float64) float64 {
	firstPair, secondPair, thirdPair := tradingPairs[0], tradingPairs[1], tradingPairs[2]
	sizeFirstCurrency := math.Min(balance / firstPair.Bid.Price, firstPair.Bid.Size)
	if (sizeFirstCurrency > thirdPair.Ask.Size) {
		sizeFirstCurrency = thirdPair.Ask.Size
	}
	sizeSecondCurrency := sizeFirstCurrency / secondPair.Bid.Price
	if (sizeSecondCurrency > secondPair.Bid.Size) {
		sizeSecondCurrency = secondPair.Bid.Size
		sizeFirstCurrency = sizeSecondCurrency * secondPair.Bid.Price
	}
	return sizeFirstCurrency
}

func AnalizeOfferProfit(tradingPairs TradingPairs, startSize float64) float64 {
	firstPair, secondPair, thirdPair := tradingPairs[0], tradingPairs[1], tradingPairs[2]
	startBalance := startSize * firstPair.Bid.Price
	finalBalance := startSize / secondPair.Bid.Price / thirdPair.Ask.Price
	return finalBalance - startBalance
}

func BuyDecision(tradingPairs TradingPairs, profit float64, pairsNames string, balance float64, moneyAmount float64) {
	if profit/balance < MIN_PROFIT || balance < moneyAmount {
		return
	}


	currentTime := time.Now()
	lastAlertTime, exists := lastAlertTimes[pairsNames]
	if exists && currentTime.Sub(lastAlertTime).Minutes() < 1 {
		return
	}
	pricesMessage := "Цены:\n"
	for i := 0; i < 2; i++ {
		pricesMessage += fmt.Sprintf("%d: Bid Price=%.6f, Size=%.6f, Seq=%d\n", 
			i+1, 
			tradingPairs[i].Bid.Price, 
			tradingPairs[i].Bid.Size, 
			tradingPairs[i].Bid.Seq)
	}
	pricesMessage += fmt.Sprintf("3: Ask Price=%.6f, Size=%.6f, Seq=%d", 
		tradingPairs[2].Ask.Price, 
		tradingPairs[2].Ask.Size, 
		tradingPairs[2].Ask.Seq)
	message := fmt.Sprintf("Нашлись сделки с профитом *%.2f*$ от баланса (%.2f%%).\n%s\n%s, ", profit, profit/balance*100, pairsNames, pricesMessage)
	alerting.SendMessage(message)
	lastAlertTimes[pairsNames] = currentTime	
}

func roundCustomStep(number, step float64) float64 {
	return math.Trunc(number/step) * step
}
