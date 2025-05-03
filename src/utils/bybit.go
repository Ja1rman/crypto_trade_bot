package utils

// all tickers curl -L -X GET 'https://api.bybit.com/v5/market/tickers?category=spot'
import (
    "context"
	"encoding/json"
	"strconv"


    bybit "github.com/bybit-exchange/bybit.go.api"

	"crypto_trading/src/analyzer"
	"crypto_trading/src/handlers"
	"crypto_trading/src/logger"
	"crypto_trading/src/config"
)

type InstrumentInfo struct {
	Symbol        string                 `json:"symbol"`
	BaseCoin      string                 `json:"baseCoin"`
	QuoteCoin     string                 `json:"quoteCoin"`
	Innovation    string                 `json:"innovation"`
	Status        string                 `json:"status"`
	MarginTrading string                 `json:"marginTrading"`
	StTag         string                 `json:"stTag"`
	LotSizeFilter struct {
		BasePrecision  string `json:"basePrecision"`
		QuotePrecision string `json:"quotePrecision"`
		MinOrderQty  string `json:"minOrderQty"`
		MaxOrderQty string `json:"maxOrderQty"`
		MinOrderAmt  string `json:"minOrderAmt"`
		MaxOrderAmt string `json:"maxOrderAmt"`
	} `json:"lotSizeFilter"`
	PriceFilter struct {
		TickSize string `json:"tickSize"`
	} `json:"priceFilter"`
	RiskParams struct {
		PriceLimitRatioX string `json:"priceLimitRatioX"`
		PriceLimitRatioY string `json:"priceLimitRatioY"`
	} `json:"riskParameters"`
}

var client = bybit.NewBybitHttpClient("", "", bybit.WithBaseURL(bybit.MAINNET))

func SetBybitSettings() *analyzer.CurrencySettings {
	allSymbols := getAllBybitTickers()
	allSymbols = filterNonActiveBybitSymbols(allSymbols)
	tickersList := formatBybitSymbols(allSymbols)
	routes, infoMap := GenerateRoutesAndInfoMap(tickersList)
	result := &analyzer.CurrencySettings{
		MIN_PROFIT: 0.002,
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
		COMMISSION: 0.0019, // комса https://www.bybit.com/ru-RU/announcement-info/
		PRICES: handlers.Prices{
			Cache: make(map[string]handlers.OrderBookData),
		},
	}
	return result
}

func getAllBybitTickers() []InstrumentInfo {
	var result []InstrumentInfo

	params := map[string]interface{}{
		"category": "spot",
	}

	response, err := client.NewUtaBybitServiceWithParams(params).GetInstrumentInfo(context.Background())
	if err != nil {
		logger.Logger.Fatalf("Ошибка при получении информации: %v\n", err)
		return result
	}

	resultMap, ok := response.Result.(map[string]interface{})
	if !ok {
		logger.Logger.Fatalf("Неверный формат результата\n")
		return result
	}

	rawList, ok := resultMap["list"].([]interface{})
	if !ok {
		logger.Logger.Fatalf("Поле list отсутствует или в неверном формате\n")
		return result
	}

	for _, rawItem := range rawList {
		itemMap, ok := rawItem.(map[string]interface{})
		if !ok {
			return result
		}

		jsonBytes, err := json.Marshal(itemMap)
		if err != nil {
			logger.Logger.Fatalf("Ошибка при сериализации itemMap: %v\n", err)
			return result
		}

		var info InstrumentInfo
		if err := json.Unmarshal(jsonBytes, &info); err != nil {
			logger.Logger.Fatalf("Ошибка при десериализации itemMap: %v\n", err)
			return result
		}

		result = append(result, info)
	}
	return result
}

func filterNonActiveBybitSymbols(symbols []InstrumentInfo) []InstrumentInfo {
	var newSymbols []InstrumentInfo
	for _, symbol := range symbols {
		if symbol.Status != "Trading" ||
		   symbol.StTag != "0" {
			continue
		}
		newSymbols = append(newSymbols, symbol)
	}
	return newSymbols
}

func formatBybitSymbols(symbols []InstrumentInfo) map[string]config.Info {
	result := make(map[string]config.Info) 
	for _, symbol := range symbols {
		tickSize, _ := strconv.ParseFloat(symbol.PriceFilter.TickSize, 64)
		basePrecision, _ := strconv.ParseFloat(symbol.LotSizeFilter.BasePrecision, 64)
		quotePrecision, _ := strconv.ParseFloat(symbol.LotSizeFilter.QuotePrecision, 64)
		info := config.Info{
			BaseCoin:  symbol.BaseCoin,
			QuoteCoin: symbol.QuoteCoin,
			Precision: config.Precision{
				BasePrecision:  1/basePrecision,
				QuotePrecision: 1/quotePrecision,
			},
			TickSize: 1/tickSize,
		}
		result[symbol.Symbol] = info
	}
	return result
}
