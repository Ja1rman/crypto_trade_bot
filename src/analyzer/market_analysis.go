package analyzer

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"time"
	"slices"
	"sync"

	"crypto_trading/src/alerting"
	"crypto_trading/src/config"
	"crypto_trading/src/handlers"

	"crypto_trading/src/logger"
	"crypto_trading/src/trade"
)

type MoneyLimits struct {
	StopBalance float64
	MaxDealPrice float64
}
type TradingPairs []handlers.OrderBookData

type StartCurrencies struct {
	sync.RWMutex
	Cache map[string]MoneyLimits
}

type CurrencySettings struct {
	MIN_PROFIT float64
	MIN_MONEY_DEAL float64
	CURRENCY_ROUTES map[string]map[string]struct{}
	SUBSCRIBE_TICKERS_LIST map[string]config.Info
	START_CURRENCIES StartCurrencies
	LAST_ALERT_SENT map[string]time.Time
	ALL_PATHS [][]string
	COMMISSION float64
}


func (currSet *CurrencySettings) LaunchInfiniteAnalyze() {
	currSet.searchTradingPaths()
	lastNotifier := time.Now().Add(3*time.Minute - time.Hour)

	for {
		start := time.Now().UnixNano()
		currSet.analyzeForAllPaths()
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

func (currSet *CurrencySettings) searchTradingPaths() {
	currSet.START_CURRENCIES.RLock()
	defer currSet.START_CURRENCIES.RUnlock()
	for currency := range currSet.START_CURRENCIES.Cache {
        currSet.tradingPathsSearching([]string{currency}, 1, 3, 5)
	}
}

func (currSet *CurrencySettings) tradingPathsSearching(currencies []string, currentDepth int, minDepth int, maxDepth int) {
	currentDepth += 1
	if currentDepth > maxDepth {
		return
	}
	for currentCurrency := range currSet.CURRENCY_ROUTES[currencies[len(currencies)-1]] {
		if _, exists := currSet.CURRENCY_ROUTES[currentCurrency];
		   !exists ||
		   slices.Contains(currencies, currentCurrency) ||
		   currentCurrency == "USDT" ||
		   currentCurrency == "USDC" {
			continue
		}
		newCurrencies := make([]string, len(currencies)+1)
        copy(newCurrencies, currencies)
        newCurrencies[len(currencies)] = currentCurrency
		if minDepth <= currentDepth {
			if _, exists := currSet.CURRENCY_ROUTES[currentCurrency][newCurrencies[0]]; exists {
				path := make([]string, len(newCurrencies)+1)
                copy(path, newCurrencies)
                path[len(newCurrencies)] = newCurrencies[0]
				currSet.ALL_PATHS = append(currSet.ALL_PATHS, path)
			}
			if _, exists := currSet.CURRENCY_ROUTES[currentCurrency]["USDT"]; exists && newCurrencies[0] == "USDC" {
				path := make([]string, len(newCurrencies)+1)
                copy(path, newCurrencies)
                path[len(newCurrencies)] = "USDT"
				currSet.ALL_PATHS = append(currSet.ALL_PATHS, path)
			}
			if _, exists := currSet.CURRENCY_ROUTES[currentCurrency]["USDC"]; exists && newCurrencies[0] == "USDT" {
				path := make([]string, len(newCurrencies)+1)
                copy(path, newCurrencies)
                path[len(newCurrencies)] = "USDC"
				currSet.ALL_PATHS = append(currSet.ALL_PATHS, path)
			}
		}
		currSet.tradingPathsSearching(newCurrencies, currentDepth, minDepth, maxDepth)
	}
}

func (currSet *CurrencySettings) analyzeForAllPaths() {
	n := 2

	var wg sync.WaitGroup
	chunkSize := (len(currSet.ALL_PATHS) + n - 1) / n

	for i := 0; i < n; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(currSet.ALL_PATHS) {
			end = len(currSet.ALL_PATHS)
		}
		chunk := currSet.ALL_PATHS[start:end]

		wg.Add(1)
		go func(currenciesChunk [][]string) {
			defer wg.Done()
			for _, currencies := range currenciesChunk {
				currSet.processAnalyze(currencies)
			}
		}(chunk)
	}

	wg.Wait()
}

func (currSet *CurrencySettings) processAnalyze(currencies []string) {
	balances := make([]float64, len(currencies))
	currSet.START_CURRENCIES.RLock()
	balances[0] = currSet.START_CURRENCIES.Cache[currencies[0]].MaxDealPrice
	currSet.START_CURRENCIES.RUnlock()
	tradingPairs := TradingPairs{}
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
		tradingPairs = append(tradingPairs, orderBook)
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
	currSet.buyDecision(tradingPairs, balances, currencies)
}

func (currSet *CurrencySettings) buyDecision(tradingPairs TradingPairs, balances []float64, currencies []string) {
	profit := balances[len(balances)-1] - balances[0]
	currSet.START_CURRENCIES.RLock()
	limits := currSet.START_CURRENCIES.Cache[currencies[0]]
	maxDealPrice := limits.MaxDealPrice
	currSet.START_CURRENCIES.RUnlock()
	coms := balances[0] * currSet.COMMISSION * float64(len(currencies)-1) 
	if (profit/balances[0] - coms < currSet.MIN_PROFIT ||
		balances[0] < currSet.MIN_MONEY_DEAL*maxDealPrice) {
		return
	}

	currentTime := time.Now()
	currenciesString := strings.Join(currencies, "/")

	if config.TRADE && limits.StopBalance <= limits.MaxDealPrice {
		logger.Logger.Printf("Попытка цикла для пары: %s\n", currenciesString)
		trade.ProcessCycle(balances[1], balances[1]/balances[0], currencies)
		currSet.checkBalance()
	}

	lastAlertTime, exists := currSet.LAST_ALERT_SENT[currenciesString]
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
		currSet.LAST_ALERT_SENT[currenciesString] = currentTime
	}
}


func (currSet *CurrencySettings) checkBalance() {
	currSet.START_CURRENCIES.Lock()
	defer currSet.START_CURRENCIES.Unlock()

	var keys []string
	for k := range currSet.START_CURRENCIES.Cache {
		keys = append(keys, k)
	}
	currencies := strings.Join(keys, ",")
	balances, err := trade.GetWalletBalance(currencies)
	if err != nil {
		logger.Logger.Printf("Ошибка получения баланса: %s\n", err.Error())
		alerting.SendMessage("Ошибка получения баланса")
		return
	}
	alerting.SendMessage(fmt.Sprintf("Баланс: %v\n", balances))
	usdBal := balances["USDT"] + balances["USDC"]
	for currency, limits := range currSet.START_CURRENCIES.Cache {
		bal := balances[currency]
		newLimits := MoneyLimits{
			StopBalance: limits.StopBalance,
			MaxDealPrice: bal * 0.995,
		}
		currSet.START_CURRENCIES.Cache[currency] = newLimits
		if currency == "USDT" || currency == "USDC" {
			bal = usdBal
		}
		if bal < limits.StopBalance {
			alerting.SendMessage("Не хватает денег!!!")
			panic("Выход за критический порог баланса. Остановка работы...")
		}
	}
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
