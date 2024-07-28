package config

import (
	"encoding/json"
	"game/domain/activity"
	"game/server"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"game/pb"
)
const (
	web_token = "threewebtoken123"
)
type ReqChangeActivity struct {
	Token  string            `json:"token"`
	Data   []pb.DZActivity   `json:"data"`
}

type ResChangeActivity struct {
	Result  int     `json:"result"`
	Desc    string  `json:"desc"`
}

func sendMsg( code int,msg string) []byte {
	Data := &ResChangeActivity{}
	Data.Result = code
	Data.Desc   = msg
	result ,_  := json.Marshal(Data)
	return result
}
/**
 * 设置活动
 */
func Change_activity(w http.ResponseWriter, r *http.Request)  {
	raw,err := ioutil.ReadAll(r.Body)
	if err!= nil {
		 w.Write( sendMsg(1,err.Error()) )
		return
	}
	data := &ReqChangeActivity{}
	err2 :=json.Unmarshal(raw,data)
	if err2!=nil{
		w.Write(sendMsg(2,err2.Error()))
		return
	}
	if data.Token != web_token {
		w.Write(sendMsg(3,"Token Error"))
		return
	}

   model := activity.ConfigActivity{}

   err1  := model.UpdateAll(data.Data)
   if err1 != nil{
   	   w.Write(sendMsg(4,err1.Error()))
	   return
   }
   w.Write(sendMsg(0,"成功了"))
}

func Activity_list(m *server.ClientMsg, sess *server.Session) []byte {
	list,err := activity.ConfigActivity{}.GetAll()
	if err!=nil{
		glog.Error(err.Error())
	}
	return server.BuildClientMsg(m.MsgId,list)

}
