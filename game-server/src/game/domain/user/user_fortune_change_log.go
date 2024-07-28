package user

import (
	"bytes"
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
)

//channel :
//1：新用户， 2：充值， 3：后台赠送，4：兑换码，5：破产补助，
//6：保管赠送抽水，8：初级场抽水，9高级场抽水，10钻石充值，11钻石赠送，12钻石开房， 13:钻石兑换码,
//7：购买豪车，14：购买头像框，15：购买vip，16：完善资料, 17新用户钻石， 18绑定代理赠送， 19开房退费, 20新手场抽水,
//21:签到1,22：签到2,23：签到3,24：签到4,25签到5,26，签到6,27签到7， 29分享， 30：私人场万一局， 31， 加一个俱乐部

type UserfortuneChangeLogReq struct {
	Chl   int `json:"chl"`
	Count int `json:"count"`
}

func UserfortuneChangeLog(channel, count int) {
	req := UserfortuneChangeLogReq{}
	req.Chl = channel
	req.Count = count

	b, err := json.Marshal(req)
	if err != nil {
		glog.Info("User_fortune_Change_Log err:", err)
		return
	}

	body := bytes.NewBuffer([]byte(b))
	res, err := http.Post("http://127.0.0.1:8089/log/fortune_log", "application/json;charset=utf-8", body)
	if err != nil {
		glog.Info("User_fortune_Change_Log err:", err)
		return
	}

	ioutil.ReadAll(res.Body)
	res.Body.Close()
	return
}
