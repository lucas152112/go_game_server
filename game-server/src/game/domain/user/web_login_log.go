package user

import (
	"bytes"
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
)

type WebLoginLogReq struct {
	UserId      string `json:"userId"`
	UserChannel string `json:"userChannel"`
}

func WebLoginLog(userId string, userChannel string) {
	req := WebLoginLogReq{}
	req.UserId = userId
	req.UserChannel = userChannel

	b, err := json.Marshal(req)
	if err != nil {
		glog.Info("WebLoginLog err:", err)
		return
	}

	body := bytes.NewBuffer([]byte(b))
	res, err := http.Post("http://127.0.0.1:8089/log/login_log", "application/json;charset=utf-8", body)
	if err != nil {
		glog.Info("WebLoginLog err:", err)
		return
	}

	ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}
