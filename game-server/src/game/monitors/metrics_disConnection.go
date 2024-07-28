package monitors

import "github.com/prometheus/client_golang/prometheus"

var disConnection *prometheus.CounterVec

func init()  {
	disConnection = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "game_disConnection",
		Help: "断线重连次数",
	},[]string{
		"Country",
	})
}

func GetDisConnection() *prometheus.CounterVec  {
	return disConnection
}

/**
 只有新增
 */
func DisCounnectionReport( Country string )  {
	disConnection.With(prometheus.Labels{"Country":Country}).Inc()
}


