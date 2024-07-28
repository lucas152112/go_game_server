package user

import (
	"encoding/json"
	"fmt"
	domainDZ "game/domain/dzgame"
	domainUser "game/domain/user"
	"game/pb"
	"game/server"
	"game/util"
	"math/rand"
	"time"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
)

const (
	K_Cash_Channel      = "101"
	Bi_Kuai_Bao_Channel = "102"
	TP_User_Channel     = "tp"
)

func NLoginHandler(m *server.ClientMsg, sess *server.Session) []byte {
	glog.Info("login ", m)

	res := &pb.NLoginAck{}
	res.Code = 1
	res.Reason = "success"

	if domainDZ.GetGameManager().IsStop() {
		res.Code = 2
		res.Reason = "server stop"
		return server.BuildClientMsg(m.MsgId, res)
	}

	reqStr := m.MsgBody.(string)
	req := &pb.NLoginReq{}
	err := json.Unmarshal([]byte(reqStr), req)
	glog.Info("NLoginHandler req:", req)
	if err != nil {
		res.Code = 2
		res.Reason = "protocal error "
		return server.BuildClientMsg(m.MsgId, res)
	}

	userName := req.Username
	password := req.Password
	imei := req.Imei
	model := req.Model
	channel := req.Channel
	version := req.Version
	channelInt := req.ChannelInt
	ip := sess.IP
	jingDu := req.JingDu
	WeiDu := req.WeiDu
	language := req.Language
	Lan := util.L().To(language)

	glog.Info("===>用户登录消息sess:", sess)

	u := &domainUser.User{}
	nickName := model

	if rand.Float64() < 0.5 {
		u.Gender = 1
		u.PhotoUrl = fmt.Sprintf("%v", rand.Int()%4)
	} else {
		u.Gender = 2
		u.PhotoUrl = fmt.Sprintf("%v", 4+rand.Int()%3)
	}

	tempID := time.Now().Unix() + int64(rand.Int()%1000)

	u.UserId = fmt.Sprintf("%v", tempID)
	u.UserName = userName
	if password == "" {
		password = imei
	}
	u.Password = genPassword(password)
	u.Nickname = nickName
	u.CreateTime = util.GetDbTime()
	u.ChannelId = channel
	glog.Info("u.ChannelId:", u.ChannelId)
	u.DeviceModel = model
	u.VersionName = version
	u.ChannelInt = channelInt
	u.Imei = imei
	u.Balance = 0
	u.Coins = 0
	u.IsGuest = true
	u.BGuest = true

	userId := u.UserId

	player := domainUser.NewPlayer()

	player.User = u
	player.SessKey = bson.NewObjectId().Hex()

	domainUser.GetPlayerManager().AddUserItem(userId, sess)

	player.NewPlayer = true

	sess.LoggedIn = true
	sess.OnLogout = player.OnLogout
	sess.Data = player
	player.User.JingDu = jingDu
	player.User.WeiDu = WeiDu
	player.User.IP = ip
	player.LoginIP = sess.IP
	player.LoginDeviceId = model
	player.Language = language

	player.OnLogoutFunc = onLogout
	player.OnExitFunc = onExitFunc

	player.SendToClientFuncNew = func(msgId int32, body interface{}) {
		sess.SendToClient(server.BuildClientMsg(int32(msgId), body))
	}

	// 登录成功
	res.Code = 1
	res.Reason = Lan.Msg("Login success")
	res.UserId = player.User.UserId
	res.ServerTime = time.Now().Unix()

	glog.Info("nLogin success userID:", userId)

	return server.BuildClientMsg(m.MsgId, res)
}

func onLogout(userId string) {
	go domainDZ.GetGameManager().OffLine(userId)
	//退出大厅
	glog.Info("onlogout ==>", userId)
}
func onExitFunc(userId string) {
	//用户超时 下线游戏..
	go domainDZ.GetGameManager().OffLine(userId)
}
