package analyzer

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"time"
	"slices"

	"crypto_trading/src/alerting"
	"crypto_trading/src/config"
	"crypto_trading/src/handlers"

	"crypto_trading/src/logger"
	"crypto_trading/src/trade"
)

type MoneyLimits struct {
	StopPrice float64
	MaxDealPrice float64
}
type TradingPairs []handlers.OrderBookData

var(
	MIN_PROFIT float64 = 0.008 + 0.0058 // около 0.58% это комса
	MIN_MONEY_DEAL float64 = 0.3
	START_CURRENCIES = map[string]MoneyLimits{
		"USDT": {800, 800.},
		"USDC": {800, 800},
		//"BTC": {0.01, 0.01},
		//"ETH": {0.5, 0.5},
	}
	lastAlertTimes = make(map[string]time.Time)
)


func LaunchInfiniteAnalyze() {
	lastNotifier := time.Now().Add(3*time.Minute - time.Hour)

	for {
		start := time.Now().UnixNano()
		tryDifferentStartCurrencies()
		//elapsed := time.Now().UnixNano() - start
		//fmt.Println(fmt.Sprintf("Время выполнения анализа: %d", elapsed))
		now := time.Now()
		if now.Sub(lastNotifier) > time.Hour {
			elapsed := now.UnixNano() - start
			go alerting.SendMessage(fmt.Sprintf("Время выполнения анализа: %d, Горутин в работе: %d", elapsed, runtime.NumGoroutine()))
			lastNotifier = now
		}
	}
}

func tryDifferentStartCurrencies() {
	for currency := range START_CURRENCIES {
		tradingPathsSearching([]string{currency}, 1, 3, 3)
	}
}

func tradingPathsSearching(currencies []string, currentDepth int, minDepth int, maxDepth int) {
	currentDepth += 1
	if currentDepth > maxDepth {
		return
	}
	for _, currentCurrency := range config.CURRENCY_ROUTES[currencies[len(currencies)-1]] {
		if _, exists := config.CURRENCY_ROUTES[currentCurrency]; !exists || slices.Contains(currencies, currentCurrency) {
			continue
		}
		currencies = append(currencies, currentCurrency)
		if minDepth <= currentDepth && slices.Contains(config.CURRENCY_ROUTES[currentCurrency], currencies[0]) { 
			ProcessAnalyze(currencies) // можно распараллелить TODO
		}
		tradingPathsSearching(currencies, currentDepth, minDepth, maxDepth)
		currencies = currencies[:len(currencies)-1]
	}
}

func ProcessAnalyze(currencies []string) {
	balances := make([]float64, len(currencies)+1)
	balances[0] = START_CURRENCIES[currencies[0]].MaxDealPrice
	currencies = append(currencies, currencies[0])
	logPairs := TradingPairs{}
	for i, firstCurr := range currencies {
		if i == len(currencies)-1 {
			break
		}
		secondCurr := currencies[i+1]
		symbol, info, err := trade.FindSymbol(firstCurr, secondCurr)
		if err != nil {
			if firstCurr != secondCurr {
				logger.Logger.Println(err, currencies)
			}
			return
		}
		orderBook := handlers.GetOrderBookData(symbol)
		logPairs = append(logPairs, orderBook)
		remainder := 0.
		balances[i+1], remainder = ConvertCurrency(balances[i], firstCurr, secondCurr, info, orderBook)
		if balances[i+1] <= 0 {
			return
		}
		if remainder > 0 {
			newBalance := balances[i]-remainder
			for j := i; j >= 0; j-- {
				nextBalance := 0.
				if j > 0 {
					nextBalance = (balances[j-1] * newBalance) / balances[j]
				}
				balances[j] = newBalance
				newBalance = nextBalance
			}
		}
	}
	BuyDecision(logPairs, balances, currencies)
}

func BuyDecision(tradingPairs TradingPairs, balances []float64, currencies []string) {
	profit := balances[len(balances)-1] - balances[0]
	if (profit/balances[0] < MIN_PROFIT ||
		balances[0] < MIN_MONEY_DEAL*START_CURRENCIES[currencies[0]].MaxDealPrice) {
		return
	}

	currentTime := time.Now()
	currenciesString := strings.Join(currencies, "/")

	//logger.Logger.Printf("Попытка цикла для пары: %s\n", currenciesString)
	// BID - цена моментальной покупки. всё поменять в trade
	//trade.ProcessCycle(startSize, tradingPairs[0].Bid.Price, currencies)

	lastAlertTime, exists := lastAlertTimes[currenciesString]
	if !exists || exists && currentTime.Sub(lastAlertTime).Minutes() >= 1 {
		pricesMessage := "Цены:\n"
		for i := 0; i < len(tradingPairs); i++ {
			pricesMessage += fmt.Sprintf("%d: Bid Price=%.10f, Size=%.10f, Seq=%d, Ts=%d\n", 
				i+1, 
				tradingPairs[i].Bid.Price, 
				tradingPairs[i].Bid.Size, 
				tradingPairs[i].Bid.Seq,
				tradingPairs[i].Bid.Ts,
			)
			pricesMessage += fmt.Sprintf("Ask Price=%.10f, Size=%.10f, Seq=%d, Ts=%d\n", 
				tradingPairs[i].Ask.Price, 
				tradingPairs[i].Ask.Size, 
				tradingPairs[i].Ask.Seq,
				tradingPairs[i].Ask.Ts,
			)
		}
		message := fmt.Sprintf("Нашлись сделки с профитом *%.2f*$ от баланса (%.2f%%).\n%s\n%s\nbalances: %v", profit, profit/balances[0]*100, currencies, pricesMessage, balances)
		alerting.SendMessage(message)
		lastAlertTimes[currenciesString] = currentTime
	}
	/*
	balances, err := trade.GetWalletBalance("USDT,USDC")
	if err != nil {
		logger.Logger.Printf("Ошибка получения баланса: %s\n", err.Error())
		alerting.SendMessage("Ошибка получения баланса")
		return
	}
	alerting.SendMessage(fmt.Sprintf("Баланс: %v\n", balances))
	for currency, limit := range START_CURRENCIES {
		if balances[currency] < limit.StopPrice {
			alerting.SendMessage("Не хватает денег!!!")
			panic("Выход за критический порог баланса. Остановка работы...")
		}
	}*/
}

func ConvertCurrency(
	balance float64,
	firstCurr string,
	secondCurr string,
	info config.Info,
	orderBook handlers.OrderBookData,
) (float64, float64) {
	newBalance := 0.
	remainder := 0.
	if info.BaseCoin == firstCurr {
		if orderBook.Bid.Size == 0 {
			return 0, balance
		}
		balanceWithShift := math.Min(balance, orderBook.Bid.Size)
		newBalance = balanceWithShift * orderBook.Bid.Price
		remainder = balance - balanceWithShift
	} else {
		if orderBook.Ask.Price == 0 || orderBook.Ask.Size == 0 {
			return 0, balance
		}
		newBalance = math.Min(balance / orderBook.Ask.Price, orderBook.Ask.Size)
		oldBalance := newBalance * orderBook.Ask.Price
		remainder = balance - oldBalance
	}
	return newBalance, remainder
}
