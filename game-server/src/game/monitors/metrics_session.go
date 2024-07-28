package monitors

import "github.com/prometheus/client_golang/prometheus"

var sessionGauge prometheus.Gauge

func init()  {
	sessionGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:"game_sessions",
		Help:"Current sessions",
	})
}

func GetSessionGauge() prometheus.Gauge  {
	return sessionGauge
}