package server

import (
	"encoding/json"
	//"game/monitors"
	"game/pb"
	"game/util"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
)

type ClientMsg struct {
	MsgId   int32       `json:"msgId"`
	MsgBody interface{} `json:"msgBody"`
}

type Session struct {
	conn      *websocket.Conn
	IP        string
	mq        chan *ClientMsg
	Data      interface{}
	LoggedIn  bool
	exitChan  chan bool
	cleanOnce sync.Once
	kickOnce  sync.Once
	OnLogout  func()
}

func newSess(conn *websocket.Conn) *Session {
	sess := &Session{}
	sess.conn = conn
	//monitors.GetSessionGauge().Inc()
	return sess
}

func (s *Session) cleanSess() {
	s.cleanOnce.Do(func() {
		glog.V(2).Info("===>清理session:", s)
		//monitors.GetSessionGauge().Dec()
		if s.conn != nil {
			s.conn.Close()
		}

		if s.mq != nil {
			close(s.mq)
		}
	})
}

func (s *Session) GetConn() *websocket.Conn {
	return s.conn
}

func (s *Session) ClearConn() {
	s.conn = nil
	//s.cleanSess() //
}

func (s *Session) Kickout() {
	s.kickOnce.Do(func() {
		glog.Info("Kickout")
		close(s.exitChan)
	})
}

func (s *Session) Run(dispatcher MsgDispatcher) {
	defer util.PrintPanicStack()
	s.IP = s.conn.Request().Header.Get("X-Real-Ip")
	if len(s.IP) == 0 {
		s.IP = strings.Split(s.conn.Request().RemoteAddr, ":")[0]
	}

	s.mq = make(chan *ClientMsg, 100)
	s.exitChan = make(chan bool)

	GetServerInstance().waitGroup.Add(1)

	glog.Info("===>打开session:", s)

	defer func() {
		glog.Info("disconnected:")
		util.PrintPanicStack()
		s.cleanSess()
		s.logout()
		GetServerInstance().waitGroup.Done()
	}()

	for {
		select {
		case <-GetServerInstance().stopChan:
			glog.Info("stopChan ")
			return
		case msg, ok := <-s.mq:
			if !ok {
				glog.Info("<-s.mq != ok")
				return
			}
			res := dispatcher.DispatchMsg(msg, s)
			if res != nil {
				s.SendToClient(res)
			}
		case <-s.exitChan:
			glog.Info("<-s.exitChan")
			return
		}
	}
}

func (s *Session) logout() {
	if s.LoggedIn && s.OnLogout != nil {
		s.OnLogout()

		s.LoggedIn = false
	}
}

func (s *Session) SendMQ(msg *ClientMsg) bool {
	ret := true

	defer func() {
		if r := recover(); r != nil {
			ret = false
		}
	}()
	s.mq <- msg

	return ret
}

func (s *Session) SendToClient(msg []byte) {
	//ret := true

	defer func() {
		if r := recover(); r != nil {
			//ret = false
			clientMsg := &ClientMsg{}
			err := json.Unmarshal(msg, clientMsg)
			if err != nil {
				glog.Error("SendToClient json.Unmarshal err:", err)
			}
			//glog.Info("Get Error ",r)
			//glog.Info("send on closed channel sess:", s.IP,r)
		}
	}()

	if msg != nil {

		//_, _, sends, sends_bytes := monitors.GetSockets()
		//sends.Inc()
		//sends_bytes.Add(float64(len(msg)))
		//

		if true {
			clientMsg := &ClientMsg{}
			err := json.Unmarshal(msg, clientMsg)
			if err != nil {
				glog.Error("====>SendToClient unmarshal failed msgId:", clientMsg.MsgId, ", conn:", s.conn)
			}

			//glog.Info("==>向客户端发送消息 msgId:", clientMsg.MsgId)
		}
		s.conn.SetWriteDeadline(time.Now().Add(time.Second))
		//glog.Info("s.conn",s.conn)
		sendMsg := string(msg)
		err := websocket.Message.Send(s.conn, sendMsg)
		if err != nil {
			glog.Info("===>发送失败err:", err, " s.conn:", s.conn)
			s.cleanSess()
			return
		}

		//glog.Info("===>SendToClient Msg", sendMsg)

		s.conn.SetWriteDeadline(time.Time{})
	}
}

func BuildClientMsg(msgId int32, body interface{}) []byte {
	if body == nil {
		return nil
	}
	msg := &ClientMsg{}
	msg.MsgId = msgId
	d, err := json.Marshal(body)
	if err != nil {
		//glog.Info("json error",err.Error())
		panic(err)
		return nil
	}

	msg.MsgBody = string(d)

	res, errRes := json.Marshal(msg)
	if errRes != nil {
		//glog.Info("json error 2",errRes.Error())
		panic(errRes)
		return nil
	}

	//debug？？
	if msgId != 202 {
		if len(d) > 500 {
			localIP := util.GetLocalIP()
			if localIP == "172.17.79.168" {
				//glog.Info("BuildClientMsg body = ", body, " msgId:", msgId)
			} else {
				//glog.Info("BuildClientMsg msgId:", msgId)
			}
		} else {
			//glog.Info("BuildClientMsg body = ", body, " msgId:", msgId)
		}
	}

	return res
}

// Build msg without double marshal
func BuildClientMsgSimplify(msgId int32, body []byte) []byte {
	if body == nil {
		return nil
	}
	msg := &ClientMsg{}
	msg.MsgId = msgId
	msg.MsgBody = string(body)

	res, errRes := json.Marshal(msg)
	if errRes != nil {
		panic(errRes)
		return nil
	}

	return res
}

func BuildClientMsg3(msgId int32, body []byte) []byte {
	if _, ok := pb.MessageId_name[msgId]; !ok {
		glog.Warning("build client msg failed client msgId:", msgId, " does not exist")
		return nil
	}

	msg := &pb.ClientMsg{}
	msg.MsgId = proto.Int32(int32(msgId))
	msg.MsgBody = body

	res, err := proto.Marshal(msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	//glog.V(1).Info("向客户端发送msgId:", util.GetMsgIdName(msgId), " 长度:", len(res))

	return res
}
