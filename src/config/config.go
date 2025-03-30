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
	TickSize float64 `json:"TickSize"`
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
				1e2, 1e6,
			},
			1e4,
		},
		"ADAUSDC": {
			"ADA", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"ADAUSDT": {
			"ADA", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"ALGOBTC": {
			"ALGO", "BTC" ,
			Precision {
				1e1, 1e10,
			},
			1e9,
		},
		"ALGOUSDT": {
			"ALGO", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"APEUSDC": {
			"APE", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"APEUSDT": {
			"APE", "USDT" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"APEXUSDC": {
			"APEX", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"APEXUSDT": {
			"APEX", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"APTUSDC": {
			"APT", "USDC" ,
			Precision {
				1e2, 1e4,
			},
			1e2,
		},
		"APTUSDT": {
			"APT", "USDT" ,
			Precision {
				1e2, 1e4,
			},
			1e2,
		},
		"ARBUSDC": {
			"ARB", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"ARBUSDT": {
			"ARB", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"ATOMUSDC": {
			"ATOM", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"ATOMUSDT": {
			"ATOM", "USDT" ,
			Precision {
				1e3, 1e6,
			},
			1e3,
		},
		"AVAXEUR": {
			"AVAX", "EUR" ,
			Precision {
				1e2, 1e4,
			},
			1e2,
		},
		"AVAXUSDC": {
			"AVAX", "USDC" ,
			Precision {
				1e3, 1e5,
			},
			1e2,
		},
		"AVAXUSDT": {
			"AVAX", "USDT" ,
			Precision {
				1e3, 1e5,
			},
			1e2,
		},
		"BBSOLSOL": {
			"BBSOL", "SOL" ,
			Precision {
				1e3, 1e7,
			},
			1e4,
		},
		"BBSOLUSDC": {
			"BBSOL", "USDC" ,
			Precision {
				1e3, 1e5,
			},
			1e2,
		},
		"BBSOLUSDT": {
			"BBSOL", "USDT" ,
			Precision {
				1e3, 1e5,
			},
			1e2,
		},
		"BCHUSDC": {
			"BCH", "USDC" ,
			Precision {
				1e3, 1e4,
			},
			1e1,
		},
		"BCHUSDT": {
			"BCH", "USDT" ,
			Precision {
				1e3, 1e4,
			},
			1e1,
		},
		"BNBUSDC": {
			"BNB", "USDC" ,
			Precision {
				1e4, 1e5,
			},
			1e1,
		},
		"BNBUSDT": {
			"BNB", "USDT" ,
			Precision {
				1e4, 1e6,
			},
			1e2,
		},
		"BONKUSDC": {
			"BONK", "USDC" ,
			Precision {
				1, 1e8,
			},
			1e8,
		},
		"BONKUSDT": {
			"BONK", "USDT" ,
			Precision {
				1e1, 1e9,
			},
			1e8,
		},
		"BRETTUSDC": {
			"BRETT", "USDC" ,
			Precision {
				1, 1e5,
			},
			1e5,
		},
		"BRETTUSDT": {
			"BRETT", "USDT" ,
			Precision {
				1e2, 1e7,
			},
			1e5,
		},
		"BTCBRL": {
			"BTC", "BRL" ,
			Precision {
				1e6, 1e6,
			},
			1,
		},
		"BTCBRZ": {
			"BTC", "BRZ" ,
			Precision {
				1e6, 1e7,
			},
			1e1,
		},
		"BTCDAI": {
			"BTC", "DAI" ,
			Precision {
				1e6, 1e7,
			},
			1e1,
		},
		"BTCEUR": {
			"BTC", "EUR" ,
			Precision {
				1e6, 1e7,
			},
			1e1,
		},
		"BTCPLN": {
			"BTC", "PLN" ,
			Precision {
				1e6, 1e6,
			},
			1,
		},
		"BTCTRY": {
			"BTC", "TRY" ,
			Precision {
				1e6, 1e6,
			},
			1,
		},
		"BTCUSDC": {
			"BTC", "USDC" ,
			Precision {
				1e6, 1e7,
			},
			1e1,
		},
		"BTCUSDE": {
			"BTC", "USDE" ,
			Precision {
				1e6, 1e7,
			},
			1e1,
		},
		"BTCUSDT": {
			"BTC", "USDT" ,
			Precision {
				1e6, 1e7,
			},
			1e1,
		},
		"CATIEUR": {
			"CATI", "EUR" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"CATIUSDC": {
			"CATI", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"CATIUSDT": {
			"CATI", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"CHZUSDC": {
			"CHZ", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"CHZUSDT": {
			"CHZ", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"DAIUSDT": {
			"DAI", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"DOGEEUR": {
			"DOGE", "EUR" ,
			Precision {
				1e2, 1e7,
			},
			1e5,
		},
		"DOGEUSDC": {
			"DOGE", "USDC" ,
			Precision {
				1e1, 1e6,
			},
			1e5,
		},
		"DOGEUSDT": {
			"DOGE", "USDT" ,
			Precision {
				1e1, 1e6,
			},
			1e5,
		},
		"DOGSEUR": {
			"DOGS", "EUR" ,
			Precision {
				1, 1e7,
			},
			1e7,
		},
		"DOGSUSDC": {
			"DOGS", "USDC" ,
			Precision {
				1, 1e7,
			},
			1e7,
		},
		"DOGSUSDT": {
			"DOGS", "USDT" ,
			Precision {
				1, 1e7,
			},
			1e7,
		},
		"DOTBTC": {
			"DOT", "BTC" ,
			Precision {
				1e2, 1e10,
			},
			1e8,
		},
		"DOTUSDC": {
			"DOT", "USDC" ,
			Precision {
				1e3, 1e6,
			},
			1e3,
		},
		"DOTUSDT": {
			"DOT", "USDT" ,
			Precision {
				1e3, 1e6,
			},
			1e3,
		},
		"ENAEUR": {
			"ENA", "EUR" ,
			Precision {
				1, 1e4,
			},
			1e4,
		},
		"ENAUSDT": {
			"ENA", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"EOSUSDC": {
			"EOS", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"EOSUSDT": {
			"EOS", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"ETHBRL": {
			"ETH", "BRL" ,
			Precision {
				1e4, 1e6,
			},
			1e2,
		},
		"ETHBTC": {
			"ETH", "BTC" ,
			Precision {
				1e5, 1e11,
			},
			1e6,
		},
		"ETHDAI": {
			"ETH", "DAI" ,
			Precision {
				1e5, 1e7,
			},
			1e2,
		},
		"ETHEUR": {
			"ETH", "EUR" ,
			Precision {
				1e5, 1e7,
			},
			1e2,
		},
		"ETHPLN": {
			"ETH", "PLN" ,
			Precision {
				1e4, 1e4,
			},
			1,
		},
		"ETHTRY": {
			"ETH", "TRY" ,
			Precision {
				1e4, 1e4,
			},
			1,
		},
		"ETHUSDC": {
			"ETH", "USDC" ,
			Precision {
				1e5, 1e7,
			},
			1e2,
		},
		"ETHUSDE": {
			"ETH", "USDE" ,
			Precision {
				1e5, 1e7,
			},
			1e2,
		},
		"ETHUSDT": {
			"ETH", "USDT" ,
			Precision {
				1e5, 1e7,
			},
			1e2,
		},
		"FETUSDC": {
			"FET", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"FETUSDT": {
			"FET", "USDT" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"FILUSDC": {
			"FIL", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"FILUSDT": {
			"FIL", "USDT" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"FLOKIUSDC": {
			"FLOKI", "USDC" ,
			Precision {
				1, 1e8,
			},
			1e8,
		},
		"FLOKIUSDT": {
			"FLOKI", "USDT" ,
			Precision {
				1e1, 1e9,
			},
			1e8,
		},
		"GMTUSDC": {
			"GMT", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"GMTUSDT": {
			"GMT", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"HFTUSDC": {
			"HFT", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"HFTUSDT": {
			"HFT", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"HMSTRUSDC": {
			"HMSTR", "USDC" ,
			Precision {
				1e2, 1e8,
			},
			1e6,
		},
		"HMSTRUSDT": {
			"HMSTR", "USDT" ,
			Precision {
				1e2, 1e8,
			},
			1e6,
		},
		"ICPUSDC": {
			"ICP", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"ICPUSDT": {
			"ICP", "USDT" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"INJUSDC": {
			"INJ", "USDC" ,
			Precision {
				1e2, 1e4,
			},
			1e2,
		},
		"INJUSDT": {
			"INJ", "USDT" ,
			Precision {
				1e2, 1e4,
			},
			1e2,
		},
		"JASMYUSDC": {
			"JASMY", "USDC" ,
			Precision {
				1e2, 1e7,
			},
			1e5,
		},
		"JASMYUSDT": {
			"JASMY", "USDT" ,
			Precision {
				1e2, 1e7,
			},
			1e5,
		},
		"KASUSDC": {
			"KAS", "USDC" ,
			Precision {
				1e2, 1e7,
			},
			1e5,
		},
		"KASUSDT": {
			"KAS", "USDT" ,
			Precision {
				1e2, 1e7,
			},
			1e5,
		},
		"LDOUSDC": {
			"LDO", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"LDOUSDT": {
			"LDO", "USDT" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"LINKEUR": {
			"LINK", "EUR" ,
			Precision {
				1e2, 1e4,
			},
			1e2,
		},
		"LINKUSDC": {
			"LINK", "USDC" ,
			Precision {
				1e3, 1e5,
			},
			1e2,
		},
		"LINKUSDT": {
			"LINK", "USDT" ,
			Precision {
				1e3, 1e5,
			},
			1e2,
		},
		"LTCBTC": {
			"LTC", "BTC" ,
			Precision {
				1e3, 1e9,
			},
			1e6,
		},
		"LTCEUR": {
			"LTC", "EUR" ,
			Precision {
				1e5, 1e7,
			},
			1e2,
		},
		"LTCUSDC": {
			"LTC", "USDC" ,
			Precision {
				1e5, 1e7,
			},
			1e2,
		},
		"LTCUSDT": {
			"LTC", "USDT" ,
			Precision {
				1e5, 1e7,
			},
			1e2,
		},
		"LUNCUSDC": {
			"LUNC", "USDC" ,
			Precision {
				1e1, 1e9,
			},
			1e8,
		},
		"LUNCUSDT": {
			"LUNC", "USDT" ,
			Precision {
				1e1, 1e9,
			},
			1e8,
		},
		"MANABTC": {
			"MANA", "BTC" ,
			Precision {
				1e1, 1e9,
			},
			1e8,
		},
		"MANAUSDC": {
			"MANA", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"MANAUSDT": {
			"MANA", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"METHETH": {
			"METH", "ETH" ,
			Precision {
				1e5, 1e9,
			},
			1e4,
		},
		"METHUSDT": {
			"METH", "USDT" ,
			Precision {
				1e5, 1e7,
			},
			1e2,
		},
		"MEWUSDC": {
			"MEW", "USDC" ,
			Precision {
				1e2, 1e8,
			},
			1e6,
		},
		"MEWUSDT": {
			"MEW", "USDT" ,
			Precision {
				1e2, 1e8,
			},
			1e6,
		},
		"MNTBTC": {
			"MNT", "BTC" ,
			Precision {
				1e2, 1e10,
			},
			1e8,
		},
		"MNTUSDC": {
			"MNT", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"MNTUSDT": {
			"MNT", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"NEAREUR": {
			"NEAR", "EUR" ,
			Precision {
				1e1, 1e4,
			},
			1e3,
		},
		"NEARUSDC": {
			"NEAR", "USDC" ,
			Precision {
				1e1, 1e4,
			},
			1e3,
		},
		"NEARUSDT": {
			"NEAR", "USDT" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"NOTUSDC": {
			"NOT", "USDC" ,
			Precision {
				1, 1e6,
			},
			1e6,
		},
		"NOTUSDT": {
			"NOT", "USDT" ,
			Precision {
				1e2, 1e8,
			},
			1e6,
		},
		"ONDOEUR": {
			"ONDO", "EUR" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"ONDOUSDC": {
			"ONDO", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"ONDOUSDT": {
			"ONDO", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"OPUSDC": {
			"OP", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"OPUSDT": {
			"OP", "USDT" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"PEPEEUR": {
			"PEPE", "EUR" ,
			Precision {
				1, 1e8,
			},
			1e8,
		},
		"PEPEUSDC": {
			"PEPE", "USDC" ,
			Precision {
				1, 1e8,
			},
			1e8,
		},
		"PEPEUSDT": {
			"PEPE", "USDT" ,
			Precision {
				1, 1e8,
			},
			1e8,
		},
		"SANDBTC": {
			"SAND", "BTC" ,
			Precision {
				1e1, 1e9,
			},
			1e8,
		},
		"SANDUSDC": {
			"SAND", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"SANDUSDT": {
			"SAND", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"SEIUSDC": {
			"SEI", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"SEIUSDT": {
			"SEI", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"SHIBEUR": {
			"SHIB", "EUR" ,
			Precision {
				1e1, 1e9,
			},
			1e8,
		},
		"SHIBUSDC": {
			"SHIB", "USDC" ,
			Precision {
				1e1, 1e9,
			},
			1e8,
		},
		"SHIBUSDT": {
			"SHIB", "USDT" ,
			Precision {
				1e1, 1e9,
			},
			1e8,
		},
		"SOLBRL": {
			"SOL", "BRL" ,
			Precision {
				1e3, 1e4,
			},
			1e1,
		},
		"SOLBTC": {
			"SOL", "BTC" ,
			Precision {
				1e3, 1e10,
			},
			1e7,
		},
		"SOLEUR": {
			"SOL", "EUR" ,
			Precision {
				1e3, 1e5,
			},
			1e2,
		},
		"SOLUSDC": {
			"SOL", "USDC" ,
			Precision {
				1e3, 1e5,
			},
			1e2,
		},
		"SOLUSDE": {
			"SOL", "USDE" ,
			Precision {
				1e4, 1e6,
			},
			1e2,
		},
		"SOLUSDT": {
			"SOL", "USDT" ,
			Precision {
				1e3, 1e5,
			},
			1e2,
		},
		"STETHEUR": {
			"STETH", "EUR" ,
			Precision {
				1e4, 1e6,
			},
			1e2,
		},
		"STETHUSDT": {
			"STETH", "USDT" ,
			Precision {
				1e5, 1e7,
			},
			1e2,
		},
		"STRKUSDC": {
			"STRK", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"STRKUSDT": {
			"STRK", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"SUIUSDC": {
			"SUI", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"SUIUSDT": {
			"SUI", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"SWELLUSDC": {
			"SWELL", "USDC" ,
			Precision {
				1, 1e5,
			},
			1e5,
		},
		"SWELLUSDT": {
			"SWELL", "USDT" ,
			Precision {
				1, 1e5,
			},
			1e5,
		},
		"TIAUSDC": {
			"TIA", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"TIAUSDT": {
			"TIA", "USDT" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"TONEUR": {
			"TON", "EUR" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"TONUSDC": {
			"TON", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"TONUSDT": {
			"TON", "USDT" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"TRUMPUSDC": {
			"TRUMP", "USDC" ,
			Precision {
				1e2, 1e4,
			},
			1e2,
		},
		"TRUMPUSDT": {
			"TRUMP", "USDT" ,
			Precision {
				1e2, 1e4,
			},
			1e2,
		},
		"TRXUSDC": {
			"TRX", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"TRXUSDT": {
			"TRX", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"UNIUSDC": {
			"UNI", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"UNIUSDT": {
			"UNI", "USDT" ,
			Precision {
				1e3, 1e6,
			},
			1e3,
		},
		"USDCBRL": {
			"USDC", "BRL" ,
			Precision {
				1e1, 1e4,
			},
			1e3,
		},
		"USDCEUR": {
			"USDC", "EUR" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"USDCUSDT": {
			"USDC", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"USDEUSDC": {
			"USDE", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"USDEUSDT": {
			"USDE", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"USDTBRL": {
			"USDT", "BRL" ,
			Precision {
				1e1, 1e4,
			},
			1e3,
		},
		"USDTBRZ": {
			"USDT", "BRZ" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"USDTEUR": {
			"USDT", "EUR" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"USDTPLN": {
			"USDT", "PLN" ,
			Precision {
				1e1, 1e4,
			},
			1e3,
		},
		"USDTTRY": {
			"USDT", "TRY" ,
			Precision {
				1e1, 1e3,
			},
			1e2,
		},
		"WBTCBTC": {
			"WBTC", "BTC" ,
			Precision {
				1e6, 1e10,
			},
			1e4,
		},
		"WBTCUSDT": {
			"WBTC", "USDT" ,
			Precision {
				1e6, 1e8,
			},
			1e2,
		},
		"WIFEUR": {
			"WIF", "EUR" ,
			Precision {
				1e1, 1e4,
			},
			1e3,
		},
		"WIFUSDC": {
			"WIF", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"WIFUSDT": {
			"WIF", "USDT" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"WLDEUR": {
			"WLD", "EUR" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"WLDUSDC": {
			"WLD", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"WLDUSDT": {
			"WLD", "USDT" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"XLMBTC": {
			"XLM", "BTC" ,
			Precision {
				1e1, 1e11,
			},
			1e10,
		},
		"XLMUSDC": {
			"XLM", "USDC" ,
			Precision {
				1e1, 1e5,
			},
			1e4,
		},
		"XLMUSDT": {
			"XLM", "USDT" ,
			Precision {
				1e1, 1e5,
			},
			1e4,
		},
		"XRPBTC": {
			"XRP", "BTC" ,
			Precision {
				1e1, 1e9,
			},
			1e8,
		},
		"XRPEUR": {
			"XRP", "EUR" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"XRPUSDC": {
			"XRP", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"XRPUSDT": {
			"XRP", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"ZKUSDC": {
			"ZK", "USDC" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"ZKUSDT": {
			"ZK", "USDT" ,
			Precision {
				1e2, 1e6,
			},
			1e4,
		},
		"ZROUSDC": {
			"ZRO", "USDC" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
		"ZROUSDT": {
			"ZRO", "USDT" ,
			Precision {
				1e2, 1e5,
			},
			1e3,
		},
	}
)
