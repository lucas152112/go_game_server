package monitors

import (
	"flag"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var port string

func init()  {

	flag.StringVar(&port, "monitorPort", ":9108", "monitors port.")
}

func Server()  {

	reg:= prometheus.NewPedanticRegistry()
	reg.MustRegister(sessionGauge)
	reg.MustRegister(socket_receives)
	reg.MustRegister(socket_receive_bytes)
	reg.MustRegister(socket_sends)
	reg.MustRegister(socket_send_bytes)
	reg.MustRegister(socket_executes)
	reg.MustRegister(socket_execute_times)
	reg.MustRegister(uptime)
	reg.MustRegister(tableGauge)
	reg.MustRegister(tableCount)
	reg.MustRegister(disConnection)      //短信重连报警

	//定义采集器
	gathers:= prometheus.Gatherers{
		prometheus.DefaultGatherer,
		reg,
	}

	handler := promhttp.HandlerFor(gathers,
		promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	})

	glog.Info("==>启动监控 默认 :9108")
	http.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer,request)
	})
	http.ListenAndServe(port,nil)
}
