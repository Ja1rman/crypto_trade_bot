package stats

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"crypto_trading/src/logger"
)

var (
	RecieveFromServerDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:    "recieve_from_server_duration",
		Help:    "Время принятия сообщения от сервера в миллисекундах",
	}, []string{"duration"}) 
	AnalyzeDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:    "analyze_duration",
		Help:    "Время выполнения анализа в микросекундах",
	}, []string{"duration"}) 
	OrdersDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:    "orders_duration",
		Help:    "Время выполнения покупки/продажи ордеров в миллисекундах",
	}, []string{"duration"}) 
	WalletBalance = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wallet_balance",
		Help: "Баланс кошелька на бирже по каждой валюте",
	}, []string{"currency"}) 
)


func PushToPrometheus(gaugeVec *prometheus.GaugeVec, job string, value float64, label string) {
	gaugeVec.WithLabelValues(label).Set(value)

	reg := prometheus.NewRegistry()
	if err := reg.Register(gaugeVec); err != nil {
		logger.Logger.Println("Ошибка регистрации в prometheus:", err)
		return
	}

	if err := push.New("http://pushgateway:9091", job).
		Gatherer(reg).
		Grouping("instance", "server").
		Add(); err != nil {
		logger.Logger.Println("Could not push metrics:", err)
		return
	}
}
