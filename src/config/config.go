package config

import (
	"os"
	
	bybit "github.com/wuhewuhe/bybit.go.api"
)

var (
	API_KEY string = os.Getenv("API_KEY")
	API_KEY_SECRET string = os.Getenv("API_KEY_SECRET")
)

type Info struct {
	BaseCoin string `json:"BaseCoin"`
	QuoteCoin string `json:"QuoteCoin"`
	Precision Precision `json:"Precision"`
}

type Precision struct {
	BasePrecision  float64 `json:"BasePrecision"`
	QuotePrecision float64 `json:"QuotePrecision"`
}

var (
	NET = bybit.TESTNET

	SUBSCRIBE_TICKERS_LIST = map[string]Info{
		"ADAEUR": {
			"ADA", "EUR",
			Precision {
				0.01, 0.000001,
			},
		},
		"ADAUSDC": {
			"ADA", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"ADAUSDT": {
			"ADA", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"ALGOBTC": {
			"ALGO", "BTC",
			Precision {
				0.1, 1e-10,
			},
		},
		"ALGOUSDT": {
			"ALGO", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"APEUSDC": {
			"APE", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"APEUSDT": {
			"APE", "USDT",
			Precision {
				0.01, 0.00001,
			},
		},
		"APEXUSDC": {
			"APEX", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"APEXUSDT": {
			"APEX", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"APTUSDC": {
			"APT", "USDC",
			Precision {
				0.01, 0.0001,
			},
		},
		"APTUSDT": {
			"APT", "USDT",
			Precision {
				0.01, 0.0001,
			},
		},
		"ARBUSDC": {
			"ARB", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"ARBUSDT": {
			"ARB", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"ATOMUSDC": {
			"ATOM", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"ATOMUSDT": {
			"ATOM", "USDT",
			Precision {
				0.001, 0.000001,
			},
		},
		"AVAXEUR": {
			"AVAX", "EUR",
			Precision {
				0.01, 0.0001,
			},
		},
		"AVAXUSDC": {
			"AVAX", "USDC",
			Precision {
				0.001, 0.00001,
			},
		},
		"AVAXUSDT": {
			"AVAX", "USDT",
			Precision {
				0.001, 0.00001,
			},
		},
		"BBSOLSOL": {
			"BBSOL", "SOL",
			Precision {
				0.001, 1e-7,
			},
		},
		"BBSOLUSDC": {
			"BBSOL", "USDC",
			Precision {
				0.001, 0.00001,
			},
		},
		"BBSOLUSDT": {
			"BBSOL", "USDT",
			Precision {
				0.001, 0.00001,
			},
		},
		"BCHUSDC": {
			"BCH", "USDC",
			Precision {
				0.001, 0.0001,
			},
		},
		"BCHUSDT": {
			"BCH", "USDT",
			Precision {
				0.001, 0.0001,
			},
		},
		"BNBUSDC": {
			"BNB", "USDC",
			Precision {
				0.001, 0.0001,
			},
		},
		"BNBUSDT": {
			"BNB", "USDT",
			Precision {
				0.00001, 1e-7,
			},
		},
		"BONKUSDC": {
			"BONK", "USDC",
			Precision {
				1, 1e-8,
			},
		},
		"BONKUSDT": {
			"BONK", "USDT",
			Precision {
				0.1, 1e-9,
			},
		},
		"BRETTUSDC": {
			"BRETT", "USDC",
			Precision {
				1, 0.00001,
			},
		},
		"BRETTUSDT": {
			"BRETT", "USDT",
			Precision {
				0.01, 1e-7,
			},
		},
		"BTCBRL": {
			"BTC", "BRL",
			Precision {
				0.00001, 0.00001,
			},
		},
		"BTCBRZ": {
			"BTC", "BRZ",
			Precision {
				0.00001, 0.000001,
			},
		},
		"BTCDAI": {
			"BTC", "DAI",
			Precision {
				0.000001, 1e-8,
			},
		},
		"BTCEUR": {
			"BTC", "EUR",
			Precision {
				0.000001, 1e-8,
			},
		},
		"BTCPLN": {
			"BTC", "PLN",
			Precision {
				0.00001, 0.00001,
			},
		},
		"BTCTRY": {
			"BTC", "TRY",
			Precision {
				0.00001, 0.00001,
			},
		},
		"BTCUSDC": {
			"BTC", "USDC",
			Precision {
				0.000001, 1e-8,
			},
		},
		"BTCUSDE": {
			"BTC", "USDE",
			Precision {
				0.000001, 1e-8,
			},
		},
		"BTCUSDT": {
			"BTC", "USDT",
			Precision {
				0.000001, 1e-8,
			},
		},
		"CATIEUR": {
			"CATI", "EUR",
			Precision {
				0.01, 0.000001,
			},
		},
		"CATIUSDC": {
			"CATI", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"CATIUSDT": {
			"CATI", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"CHZUSDC": {
			"CHZ", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"CHZUSDT": {
			"CHZ", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"DAIUSDT": {
			"DAI", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"DOGEEUR": {
			"DOGE", "EUR",
			Precision {
				0.01, 1e-7,
			},
		},
		"DOGEUSDC": {
			"DOGE", "USDC",
			Precision {
				0.1, 0.000001,
			},
		},
		"DOGEUSDT": {
			"DOGE", "USDT",
			Precision {
				0.1, 0.000001,
			},
		},
		"DOGSEUR": {
			"DOGS", "EUR",
			Precision {
				1, 1e-7,
			},
		},
		"DOGSUSDC": {
			"DOGS", "USDC",
			Precision {
				1, 1e-7,
			},
		},
		"DOGSUSDT": {
			"DOGS", "USDT",
			Precision {
				1, 1e-7,
			},
		},
		"DOTBTC": {
			"DOT", "BTC",
			Precision {
				0.01, 1e-10,
			},
		},
		"DOTUSDC": {
			"DOT", "USDC",
			Precision {
				0.001, 0.000001,
			},
		},
		"DOTUSDT": {
			"DOT", "USDT",
			Precision {
				0.001, 0.000001,
			},
		},
		"ENAEUR": {
			"ENA", "EUR",
			Precision {
				1, 0.0001,
			},
		},
		"ENAUSDT": {
			"ENA", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"EOSUSDC": {
			"EOS", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"EOSUSDT": {
			"EOS", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"ETHBRL": {
			"ETH", "BRL",
			Precision {
				0.0001, 0.000001,
			},
		},
		"ETHBTC": {
			"ETH", "BTC",
			Precision {
				0.00001, 1e-11,
			},
		},
		"ETHDAI": {
			"ETH", "DAI",
			Precision {
				0.00001, 1e-7,
			},
		},
		"ETHEUR": {
			"ETH", "EUR",
			Precision {
				0.00001, 1e-7,
			},
		},
		"ETHPLN": {
			"ETH", "PLN",
			Precision {
				0.0001, 0.0001,
			},
		},
		"ETHTRY": {
			"ETH", "TRY",
			Precision {
				0.0001, 0.0001,
			},
		},
		"ETHUSDC": {
			"ETH", "USDC",
			Precision {
				0.00001, 1e-7,
			},
		},
		"ETHUSDE": {
			"ETH", "USDE",
			Precision {
				0.00001, 1e-7,
			},
		},
		"ETHUSDT": {
			"ETH", "USDT",
			Precision {
				0.00001, 1e-7,
			},
		},
		"FETUSDC": {
			"FET", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"FETUSDT": {
			"FET", "USDT",
			Precision {
				0.01, 0.00001,
			},
		},
		"FILUSDC": {
			"FIL", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"FILUSDT": {
			"FIL", "USDT",
			Precision {
				0.01, 0.00001,
			},
		},
		"FLOKIUSDC": {
			"FLOKI", "USDC",
			Precision {
				1, 1e-8,
			},
		},
		"FLOKIUSDT": {
			"FLOKI", "USDT",
			Precision {
				0.1, 1e-9,
			},
		},
		"GMTUSDC": {
			"GMT", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"GMTUSDT": {
			"GMT", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"HFTUSDC": {
			"HFT", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"HFTUSDT": {
			"HFT", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"HMSTRUSDC": {
			"HMSTR", "USDC",
			Precision {
				0.01, 1e-8,
			},
		},
		"HMSTRUSDT": {
			"HMSTR", "USDT",
			Precision {
				0.01, 1e-8,
			},
		},
		"ICPUSDC": {
			"ICP", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"ICPUSDT": {
			"ICP", "USDT",
			Precision {
				0.01, 0.00001,
			},
		},
		"INJUSDC": {
			"INJ", "USDC",
			Precision {
				0.01, 0.0001,
			},
		},
		"INJUSDT": {
			"INJ", "USDT",
			Precision {
				0.01, 0.0001,
			},
		},
		"JASMYUSDC": {
			"JASMY", "USDC",
			Precision {
				0.01, 1e-7,
			},
		},
		"JASMYUSDT": {
			"JASMY", "USDT",
			Precision {
				0.01, 1e-7,
			},
		},
		"KASUSDC": {
			"KAS", "USDC",
			Precision {
				0.01, 1e-7,
			},
		},
		"KASUSDT": {
			"KAS", "USDT",
			Precision {
				0.01, 1e-7,
			},
		},
		"LDOUSDC": {
			"LDO", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"LDOUSDT": {
			"LDO", "USDT",
			Precision {
				0.01, 0.00001,
			},
		},
		"LINKEUR": {
			"LINK", "EUR",
			Precision {
				0.01, 0.0001,
			},
		},
		"LINKUSDC": {
			"LINK", "USDC",
			Precision {
				0.001, 0.00001,
			},
		},
		"LINKUSDT": {
			"LINK", "USDT",
			Precision {
				0.001, 0.00001,
			},
		},
		"LTCBTC": {
			"LTC", "BTC",
			Precision {
				0.01, 1e-8,
			},
		},
		"LTCEUR": {
			"LTC", "EUR",
			Precision {
				0.00001, 1e-7,
			},
		},
		"LTCUSDC": {
			"LTC", "USDC",
			Precision {
				0.00001, 1e-7,
			},
		},
		"LTCUSDT": {
			"LTC", "USDT",
			Precision {
				0.00001, 1e-7,
			},
		},
		"LUNCUSDC": {
			"LUNC", "USDC",
			Precision {
				0.001, 1e-11,
			},
		},
		"LUNCUSDT": {
			"LUNC", "USDT",
			Precision {
				0.001, 1e-11,
			},
		},
		"MANABTC": {
			"MANA", "BTC",
			Precision {
				0.1, 1e-9,
			},
		},
		"MANAUSDC": {
			"MANA", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"MANAUSDT": {
			"MANA", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"METHETH": {
			"METH", "ETH",
			Precision {
				0.00001, 1e-9,
			},
		},
		"METHUSDT": {
			"METH", "USDT",
			Precision {
				0.00001, 1e-7,
			},
		},
		"MEWUSDC": {
			"MEW", "USDC",
			Precision {
				0.01, 1e-8,
			},
		},
		"MEWUSDT": {
			"MEW", "USDT",
			Precision {
				0.01, 1e-8,
			},
		},
		"MNTBTC": {
			"MNT", "BTC",
			Precision {
				0.01, 1e-10,
			},
		},
		"MNTUSDC": {
			"MNT", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"MNTUSDT": {
			"MNT", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"NEAREUR": {
			"NEAR", "EUR",
			Precision {
				0.1, 0.0001,
			},
		},
		"NEARUSDC": {
			"NEAR", "USDC",
			Precision {
				0.1, 0.0001,
			},
		},
		"NEARUSDT": {
			"NEAR", "USDT",
			Precision {
				0.01, 0.00001,
			},
		},
		"NOTUSDC": {
			"NOT", "USDC",
			Precision {
				1, 0.000001,
			},
		},
		"NOTUSDT": {
			"NOT", "USDT",
			Precision {
				0.01, 1e-8,
			},
		},
		"ONDOEUR": {
			"ONDO", "EUR",
			Precision {
				0.01, 0.000001,
			},
		},
		"ONDOUSDC": {
			"ONDO", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"ONDOUSDT": {
			"ONDO", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"OPUSDC": {
			"OP", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"OPUSDT": {
			"OP", "USDT",
			Precision {
				0.01, 0.00001,
			},
		},
		"PEPEEUR": {
			"PEPE", "EUR",
			Precision {
				1, 1e-8,
			},
		},
		"PEPEUSDC": {
			"PEPE", "USDC",
			Precision {
				1, 1e-8,
			},
		},
		"PEPEUSDT": {
			"PEPE", "USDT",
			Precision {
				1, 1e-8,
			},
		},
		"SANDBTC": {
			"SAND", "BTC",
			Precision {
				0.1, 1e-9,
			},
		},
		"SANDUSDC": {
			"SAND", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"SANDUSDT": {
			"SAND", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"SEIUSDC": {
			"SEI", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"SEIUSDT": {
			"SEI", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"SHIBEUR": {
			"SHIB", "EUR",
			Precision {
				0.1, 1e-9,
			},
		},
		"SHIBUSDC": {
			"SHIB", "USDC",
			Precision {
				0.1, 1e-9,
			},
		},
		"SHIBUSDT": {
			"SHIB", "USDT",
			Precision {
				0.1, 1e-9,
			},
		},
		"SOLBRL": {
			"SOL", "BRL",
			Precision {
				0.001, 0.0001,
			},
		},
		"SOLBTC": {
			"SOL", "BTC",
			Precision {
				0.01, 1e-9,
			},
		},
		"SOLEUR": {
			"SOL", "EUR",
			Precision {
				0.001, 0.00001,
			},
		},
		"SOLUSDC": {
			"SOL", "USDC",
			Precision {
				0.001, 0.00001,
			},
		},
		"SOLUSDE": {
			"SOL", "USDE",
			Precision {
				0.0001, 0.000001,
			},
		},
		"SOLUSDT": {
			"SOL", "USDT",
			Precision {
				0.001, 0.00001,
			},
		},
		"STETHEUR": {
			"STETH", "EUR",
			Precision {
				0.0001, 0.000001,
			},
		},
		"STETHUSDT": {
			"STETH", "USDT",
			Precision {
				0.00001, 1e-7,
			},
		},
		"STRKUSDC": {
			"STRK", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"STRKUSDT": {
			"STRK", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"SUIUSDC": {
			"SUI", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"SUIUSDT": {
			"SUI", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"SWELLUSDC": {
			"SWELL", "USDC",
			Precision {
				1, 0.00001,
			},
		},
		"SWELLUSDT": {
			"SWELL", "USDT",
			Precision {
				1, 0.00001,
			},
		},
		"TIAUSDC": {
			"TIA", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"TIAUSDT": {
			"TIA", "USDT",
			Precision {
				0.01, 0.00001,
			},
		},
		"TONEUR": {
			"TON", "EUR",
			Precision {
				0.01, 0.00001,
			},
		},
		"TONUSDC": {
			"TON", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"TONUSDT": {
			"TON", "USDT",
			Precision {
				0.01, 0.00001,
			},
		},
		"TRUMPUSDC": {
			"TRUMP", "USDC",
			Precision {
				0.01, 0.0001,
			},
		},
		"TRUMPUSDT": {
			"TRUMP", "USDT",
			Precision {
				0.01, 0.0001,
			},
		},
		"TRXUSDC": {
			"TRX", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"TRXUSDT": {
			"TRX", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"UNIUSDC": {
			"UNI", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"UNIUSDT": {
			"UNI", "USDT",
			Precision {
				0.001, 0.000001,
			},
		},
		"USDCBRL": {
			"USDC", "BRL",
			Precision {
				0.1, 0.0001,
			},
		},
		"USDCEUR": {
			"USDC", "EUR",
			Precision {
				0.01, 0.000001,
			},
		},
		"USDCUSDT": {
			"USDC", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"USDEUSDC": {
			"USDE", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"USDEUSDT": {
			"USDE", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"USDTBRL": {
			"USDT", "BRL",
			Precision {
				0.1, 0.0001,
			},
		},
		"USDTBRZ": {
			"USDT", "BRZ",
			Precision {
				0.01, 0.00001,
			},
		},
		"USDTEUR": {
			"USDT", "EUR",
			Precision {
				0.01, 0.000001,
			},
		},
		"USDTPLN": {
			"USDT", "PLN",
			Precision {
				1, 0.001,
			},
		},
		"USDTTRY": {
			"USDT", "TRY",
			Precision {
				1, 0.01,
			},
		},
		"WBTCBTC": {
			"WBTC", "BTC",
			Precision {
				0.0001, 1e-8,
			},
		},
		"WBTCUSDT": {
			"WBTC", "USDT",
			Precision {
				0.000001, 1e-8,
			},
		},
		"WIFEUR": {
			"WIF", "EUR",
			Precision {
				0.1, 0.0001,
			},
		},
		"WIFUSDC": {
			"WIF", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"WIFUSDT": {
			"WIF", "USDT",
			Precision {
				0.01, 0.00001,
			},
		},
		"WLDEUR": {
			"WLD", "EUR",
			Precision {
				0.01, 0.00001,
			},
		},
		"WLDUSDC": {
			"WLD", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"WLDUSDT": {
			"WLD", "USDT",
			Precision {
				0.01, 0.00001,
			},
		},
		"XLMBTC": {
			"XLM", "BTC",
			Precision {
				0.1, 1e-11,
			},
		},
		"XLMUSDC": {
			"XLM", "USDC",
			Precision {
				0.1, 0.00001,
			},
		},
		"XLMUSDT": {
			"XLM", "USDT",
			Precision {
				0.1, 0.00001,
			},
		},
		"XRPBTC": {
			"XRP", "BTC",
			Precision {
				0.1, 1e-9,
			},
		},
		"XRPEUR": {
			"XRP", "EUR",
			Precision {
				0.01, 0.000001,
			},
		},
		"XRPUSDC": {
			"XRP", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"XRPUSDT": {
			"XRP", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"ZKUSDC": {
			"ZK", "USDC",
			Precision {
				0.01, 0.000001,
			},
		},
		"ZKUSDT": {
			"ZK", "USDT",
			Precision {
				0.01, 0.000001,
			},
		},
		"ZROUSDC": {
			"ZRO", "USDC",
			Precision {
				0.01, 0.00001,
			},
		},
		"ZROUSDT": {
			"ZRO", "USDT",
			Precision {
				0.01, 0.00001,
			},
		},
	}
)
