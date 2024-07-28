package user

import (
	"fmt"
	domainCdKey "game/domain/cdkey"
	"game/domain/core"
	hall "game/domain/hall"
	domainPrize "game/domain/prize"
	"game/pb"
	"game/util"
	"math/rand"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type GamePlayer struct {
	User                *User
	NewPlayer           bool
	SendToClientFuncNew func(msgId int32, body interface{})
	SendToClientFunc    func(msgId int32, body proto.Message)
	MatchRecord         *MatchRecord
	DZMatchRecord       *DZMatchRecord
	LastGameId          int
	LoginTime           time.Time
	//UserLog             *UserLog
	LastChatTime    time.Time
	CDKeyGainRecord *domainCdKey.CDKeyGainRecord
	LoginIP         string
	LoginDeviceId   string
	SessKey         string
	LastExp         int
	OnLogoutFunc    func(userId string)
	OnExitFunc      func(userId string)
	SignInRecord    *domainPrize.SignInRecord
	UserTasks       *domainPrize.UserTasks
	LastActiveTime  time.Time // 最后活动时间
	Language        int
	CurrentClubHall int //当前所在的大厅
	exitChain       chan int
	stopLoop        bool
	isStarting      bool
}

func (this *GamePlayer) loop() {
	tag := "GamePlayer Loop Stop"
	var step int64 = 15
	for {
		//停止信号
		if this.stopLoop {
			glog.Info(tag, " With stopLoop  ==> ", this.User.UserId)
			return
		}

		//活动时间 15秒
		if time.Now().Unix()-this.LastActiveTime.Unix() > step {
			if this != nil && this.OnExitFunc != nil && this.User != nil {
				if this.User.BGuest {
					return
				}
				fmt.Println("gameplayer onexitfunc userID:", this.User.UserId)
				userID := this.User.UserId
				channel := this.User.ChannelId

				ttime, err := time.Parse("2006-01-02 15:04:05", this.User.LastLogin)
				if err == nil {
					tt := ttime.Unix()
					t := time.Now().Unix() - tt
					go hall.AddActionChannelLog(userID, channel, hall.ActionChannelRemainLength, t, false)
				}

				this.OnExitFunc(userID)

				//this.OnExitFunc = nil
			}

			this.Stop()
			glog.Info("gameplayer stop")
			return
		}

		select {
		// 以心跳的粒度 粗粒度
		case <-time.After(2 * time.Second):
			break
		case <-this.exitChain:
			return
		}
	}
}

func NewPlayer() *GamePlayer {
	p := &GamePlayer{}
	p.Start()
	return p
}

func (this *GamePlayer) Start() {
	if this.isStarting == false {
		this.LastActiveTime = time.Now()
		this.exitChain = make(chan int)
		go this.loop()
		this.isStarting = true
		this.stopLoop = false
	}
}

func (this *GamePlayer) ReStart() {
	this.LastActiveTime = time.Now()
	this.exitChain = make(chan int)
	go this.loop()
	this.isStarting = true
	this.stopLoop = false
}

//Stop ...
func (this *GamePlayer) Stop() {
	this.stopLoop = true
}

func GetPlayer(p interface{}) *GamePlayer {
	if p == nil {
		return nil
	}
	switch player := p.(type) {
	case *GamePlayer:
		return player
	}
	return nil
}

var checkTotalCoinsV = []int64{5000, 10000, 20000, 50000, 100000, 200000, 400000, 800000, 1000000, 1500000, 2000000, 3000000, 5000000, 10000000, 20000000, 30000000}

func checkWinCoins(oldCoins, newCoins int64) []int64 {
	out := []int64{}
	for i := 0; i < len(checkTotalCoinsV); i++ {

		if checkTotalCoinsV[i] > newCoins {
			break
		}

		if checkTotalCoinsV[i] >= oldCoins && checkTotalCoinsV[i] <= newCoins {
			out = append(out, checkTotalCoinsV[i])
		}
	}

	return out
}

func checkPlayWinGames(winGames int) int {
	if winGames == 5 {
		return winGames
	} else if winGames == 10 {
		return winGames
	} else if winGames == 15 {
		return winGames
	} else if winGames == 20 {
		return winGames
	} else if winGames == 25 {
		return winGames
	} else if winGames == 40 {
		return winGames
	} else if winGames == 60 {
		return winGames
	} else if winGames == 80 {
		return winGames
	} else if winGames == 100 {
		return winGames
	} else if winGames == 150 {
		return winGames
	} else if winGames == 200 {
		return winGames
	} else if winGames == 250 {
		return winGames
	} else if winGames == 300 {
		return winGames
	} else if winGames == 400 {
		return winGames
	} else if winGames == 500 {
		return winGames
	}

	return 0
}

type playerPlayTaskInfo struct {
	TaskType int   `json:"taskType"` // 1 playTimes, 2 winTimes, 3 winCoins
	Count    int64 `json:"count"`
}

type playerPlayTaskCompleteNotify struct {
	TaskList []playerPlayTaskInfo `json:"taskList"`
}

func (p *GamePlayer) checkTaskFuc(windGames int, playGames int, old, totalWinCoins int64) {
	glog.Info("checkTaskFuc userID:", p.User.UserId, "winGame:", windGames, "playGames:", playGames, "oldWinCoins:", old, "newWinCoins:", totalWinCoins)
	notify := playerPlayTaskCompleteNotify{}
	win := checkPlayWinGames(windGames)
	if win != 0 {
		notify.TaskList = append(notify.TaskList, playerPlayTaskInfo{2, int64(win)})
	}

	out := checkWinCoins(old, totalWinCoins)
	for i := 0; i < len(out); i++ {
		notify.TaskList = append(notify.TaskList, playerPlayTaskInfo{3, out[i]})
	}

	if playGames == 5 {
		notify.TaskList = append(notify.TaskList, playerPlayTaskInfo{1, 5})
	} else if playGames == 10 {
		notify.TaskList = append(notify.TaskList, playerPlayTaskInfo{1, 10})
	}

	glog.Info("checkTaskFuc userID:", p.User.UserId, "notify len:", len(notify.TaskList))
	if len(notify.TaskList) > 0 {
		p.SendToClientNew(int32(pb.MessageIDPlayingTaskNotify), &notify)
	}
}

//AddWinGames ...
func (p *GamePlayer) AddWinGames(bWin bool, coins int64) bool {
	oldWinCoins := p.User.TotalWinCoins
	if bWin {
		p.User.WinGames++
		p.User.TotalWinCoins += coins
	}
	p.User.TotalGames++

	id := p.User.ID
	winGames := p.User.WinGames
	playGames := p.User.TotalGames
	totalWinCoins := p.User.TotalWinCoins
	go changeGameTimes(id, winGames, playGames, totalWinCoins)
	go p.checkTaskFuc(winGames, playGames, oldWinCoins, totalWinCoins)

	return true
}

func (p *GamePlayer) ChangeScore(score pb.Currency, reason string) bool {
	if score < 0 {
		temp := p.User.Balance + score
		if temp < 0 {
			return false
		}
	}

	old := int64(p.User.Balance)
	p.User.Balance += score

	temp := p.User.Balance
	changeScore(p.User.ID, temp)

	logNew := &core.BalanceChangeLog{}
	logNew.UserID = p.User.UserId
	logNew.Reason = reason
	logNew.Old = int64(old)
	logNew.Value = int64(score)
	logNew.Result = int64(p.User.Balance)

	go SaveBalanceChangeLogNew(logNew)

	msg := &pb.DZWalletUpdateNotify{}
	msg.Balance = p.User.Balance
	msg.Frozen = 0
	msg.Coins = p.User.Coins
	msg.Beans = p.User.Bean
	GetPlayerManager().SendClientMsgNew(p.User.UserId, int32(pb.MessageId_DZ_BALANCE_NOTIFY), msg)

	//go rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_Balance, p.User.UserId, p.User.Nickname, p.User.PhotoUrl, p.User.Balance)
	return true
}

func (p *GamePlayer) ResetScore(score pb.Currency, reason string) bool {
	old := p.User.Balance
	p.User.Balance = score

	temp := p.User.Balance
	changeScore(p.User.ID, temp)

	logNew := &core.BalanceChangeLog{}
	logNew.UserID = p.User.UserId
	logNew.Reason = reason
	logNew.Old = int64(old)
	logNew.Value = int64(score)
	logNew.Result = int64(p.User.Balance)

	go SaveBalanceChangeLogNew(logNew)

	msg := &pb.DZWalletUpdateNotify{}
	msg.Balance = p.User.Balance
	msg.Frozen = 0
	msg.Coins = p.User.Coins
	msg.Beans = p.User.Bean
	GetPlayerManager().SendClientMsgNew(p.User.UserId, int32(pb.MessageId_DZ_BALANCE_NOTIFY), msg)

	//go rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_Balance, p.User.UserId, p.User.Nickname, p.User.PhotoUrl, p.User.Balance)
	return true
}

func (p *GamePlayer) ChangeCoinsForcely(score pb.Currency, reason string) bool {
	old := p.User.Coins
	p.User.Coins += score

	temp := p.User.Coins
	changeCoins(p.User.ID, temp)

	logNew := &core.CoinChangeLog{}
	logNew.UserID = p.User.UserId
	logNew.Reason = reason
	logNew.Old = int64(old)
	logNew.Value = int64(score)
	logNew.Result = int64(p.User.Coins)

	go SaveCoinChangeLog(logNew)

	msg := &pb.DZWalletUpdateNotify{}
	msg.Balance = p.User.Balance
	msg.Frozen = 0
	msg.Coins = p.User.Coins
	msg.Beans = p.User.Bean
	GetPlayerManager().SendClientMsgNew(p.User.UserId, int32(pb.MessageId_DZ_BALANCE_NOTIFY), msg)
	if p.User.IdType != int(util.IdType_Robot) {
		go AddRank(RankingBalance, p.User.UserId, int(score))
	}
	return true
}

func (p *GamePlayer) ChangeCoins(score pb.Currency, reason string) bool {
	if p.User.IdType == int(util.IdType_Robot) {
		return true
	}

	if score < 0 {
		temp := p.User.Coins + score
		if temp < 0 {
			return false
		}
	}
	old := p.User.Coins
	p.User.Coins += score

	temp := p.User.Coins
	changeCoins(p.User.ID, temp)

	logNew := &core.CoinChangeLog{}
	logNew.UserID = p.User.UserId
	logNew.Reason = reason
	logNew.Old = int64(old)
	logNew.Value = int64(score)
	logNew.Result = int64(p.User.Coins)

	go SaveCoinChangeLog(logNew)

	msg := &pb.DZWalletUpdateNotify{}
	msg.Balance = p.User.Balance
	msg.Frozen = 0
	msg.Coins = p.User.Coins
	msg.Beans = p.User.Bean
	GetPlayerManager().SendClientMsgNew(p.User.UserId, int32(pb.MessageId_DZ_BALANCE_NOTIFY), msg)
	if p.User.IdType != int(util.IdType_Robot) {
		go AddRank(RankingBalance, p.User.UserId, int(score))
	}
	return true
}

//ChangeBeans ...
func (p *GamePlayer) ChangeBeans(score pb.Currency, reason string) bool {
	if p.User.IdType == int(util.IdType_Robot) {
		return true
	}

	if score < 0 {
		temp := p.User.Bean + score
		if temp < 0 {
			return false
		}
	}

	old := p.User.Bean
	p.User.Bean += score

	temp := p.User.Bean
	changeBeans(p.User.ID, temp)

	log := &core.BeanChangeLog{}
	log.UserID = p.User.UserId
	log.Reason = reason
	log.Old = int64(old)
	log.Value = int64(score)
	log.Result = int64(p.User.Bean)

	go SaveBeanChangeLog(log)

	msg := &pb.DZWalletUpdateNotify{}
	msg.Balance = p.User.Balance
	msg.Frozen = 0
	msg.Coins = p.User.Coins
	msg.Beans = p.User.Bean
	GetPlayerManager().SendClientMsgNew(p.User.UserId, int32(pb.MessageId_DZ_BALANCE_NOTIFY), msg)
	if p.User.IdType != int(util.IdType_Robot) {
		go AddRank(RankingBalance, p.User.UserId, int(score))
	}
	return true
}

func (p *GamePlayer) ResetCoins(score pb.Currency, reason string) bool {
	old := p.User.Coins
	p.User.Coins = score

	temp := p.User.Coins
	changeCoins(p.User.ID, temp)

	logNew := &core.CoinChangeLog{}
	logNew.UserID = p.User.UserId
	logNew.Reason = reason
	logNew.Old = int64(old)
	logNew.Value = int64(score)
	logNew.Result = int64(p.User.Coins)

	go SaveCoinChangeLog(logNew)

	msg := &pb.DZWalletUpdateNotify{}
	msg.Balance = p.User.Balance
	msg.Frozen = 0
	msg.Coins = p.User.Coins
	msg.Beans = p.User.Bean
	GetPlayerManager().SendClientMsgNew(p.User.UserId, int32(pb.MessageId_DZ_BALANCE_NOTIFY), msg)
	if p.User.IdType != int(util.IdType_Robot) {
		go AddRank(RankingBalance, p.User.UserId, int(score))
	}
	return true
}

func (p *GamePlayer) ChangeFirstClubId(clubId int) bool {
	p.User.FirstClubId = clubId
	changeFirstClub(p.User.ID, clubId)

	return true
}

func (p *GamePlayer) SendToClient(msgId int32, body proto.Message) {
	if p.SendToClientFunc != nil {
		p.SendToClientFunc(msgId, body)
	}
}

//发送信息
func (p *GamePlayer) SendToClientNew(msgId int32, body interface{}) {
	if p.SendToClientFuncNew != nil {
		p.SendToClientFuncNew(msgId, body)
	}
}

func (p *GamePlayer) OnLogin() bool {
	p.LoginTime = time.Now()
	go SaveUserLastLoginInfo(p.User.UserId) //最后登陆时间

	p.DZMatchRecord = FindDZMatchRecord(p.User.UserId)
	glog.Info("===>GamePlayer OnLogin DZMatchRecord ", p.DZMatchRecord)
	return true
}

func (p *GamePlayer) OnLogout() {
	p.stopLoop = true
	p.SendToClientFunc = nil
	GetPlayerManager().DelItem(p.User.UserId, false)
	GetBackgroundUserManager().DelUser(p.User.UserId)
	if p.OnLogoutFunc != nil {
		p.OnLogoutFunc(p.User.UserId)
		p.OnLogoutFunc = nil
	}

	glog.Info("===>GamePlayer OnLogout userId:", p.User.UserId, " sessKey:", p.SessKey, " LoginIP:", p.LoginIP, " loginTime:", p.LoginTime)

	SaveDZMatchRecord(p.DZMatchRecord)
}

func (p *GamePlayer) saveLoginRecord() {
	loginRecord := &LoginRecord{}
	loginRecord.UserId = p.User.UserId
	loginRecord.UserName = p.User.UserName
	loginRecord.Channel = p.User.ChannelId
	loginRecord.LoginTime = p.LoginTime
	loginRecord.LogoutTime = time.Now()
	loginRecord.LoginIP = p.LoginIP
	loginRecord.DeviceId = p.LoginDeviceId
	InsertLoginRecord(loginRecord)
}

func (p *GamePlayer) GetDeskUserInfo() *pb.NDeskUserDef {
	userInfo := &pb.NDeskUserDef{}
	userInfo.UserId = p.User.UserId
	userInfo.Nickname = p.User.Nickname
	userInfo.Gender = p.User.Gender
	userInfo.Signiture = p.User.Signiture
	userInfo.PhotoUrl = p.User.PhotoUrl
	userInfo.JingDu = p.User.JingDu
	userInfo.WeiDu = p.User.WeiDu
	userInfo.IP = p.User.IP

	userInfo.Gold = int64(p.User.Coins)
	userInfo.Diamond = int(p.User.Balance)

	match := pb.NUserMatchRecordDef{}
	if p.MatchRecord != nil {
		match.PlayWinCount = p.MatchRecord.WinTimes
		match.PlayTotalCount = p.MatchRecord.LoseTimes + p.MatchRecord.LoseTimes
		match.TheDayWinGold = p.MatchRecord.CurDayEarnGold
		match.TheWeekWinGold = p.MatchRecord.CurWeekEarnGold
		match.MaxCards = p.MatchRecord.MaxCards
	}

	userInfo.MatchRecord = match
	return userInfo
}

func (p *GamePlayer) DZGetDeskUserInfo() *pb.DZDeskUserDef {
	userInfo := &pb.DZDeskUserDef{}
	userInfo.UserId = p.User.UserId
	userInfo.Nickname = p.User.Nickname
	userInfo.Gender = p.User.Gender
	userInfo.Signiture = p.User.Signiture
	userInfo.PhotoUrl = p.User.PhotoUrl
	userInfo.JingDu = p.User.JingDu
	userInfo.WeiDu = p.User.WeiDu
	userInfo.IP = p.User.IP
	userInfo.IdType = p.User.IdType
	if p.MatchRecord != nil {
		userInfo.MatchInfo = *p.DZMatchRecord.BuildMessage()
		if p.User.IdType == int(util.IdType_Robot) {
			userInfo.MatchInfo.VPIP = rand.Int()%20 + 20
			userInfo.MatchInfo.GameTimes = rand.Int() % 50
			userInfo.MatchInfo.PlayTimes = rand.Int() % 200
			userInfo.MatchInfo.WinRate = rand.Int()%20 + 10
		}
	} else {
		userInfo.MatchInfo = pb.DZMatchRecordDef{}
	}
	return userInfo
}

func (p *GamePlayer) DZGetClubUserInfo() *pb.DZClubUser {
	userInfo := &pb.DZClubUser{}
	userInfo.UserId = p.User.UserId
	userInfo.Nickname = p.User.Nickname
	userInfo.HeadUrl = p.User.PhotoUrl
	return userInfo
}
