package analyzer

import (
	"fmt"
	"math"
	"strings"
	"time"

	"crypto_trading/src/alerting"
	"crypto_trading/src/handlers"
	"crypto_trading/src/logger"
	"crypto_trading/src/trade"
)

var(
	MIN_PROFIT float64 = 0.0158 // около 0.58% это комса
	MIN_MONEY_DEAL float64 = 0.3
	START_CURRENCIES = map[string]MoneyLimits{
		"USDT": MoneyLimits {600, 600},
		//"USDC": MoneyLimits {600, 600},
		//"BTC": MoneyLimits {0.01, 0.01},
		//"ETH": MoneyLimits {0.5, 0.5},
	}
	lastAlertTimes = make(map[string]time.Time)
)

type MoneyLimits struct {
	StopPrice float64
	MaxDealPrice float64
}
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
	for currency, limit := range START_CURRENCIES {
		tradingPathsSearching(currency, limit.MaxDealPrice)
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
			pairsNames := []string{currency, secondCurrency, thirdCurrency}
			BuyDecision(tradingPairs, profit, pairsNames, startSize, MIN_MONEY_DEAL * balance)
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

func BuyDecision(tradingPairs TradingPairs, profit float64, pairsNames []string, startSize float64, moneyAmount float64) {
	orderSize := startSize * tradingPairs[0].Bid.Price
	if profit/orderSize < MIN_PROFIT || orderSize < moneyAmount {
		return
	}

	currentTime := time.Now()
	pairsNamesString := strings.Join(pairsNames, "/")

	logger.Logger.Printf("Попытка цикла для пары: %s\n", pairsNamesString)
	trade.ProcessCycle(startSize, tradingPairs[0].Bid.Price, pairsNames)

	lastAlertTime, exists := lastAlertTimes[pairsNamesString]
	if !exists || exists && currentTime.Sub(lastAlertTime).Minutes() >= 1 {
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
		message := fmt.Sprintf("Нашлись сделки с профитом *%.2f*$ от баланса (%.2f%%).\n%s\n%s, ", profit, profit/orderSize*100, pairsNames, pricesMessage)
		alerting.SendMessage(message)
		lastAlertTimes[pairsNamesString] = currentTime
	}

	balances, err := trade.GetWalletBalance()
	if err != nil {
		logger.Logger.Printf("Ошибка получения баланса: %s\n", err.Error())
		alerting.SendMessage("Ошибка получения баланса")
	}
	alerting.SendMessage(fmt.Sprintf("Баланс: %v\n", balances))
	for currency, limit := range START_CURRENCIES {
		if balances[currency] < limit.StopPrice {
			alerting.SendMessage("Не хватает денег!!!")
			panic("Выход за критический порог баланса. Остановка работы...")
		}
	}
}
