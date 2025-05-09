package analyzer

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
	"slices"
	"sync"

	"crypto_trading/src/alerting"
	"crypto_trading/src/config"
	"crypto_trading/src/handlers"
	"crypto_trading/src/stats"
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
	PRICES handlers.Prices
	MARKET string
}

func (currSet *CurrencySettings) LaunchInfiniteAnalyze() {
	currSet.searchTradingPaths()
	var maxAnalyzeTime int64 = 0

	for {
		start := time.Now().UnixMicro()
		currSet.analyzeForAllPaths()
		now := time.Now()
		elapsed := now.UnixMicro() - start
		if elapsed > maxAnalyzeTime {
			maxAnalyzeTime = elapsed
		}
		if now.Unix() % 15 == 0 {
			go stats.PushToPrometheus(stats.AnalyzeDuration, currSet.MARKET, float64(maxAnalyzeTime), "duration")
			// test
			randomFloat := 600 + rand.Float64()*(1200-600)
			go stats.PushToPrometheus(stats.WalletBalance, currSet.MARKET, randomFloat, "USDT")
			randomFloat = 3 + rand.Float64()*(13-3)
			go stats.PushToPrometheus(stats.WalletBalance, currSet.MARKET, randomFloat, "BTC")
			go stats.PushToPrometheus(stats.OrdersDuration, currSet.MARKET, float64(elapsed*3), "duration")

			maxAnalyzeTime = 0
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
		symbol, info, err := currSet.FindSymbol(firstCurr, secondCurr)
		if err != nil {
			if firstCurr != secondCurr {
				logger.Logger.Println(err, currencies)
			}
			return
		}
		orderBook := currSet.PRICES.GetOrderBookData(symbol)
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
		// start := time.Now().UnixMilli()
		// trade.ProcessCycle(balances[1], balances[1]/balances[0], currencies, currSet.SUBSCRIBE_TICKERS_LIST)
		// now := time.Now()
		// elapsed := now.UnixMilli() - start
		//go stats.PushToPrometheus(stats.OrdersDuration, currSet.MARKET, float64(elapsed), "duration")

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
	balances, err := currSet.GetWalletBalance(currencies)
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
		go stats.PushToPrometheus(stats.WalletBalance, currSet.MARKET, bal, currency)
		if bal < limits.StopBalance {
			alerting.SendMessage("Не хватает денег!!!")
			panic("Выход за критический порог баланса. Остановка работы...")
		}
	}
}

func (currSet *CurrencySettings) ProcessCycle(startSize float64, firstOrderPrice float64, pairsNames []string, exchange string) {
	qty, err := currSet.ProcessFirstPair(startSize, firstOrderPrice, pairsNames)
	if err!= nil || qty <= 0.0 {
		logger.Logger.Printf("Не получилось совершить сделку, error: %s\n", err)
		return
	}
	err = currSet.SellAll(pairsNames, qty)
	if err != nil {
		logger.Logger.Println(err)
	}
}

func (currSet *CurrencySettings) SellAll(pairsNames []string, qty float64) error {
	errMsg := fmt.Sprintf("Не удалось продать все валюты, pairsNames: %v", pairsNames)
	for i, pairName := range pairsNames {
		if i == 0 || i == len(pairsNames)-1 {
			continue
		}
		symbol, info, err := currSet.FindSymbol(pairsNames[i], pairsNames[i+1])
		if err != nil {
			return fmt.Errorf("%s ошибка: %s, %v", errMsg, err, info)
		}
		qty, err = currSet.CreateSellOrder(info, qty, symbol, pairName)
		if err != nil {
			return fmt.Errorf("%s ошибка: %s", errMsg, err)
		}
	}
	return nil
}

func (currSet *CurrencySettings) CreateSellOrder(info config.Info, size float64, symbol string, coin string) (float64, error) {
	coinType := "baseCoin"
	side := "Sell"
	if info.BaseCoin == coin {
		size = RoundCustomStep(size, info.Precision.BasePrecision)
	} else {
		coinType = "quoteCoin"
		side = "Buy"
		size = RoundCustomStep(size, info.Precision.QuotePrecision)
	}
	inversion := true
	if coinType == "quoteCoin" {
		inversion = false
	}
	return currSet.CreateOrderAndGetResult(symbol, 1., size, coinType, side, "Market", inversion)
}

func (currSet *CurrencySettings) ProcessFirstPair(size float64, price float64, pairsNames []string) (float64, error) {
	symbol, info, err := currSet.FindSymbol(pairsNames[0], pairsNames[1])
	if err != nil {
		return 0.0, err
	}
	return currSet.createFirstOrder(info, size, price, symbol, pairsNames[1])
}

func (currSet *CurrencySettings) createFirstOrder(info config.Info, size float64, price float64, symbol string, coin string) (float64, error) {
	side := "Buy"
	coinType := "baseCoin"
	if info.BaseCoin != coin {
		side = "Sell"
		size = size * price
		price = 1.0 / price
	}
	size = RoundCustomStep(size, info.Precision.BasePrecision)
	price = RoundCustomStep(price, info.TickSize)
	inversion := true
	if side == "Buy" {
		inversion = false
	}
	return currSet.CreateOrderAndGetResult(symbol, price, size, coinType, side, "Limit", inversion)
}

func (currSet *CurrencySettings) CreateOrderAndGetResult(symbol string, price float64, qty float64, coinType string, side string, orderType string, inversion bool) (float64, error) {
	if currSet.MARKET == "bybit" {
		orderId, err := trade.CreateBybitOrder(symbol, price, qty, coinType, side, orderType)
		if err != nil {
			return 0.0, err
		}
		orderInfo, err := trade.GetBybitOrderInfo(orderId)
		if err != nil {
			return 0.0, err
		}
		if len(orderInfo) == 0 {
			return 0.0, fmt.Errorf("orderInfo not found")
		}
		return trade.SumBybitQty(orderInfo, inversion), nil
	}
	if currSet.MARKET == "mexc" {
		orderId, err := trade.CreateMexcOrder(symbol, price, qty, coinType, side, orderType)
		if err != nil {
			return 0.0, err
		}
		orderInfo, err := trade.GetMexcOrderInfo(symbol, orderId)
		if err != nil {
			return 0.0, err
		}
		return trade.SumMexcQty(orderInfo, inversion), nil
	}
	return 0.0, fmt.Errorf("Неизвестный тип биржи")
}

func (currSet *CurrencySettings) GetWalletBalance(coins string) (map[string]float64, error) {
	if currSet.MARKET == "bybit" {
		return trade.GetBybitWalletBalance(coins)
	}
	if currSet.MARKET == "mexc" {
		return trade.GetMexcWalletBalance()
	}
	return map[string]float64{}, fmt.Errorf("Неизвестный тип биржи")
}

func (currSet *CurrencySettings) FindSymbol(firstCurrency string, secondCurrency string) (string, config.Info, error) {
	symbol := firstCurrency + secondCurrency
	if info, exists := currSet.SUBSCRIBE_TICKERS_LIST[symbol]; exists {
		return symbol, info, nil
	}
	symbol = secondCurrency + firstCurrency
	if info, exists := currSet.SUBSCRIBE_TICKERS_LIST[symbol]; exists {
		return symbol, info, nil
	}
	return "", config.Info{}, fmt.Errorf("pairs %s not found in config", symbol)
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

func RoundCustomStep(number, step float64) float64 {
	return math.Floor(number*step) / step
}
