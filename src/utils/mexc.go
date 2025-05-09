package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"crypto_trading/src/analyzer"
	"crypto_trading/src/config"
	"crypto_trading/src/handlers"
	"crypto_trading/src/logger"
)

type ExchangeInfo struct {
	Timezone        string   `json:"timezone"`
	ServerTime      int64    `json:"serverTime"`
	RateLimits      []any    `json:"rateLimits"`
	ExchangeFilters []any    `json:"exchangeFilters"`
	Symbols         []Symbol `json:"symbols"`
}

type Symbol struct {
	Symbol                   string   `json:"symbol"`
	Status                   string   `json:"status"`
	BaseAsset                string   `json:"baseAsset"`
	BaseAssetPrecision       int      `json:"baseAssetPrecision"`
	QuoteAsset               string   `json:"quoteAsset"`
	QuotePrecision           int      `json:"quotePrecision"`
	QuoteAssetPrecision      int      `json:"quoteAssetPrecision"`
	BaseCommissionPrecision  int      `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision int      `json:"quoteCommissionPrecision"`
	OrderTypes               []string `json:"orderTypes"`
	QuoteOrderQtyMarketAllowed bool   `json:"quoteOrderQtyMarketAllowed"`
	IsSpotTradingAllowed     bool     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed   bool     `json:"isMarginTradingAllowed"`
	QuoteAmountPrecision     string   `json:"quoteAmountPrecision"`
	BaseSizePrecision        string   `json:"baseSizePrecision"`
	Permissions              []string `json:"permissions"`
	Filters                  []any    `json:"filters"`
	MaxQuoteAmount           string   `json:"maxQuoteAmount"`
	MakerCommission          string   `json:"makerCommission"`
	TakerCommission          string   `json:"takerCommission"`
	TradeSideType            int   `json:"tradeSideType"`
}

func SetMexcSettings() *analyzer.CurrencySettings {
	allSymbols := getAllMexcSymbols()
	allSymbols = filterNonActiveMexcSymbols(allSymbols)
	tickersList := formatMexcSymbols(allSymbols)
	routes, infoMap := GenerateRoutesAndInfoMap(tickersList)
	result := &analyzer.CurrencySettings{
		MIN_PROFIT: 0.0001,
		MIN_MONEY_DEAL: 0.3,
		CURRENCY_ROUTES: routes,
		SUBSCRIBE_TICKERS_LIST: infoMap,
		START_CURRENCIES: analyzer.StartCurrencies{
			Cache: map[string]analyzer.MoneyLimits{
			"USDT": {StopBalance: 800, MaxDealPrice: 800},
			"USDC": {StopBalance: 800, MaxDealPrice: 800},
			//"BTC": {0.01, 0.01},
			//"ETH": {0.5, 0.5},
			//"EUR": {800, 800},
		}},
		COMMISSION: 0.0005, // комса https://www.bybit.com/ru-RU/announcement-info/
		PRICES: handlers.Prices{
			Cache: make(map[string]handlers.OrderBookData),
		},
		MARKET: "mexc",
	}
	return result
}

func getAllMexcSymbols() []Symbol {
	url := "https://api.mexc.com/api/v3/exchangeInfo"

	resp, err := http.Get(url)
	if err != nil {
		logger.Logger.Fatalf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Fatalf("неожиданный статус ответа: %s", resp.Status)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Fatalf("ошибка при чтении тела ответа: %v", err)
	}
	var exchangeInfo ExchangeInfo
	err = json.Unmarshal(body, &exchangeInfo)
	if err != nil {
		logger.Logger.Fatalf("ошибка при парсинге JSON: %v", err)
	}
	return exchangeInfo.Symbols
}

func filterNonActiveMexcSymbols(symbols []Symbol) []Symbol {
	var newSymbols []Symbol
	for _, symbol := range symbols {
		if symbol.TradeSideType != 1 ||
		   symbol.Status != "1" ||
		   !symbol.IsSpotTradingAllowed {
			continue
		}
		newSymbols = append(newSymbols, symbol)
	}
	return newSymbols
}

func formatMexcSymbols(symbols []Symbol) map[string]config.Info {
	result := make(map[string]config.Info) 
	for _, symbol := range symbols {
		info := config.Info{
			BaseCoin:  symbol.BaseAsset,
			QuoteCoin: symbol.QuoteAsset,
			Precision: config.Precision{
				BasePrecision:  float64(10^symbol.BaseAssetPrecision),
				QuotePrecision: float64(10^symbol.QuoteAssetPrecision),
			},
			TickSize: float64(10^symbol.QuotePrecision),
		}
		result[symbol.Symbol] = info
	}
	return result
}
