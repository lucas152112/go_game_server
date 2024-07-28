package user

import (
	"encoding/json"
	domainUser "game/domain/user"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"game/pb"
)

type stClientPageReq struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}

func ClientPageHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error("ClientPageHandler ReadAll err:", err)
		return
	}

	req := &stClientPageReq{}
	err = json.Unmarshal(b, req)
	if err != nil {
		glog.Error("ClientPageHandler Unmarshal err:", err)
		return
	}
	glog.Info("ClientPageHandler, req:", req)

	userId := req.UserId
	_, err = domainUser.FindByUserId(userId)
	if err != nil {
		glog.Error("ClientPageHandler no find user")
		return
	}

	notify := &pb.DZClientPageNotify{}
	domainUser.GetPlayerManager().SendClientMsgNew(userId, int32(pb.MessageId_DZ_CLIENT_PAGE_NOTIFY), notify)
	glog.Info("ClientPageHandler, success")
	return
}
