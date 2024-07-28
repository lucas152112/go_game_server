package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"game/monitors"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gorilla/handlers"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"golang.org/x/net/websocket"
)

type MsgDispatcher interface {
	RegisterHandlers(r *mux.Router)
	DispatchMsg(msg *ClientMsg, sess *Session) []byte
}

type GameServer struct {
	dispatcher    MsgDispatcher
	sigChan       chan os.Signal
	waitGroup     *sync.WaitGroup
	stopChan      chan bool
	stopOnce      sync.Once
	refuseService int32
}

var s *GameServer

var bindHost string

func init() {
	flag.StringVar(&bindHost, "bindHost", ":9011", "bind server host.")

	s = &GameServer{}
	s.sigChan = make(chan os.Signal, 1)
	s.waitGroup = &sync.WaitGroup{}
	s.stopChan = make(chan bool)
}

func GetServerInstance() *GameServer {
	return s
}

func (s *GameServer) GetSigChan() chan os.Signal {
	return s.sigChan
}

func (s *GameServer) GetStopChan() chan bool {
	return s.stopChan
}

func (s *GameServer) StartServer(dispatcher MsgDispatcher) {
	s.dispatcher = dispatcher

	r := mux.NewRouter()
	http.Handle("/", r)
	r.Handle("/ws/", websocket.Server{Handler: s.handleClient, Handshake: nil})

	s.dispatcher.RegisterHandlers(r)
	bindHost = viper.GetString("gameport")
	glog.Info("===>启动Game服务", bindHost)

	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()

	glog.Fatal(http.ListenAndServe(fmt.Sprintf("%v", bindHost), handlers.CORS()(r)))
}

func (s *GameServer) StopServer() {
	s.stopOnce.Do(func() {
		go func() {
			s.sigChan <- syscall.SIGKILL
		}()
	})
}

func (s *GameServer) WaitStopServer() {
	glog.Info("==>Start WaitStopServer")
	defer glog.Info("==>WaitStopServer done.")

	close(s.stopChan)
	s.waitGroup.Wait()
}

func (s *GameServer) IsRefuseService() bool {
	return atomic.AddInt32(&s.refuseService, 0) > 0
}

func (s *GameServer) SetRefuseService() {
	atomic.AddInt32(&s.refuseService, 1)
}

func (s *GameServer) handleClient(conn *websocket.Conn) {
	if s.IsRefuseService() {
		glog.Info("==>正在停止服务，拒绝连接...")
		conn.Close()
		return
	}

	sess := newSess(conn)
	go sess.Run(s.dispatcher)

	defer sess.cleanSess()

	for {
		var data []byte
		conn.SetReadDeadline(time.Now().Add(time.Minute * 10))
		err := websocket.Message.Receive(conn, &data)
		if err != nil {
			glog.Info("error receiving msg:", err, ", conn:", conn)
			break
		}

		conn.SetReadDeadline(time.Time{})

		msg := &ClientMsg{}

		errMsg := json.Unmarshal(data, msg)
		if errMsg != nil {
			glog.Error("unmarshal client msg failed!==>", string(data))
			break
		}

		if msg.MsgId != 202 {
			//glog.Info("handleClient:", msg.MsgId, " sess:", conn)
		}

		//监控
		receives, receives_bytes, _, _ := monitors.GetSockets()
		receives.Inc()
		receives_bytes.Add(float64(len(data)))

		sess.mq <- msg
	}

}
