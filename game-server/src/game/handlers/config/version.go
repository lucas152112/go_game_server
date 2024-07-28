package config

import (
	"game/config"
	"encoding/json"
	"game/conf"
	"game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"game/pb"
)

func GetVersionHandler(w http.ResponseWriter, r *http.Request)  {
	glog.Info(config.GetVersion())
	w.Write([]byte(config.GetVersion()))
}

/**
 被动通知
 */

type VersionInfo struct {
	Version       string `json:"client_version"`     //版本信息
	ForcedUpdate  int    `json:"forced_update"`      //是否强制更新
	Desc          string `json:"desc"`               //跟新内容
}

type VersionNotify struct {
	pb.ProBaseResponse
	VersionInfo
}

type NotifyChangeRes struct {
	pb.ProBaseResponse
}


/**
 POST JSON
 */
func NotifyVersionChange( w http.ResponseWriter, r *http.Request) {
	glog.Info("NotifyVersionChange Get it .")
	req := &VersionInfo{}
	res := &NotifyChangeRes{}
	if r.FormValue("token") != "majiangwebtoken123" {
		res.Error("Token Error",1)

		result,_ :=json.Marshal(res)
		w.Write(result)
		return
	}
	raw,err :=ioutil.ReadAll(r.Body)
	if err!=nil{
		res.Error("Body Error ",2)

		result,_ :=json.Marshal(res)
		w.Write(result)
		return
	}
	glog.Info("NotifyVersionChange Req:",string(raw))

	err = json.Unmarshal(raw,&req)
	if err!= nil{
		res.Error("json parse ",2)

		result,_ :=json.Marshal(res)
		w.Write(result)
		return
	}

	if req.Version ==""{
		res.Error("Version None",2)
		result,_ :=json.Marshal(res)
		w.Write(result)
		return
	}

	notify:=&VersionNotify{}
	notify.Success("Success")
	notify.Version = req.Version
	notify.ForcedUpdate = req.ForcedUpdate
	notify.Desc = req.Desc

	//通知
	go user.GetPlayerManager().BroadcastClientMsg(pb.MessageId_Client_Version_Change_Notify,notify)

	res.Success("success")
	result,_ :=json.Marshal(res)
	w.Write(result)
	return
}

func GetNowVersion(m *server.ClientMsg, sess *server.Session) []byte {
	res:= &VersionNotify{}
	//需要配置生产环境...
	if conf.Get().VersionManage.Addr == ""{
		res.Error("Notify Server Not Config",4)
		return server.BuildClientMsg(m.MsgId,res)
	}
	url := conf.Get().VersionManage.Addr
	resp,err := http.Get(url)
	glog.Info("GET New Version :",url)
	if err!=nil{
		res.Error("Http Request error",1)
		return server.BuildClientMsg(m.MsgId,res)
	}
	raw,err:= ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	glog.Info("Get New Version Done;",string(raw))
	if err!= nil{
		res.Error("",2)
		return server.BuildClientMsg(m.MsgId,res)
	}
	body:=&VersionInfo{}
	err =json.Unmarshal(raw,&body)
	if err!=nil{
		res.Error("json parese error",3)
		return server.BuildClientMsg(m.MsgId,res)
	}
	res.Version = body.Version
	res.ForcedUpdate = body.ForcedUpdate
	res.Desc = body.Desc
	res.Success("Success")
	return server.BuildClientMsg(m.MsgId,res)
}
