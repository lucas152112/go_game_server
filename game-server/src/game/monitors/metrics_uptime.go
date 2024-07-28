package monitors

import "github.com/prometheus/client_golang/prometheus"

var uptime prometheus.Gauge


func init() {
	uptime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:"game_uptime",
		Help:"game start time ",
	})
	uptime.SetToCurrentTime()
}
