package user

import (
	"bytes"
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
)

type WebPayLogReq struct {
	UserId      string `json:"userId"`
	PayType     string `json:"payType"`
	UserChannel string `json:"userChannel"`
	PayCount    int    `json:"payCount"`
}

func WebPayLog(userId string, payCount int, userChannel string, payType string) {
	req := WebPayLogReq{}
	req.UserId = userId
	req.UserChannel = userChannel
	req.PayType = payType
	req.PayCount = payCount

	b, err := json.Marshal(req)
	if err != nil {
		glog.Info("WebPayLog err:", err)
		return
	}

	body := bytes.NewBuffer([]byte(b))
	res, err := http.Post("http://127.0.0.1:8089/log/pay_log", "application/json;charset=utf-8", body)
	if err != nil {
		glog.Info("WebPayLog err:", err)
		return
	}

	ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}
