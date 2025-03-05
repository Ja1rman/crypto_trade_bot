package utils

// all tickers curl -L -X GET 'https://api.bybit.com/v5/market/tickers?category=spot'
import (
	"fmt"
    "context"
	"os"
	"encoding/json"

    bybit "github.com/wuhewuhe/bybit.go.api"

	"crypto_trading/src/logger"
	"crypto_trading/src/config"
)

type InstrumentInfo struct {
	Symbol      string `json:"symbol"`
	BaseCoin    string `json:"baseCoin"`
	QuoteCoin   string `json:"quoteCoin"`
	LotSizeFilter struct {
		BasePrecision  string `json:"basePrecision"`
		QuotePrecision string `json:"quotePrecision"`
	} `json:"lotSizeFilter"`
}

type InstrumentInfoResult struct {
	Category string           `json:"category"`
	List     []InstrumentInfo `json:"list"`
}


func writeToFile(values interface{}, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		logger.Logger.Fatal(err)
	}

	_, err = file.Write(data)
	if err != nil {
		logger.Logger.Fatal(err)
	}
}

func getTickers(apiKey string, apiSecret string) {
	client := bybit.NewBybitHttpClient(apiKey, apiSecret, bybit.WithBaseURL(bybit.MAINNET))
	params := map[string]interface{}{"category": "spot"}
	results, _ := client.NewUtaBybitServiceWithParams(params).GetMarketTickers(context.Background())
	writeToFile(results, "results.json")
}

func GetCurrenciesPrecision() {
	client := bybit.NewBybitHttpClient("", "", bybit.WithBaseURL(bybit.TESTNET))
	precisionMap := make(map[string]string)

	for symbol := range config.SUBSCRIBE_TICKERS_LIST {
		params := map[string]interface{}{"category": "spot", "symbol": symbol}
		response, err := client.NewUtaBybitServiceWithParams(params).GetInstrumentInfo(context.Background())
		if err != nil {
			fmt.Printf("Ошибка при получении информации для символа %s: %v\n", symbol, err)
			continue
		}

		if response.RetCode != 0 {
			fmt.Printf("API вернуло ошибку для символа %s: %s\n", symbol, response.RetMsg)
			continue
		}

		var result InstrumentInfoResult
		resultBytes, err := json.Marshal(response.Result)
		if err != nil {
			fmt.Printf("Ошибка при преобразовании result в JSON: %v\n", err)
			continue
		}
		err = json.Unmarshal(resultBytes, &result)
		if err != nil {
			fmt.Printf("Ошибка при разборе result: %v\n", err)
			continue
		}

		// Обрабатываем список инструментов
		for _, instrument := range result.List {
			// Записываем basePrecision для baseCoin
			if _, exists := precisionMap[instrument.BaseCoin]; !exists {
				precisionMap[instrument.BaseCoin] = instrument.LotSizeFilter.BasePrecision
			}

			// Записываем quotePrecision для quoteCoin
			if _, exists := precisionMap[instrument.QuoteCoin]; !exists {
				precisionMap[instrument.QuoteCoin] = instrument.LotSizeFilter.QuotePrecision
			}
		}
	}

	// Выводим результат
	writeToFile(precisionMap, "results.json")
}

