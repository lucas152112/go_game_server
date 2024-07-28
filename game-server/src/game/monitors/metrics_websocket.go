package monitors

import "github.com/prometheus/client_golang/prometheus"

var (
	socket_receives prometheus.Counter
	socket_receive_bytes prometheus.Counter
	
	socket_sends prometheus.Counter
	socket_send_bytes prometheus.Counter

	socket_executes     prometheus.Counter
	socket_execute_times prometheus.Counter
)

func init()  {
	socket_receives = prometheus.NewCounter(prometheus.CounterOpts{
		Name:"game_websocket_receives",
		Help:"websocket receives request num",
	})
	socket_receive_bytes = prometheus.NewCounter(prometheus.CounterOpts{
		Name:"game_websocket_receive_bytes",
		Help:"websocket receive ",
	})
	socket_sends = prometheus.NewCounter(prometheus.CounterOpts{
		Name:"game_websocket_sends",
		Help:"websocket sends",
	})
	socket_send_bytes = prometheus.NewCounter(prometheus.CounterOpts{
		Name:"game_websocket_send_bytes",
		Help:"websocket send bytes",
	})

	socket_executes = prometheus.NewCounter(prometheus.CounterOpts{
		Name:"game_websocket_executes",
		Help:"执行的次数 number",
	})

	socket_execute_times = prometheus.NewCounter(prometheus.CounterOpts{
		Name:"game_websocket_execute_times",
		Help:"执行的时间 毫秒",
	})
}

func GetSockets() (prometheus.Counter,prometheus.Counter,prometheus.Counter,prometheus.Counter)  {
	return socket_receives,socket_receive_bytes,socket_sends,socket_send_bytes
}

func GetSocketsExecutes()(prometheus.Counter,prometheus.Counter)  {
	return socket_executes,socket_execute_times
}



