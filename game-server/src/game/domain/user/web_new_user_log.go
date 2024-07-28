package user

import (
	"bytes"
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
)

type NewUserLogReq struct {
	UserId string `json:"userId"`
	Chl    string `json:"chl"`
}

func WebNewUserLog(userId string, channel string) {
	req := NewUserLogReq{}
	req.UserId = userId
	req.Chl = channel

	b, err := json.Marshal(req)
	if err != nil {
		glog.Info("WebNewUserLog err:", err)
		return
	}

	body := bytes.NewBuffer([]byte(b))
	res, err := http.Post("http://127.0.0.1:8089/log/new_user", "application/json;charset=utf-8", body)
	if err != nil {
		glog.Info("WebNewUserLog err:", err)
		return
	}

	ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}
