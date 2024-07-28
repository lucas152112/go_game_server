package user

import (
	"fmt"
	"game/pb"
	"game/util"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ProductFirstState struct {
	ProductId     string `bson:"productId"`     //商品ID
	ProFirstMonth int    `bson:"proFirstMonth"` //当前商品首冲月份

}

type UserFirstCharge struct {
	UserId             string                        `bson:"userId"` //用户ID
	ProductFirstStates map[string]*ProductFirstState `bson:"productFirstStates"`
}

type UserFortuneManager struct {
	sync.RWMutex
	fortune              map[string]*UserFortune
	UpdateGoldInGameFunc func(userId string)
}

var userFortuneManager *UserFortuneManager

func init() {
	userFortuneManager = &UserFortuneManager{}
	userFortuneManager.fortune = make(map[string]*UserFortune)
}

func GetUserFortuneManager() *UserFortuneManager {
	return userFortuneManager
}

func (m *UserFortuneManager) LoadUserFortune(userId string, isRobot bool) bool {
	m.RLock()
	_, ok := m.fortune[userId]
	if ok {
		m.RUnlock()
		return true
	}

	m.RUnlock()

	f, err := FindUserFortune(userId)
	if err != nil && err != mgo.ErrNotFound {
		glog.Error(err)
		return false
	}

	m.Lock()

	f.IsRobot = isRobot
	m.fortune[userId] = f

	m.Unlock()

	if !util.CompareDate(f.LastPayTime, time.Now()) {
		f.GainedFirstRechargeBonus = false
	}

	return true
}

func (m *UserFortuneManager) UnloadUserFortune(userId string) {
	m.Lock()
	defer m.Unlock()

	f, ok := m.fortune[userId]
	if !ok || f == nil {
		return
	}

	delete(m.fortune, userId)

	if f.IsRobot {
		return
	}

	if f != nil {
		err := SaveFortune(f)
		if err != nil {
			glog.Error("保存用户财富信息失败userId:", userId, " f:", f)
		}
	}
}

func (m *UserFortuneManager) GetUserFortune(userId string) (UserFortune, bool) {
	m.RLock()
	defer m.RUnlock()

	f := m.fortune[userId]
	if f == nil {
		return UserFortune{}, false
	}

	return *f, true
}

func (m *UserFortuneManager) GetUserInfoFortune(userId string) (*UserFortune, bool) {
	m.RLock()
	f := m.fortune[userId]
	defer m.RUnlock()
	if f == nil {
		ff, err := FindUserFortune(userId)
		if err == nil {
			return ff, true
		}
		if err != nil && err != mgo.ErrNotFound {
			glog.Error(err)
			return nil, false
		}

		if err == mgo.ErrNotFound {
			return &UserFortune{}, true
		}
	}

	return f, true
}

func (m *UserFortuneManager) checkUserFortune(userId string) {
	_, ok := m.GetUserFortune(userId)
	if !ok {
		m.LoadUserFortune(userId, false)
	}
}

func (m *UserFortuneManager) SaveUserFortune(userId string) {
	m.Lock()
	f := m.fortune[userId]
	m.Unlock()

	if f == nil {
		return
	}

	SaveFortune(f)
}

func (m *UserFortuneManager) EarnFortune(userId string, gold int64, diamond, score int, isRecharge bool, reason string) bool {
	m.checkUserFortune(userId)

	m.Lock()
	defer m.Unlock()

	glog.Info("===>EarnFortune userId:", userId, " gold:", gold, " diamond:", diamond, " reason:", reason)

	fortune := m.fortune[userId]
	if fortune == nil {
		glog.Error("EarnFortune failed userId:", userId, " gold:", gold, " reason:", reason)
		return false
	}

	fortune.Diamond += diamond
	fortune.Gold += int64(gold)

	fortune.Score += score

	if fortune.Gold < 0 {
		fortune.Gold = 0
	}
	if fortune.Diamond < 0 {
		fortune.Diamond = 0
	}

	l := &FortuneLog{}
	l.UserId = userId
	l.Gold = int(gold)
	l.CurGold = fortune.Gold
	l.Diamond = diamond
	l.CurDiamond = fortune.Diamond
	l.Reason = reason
	SaveEarnFortuneLog(l)

	if gold != 0 {
		m.updateGoldInGame(userId)
		SaveGold(userId, fortune.Gold, fortune.IsRobot)
	}

	if diamond != 0 {
		m.updateGoldInGame(userId)
		SaveDiamond(userId, fortune.Diamond, fortune.IsRobot)
	}

	m.updateUserFortune(userId, 0)

	return true
}

func (m *UserFortuneManager) updateGoldInGame(userId string) {
	if m.UpdateGoldInGameFunc != nil {
		go m.UpdateGoldInGameFunc(userId)
	}
}

func (m *UserFortuneManager) ConsumeGold(userId string, gold int64, consumeAllIfNotEnough bool, reason string) (int64, int, bool) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return 0, 0, false
	}

	if gold > 0 {
		if fortune.Gold < int64(gold) {
			if consumeAllIfNotEnough {
				oldGold := fortune.Gold
				fortune.Gold = 0
				gold = int64(oldGold)
			} else {
				return 0, 0, false
			}
		} else {
			fortune.Gold -= int64(gold)
		}
	}

	if !fortune.IsRobot {
		l := &FortuneLog{}
		l.UserId = userId
		l.Gold = int(gold)
		l.CurGold = fortune.Gold
		l.Diamond = 0
		l.CurDiamond = fortune.Diamond
		l.Reason = reason
		SaveConsumeFortuneLog(l)
	}

	if gold > 0 {
		m.updateGoldInGame(userId)
		SaveGold(userId, fortune.Gold, fortune.IsRobot)
	}

	m.updateUserFortune(userId, 0)

	return fortune.Gold, int(gold), true
}

func (m *UserFortuneManager) ConsumeGoldNoMsg(userId string, gold int64, consumeAllIfNotEnough bool, reason string) (int64, int, bool) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return 0, 0, false
	}

	if gold > 0 {
		if fortune.Gold < int64(gold) {
			if consumeAllIfNotEnough {
				oldGold := fortune.Gold
				fortune.Gold = 0
				gold = int64(oldGold)
			} else {
				return 0, 0, false
			}
		} else {
			fortune.Gold -= int64(gold)
		}
	}

	if !fortune.IsRobot {
		l := &FortuneLog{}
		l.UserId = userId
		l.Gold = int(gold)
		l.CurGold = fortune.Gold
		l.Diamond = 0
		l.CurDiamond = fortune.Diamond
		l.Reason = reason
		SaveConsumeFortuneLog(l)
	}

	if gold > 0 {
		m.updateGoldInGame(userId)
		SaveGold(userId, fortune.Gold, fortune.IsRobot)
	}

	return fortune.Gold, int(gold), true
}

func (m *UserFortuneManager) ConsumeDiamond(userId string, diamond int, reason string) (int, bool) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return 0, false
	}

	if diamond > 0 {
		if fortune.Diamond < diamond {
			return 0, false
		} else {
			fortune.Diamond -= diamond
		}
	}

	l := &FortuneLog{}
	l.UserId = userId
	l.Gold = 0
	l.CurGold = 0
	l.Diamond = diamond
	l.CurDiamond = fortune.Diamond
	l.Reason = reason
	SaveConsumeFortuneLog(l)

	m.updateUserFortune(userId, 0)

	return fortune.Diamond, true
}

func (m *UserFortuneManager) EarnDiamond(userId string, diamond int, reason string) (int, bool) {
	m.checkUserFortune(userId)
	m.Lock()
	defer m.Unlock()
	fortune := m.fortune[userId]
	if fortune == nil {
		glog.Error("EarnDiamondOnline 失败usreId :", userId)
		return 0, false
	}

	fortune.Diamond += diamond

	l := &FortuneLog{}
	l.UserId = userId
	l.Gold = 0
	l.CurGold = 0
	l.Diamond = diamond
	l.CurDiamond = fortune.Diamond
	l.Reason = reason
	SaveConsumeFortuneLog(l)

	return fortune.Diamond, true
}

func (m *UserFortuneManager) ConsumeThirdCurrency(userId string, thirdCurrency int, reason string) (int, bool) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return 0, false
	}

	if thirdCurrency > 0 {
		if fortune.ThirdCurrency < thirdCurrency {
			return 0, false
		} else {
			fortune.ThirdCurrency -= thirdCurrency
		}
	}

	l := &FortuneLog{}
	l.UserId = userId
	l.Gold = 0
	l.CurGold = 0
	l.Diamond = 0
	l.CurDiamond = 0
	l.ThirdCurrency = thirdCurrency
	l.CurThirdCurrency = fortune.ThirdCurrency
	l.Reason = reason
	SaveConsumeFortuneLog(l)

	m.updateUserFortune(userId, 0)

	return fortune.Diamond, true
}

func (m *UserFortuneManager) EarnThirdCurrency(userId string, thirdCurrency int, reason string) (int, bool) {
	m.checkUserFortune(userId)
	m.Lock()
	defer m.Unlock()
	fortune := m.fortune[userId]
	if fortune == nil {
		glog.Error("EarnThirdCurrency 失败usreId :", userId)
		return 0, false
	}

	fortune.ThirdCurrency += thirdCurrency

	l := &FortuneLog{}
	l.UserId = userId
	l.Gold = 0
	l.CurGold = 0
	l.Diamond = 0
	l.CurDiamond = 0
	l.ThirdCurrency = thirdCurrency
	l.CurThirdCurrency = fortune.ThirdCurrency
	l.Reason = reason
	SaveConsumeFortuneLog(l)

	return fortune.Diamond, true
}

//赠送钻石
func (m *UserFortuneManager) GiftDiamondUpdateSafeBoxLog(userId string, pwd string, toUserId string, diamond int) (bool, int, string, *BoxLog) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false, 0, "用户信息错误", nil
	}

	if !fortune.SafeBox.IsOpen {
		return false, 0, "保管箱未激活，请到商场购买", nil
	}

	if fortune.SafeBox.Pwd != pwd {
		return false, 0, "保管箱密码错误", nil
	}

	if fortune.Diamond-diamond < 0 {
		return false, 0, "Insufficient diamonds", nil
	}

	fortune.Diamond -= diamond
	log := BoxLog{"赠送钻石扣除", 0, diamond, int64(fortune.Diamond), time.Now().Local(), toUserId, "0"}
	if fortune.SafeBox.BoxLogs == nil {
		fortune.SafeBox.BoxLogs = map[string]BoxLog{}
	}
	logId := strconv.Itoa(len(fortune.SafeBox.BoxLogs) + 1)
	log.LogId = logId
	fortune.SafeBox.BoxLogs[logId] = log
	m.fortune[userId] = fortune

	return true, fortune.Diamond, "赠送钻石成功！", &log
}

//存保管箱离线
func (m *UserFortuneManager) AddDiamondUpdateSafeBoxLog(fromUserId string, toUserId string, diamond int) bool {
	m.checkUserFortune(toUserId)
	m.Lock()
	defer m.Unlock()
	fortune := m.fortune[toUserId]
	if fortune == nil {
		glog.Error("赠送用户存保管箱离线失败usreId :", toUserId)
		return false
	}

	fortune.Diamond += diamond

	log := BoxLog{"赠送钻石", 0, diamond, int64(fortune.Diamond), time.Now().Local(), fromUserId, "0"}
	if fortune.SafeBox.BoxLogs == nil {
		fortune.SafeBox.BoxLogs = map[string]BoxLog{}
	}
	logId := strconv.Itoa(len(fortune.SafeBox.BoxLogs) + 1)
	log.LogId = logId
	fortune.SafeBox.BoxLogs[logId] = log
	m.fortune[toUserId] = fortune
	err := SaveFortune(fortune)
	if err != nil {
		glog.Error("赠送钻石失败userId:", toUserId, " f:", fortune)
		return false
	}
	return true
}

func (m *UserFortuneManager) AddGold2Mongo(userId string, gold int64, reason string) (int64, bool) {
	m.Lock()
	defer m.Unlock()

	fortune, err := FindUserFortune(userId)
	if err != nil {
		glog.Info("FindUserFortune error, err=", err)
		return 0, false
	}

	fortune.Gold += gold
	SaveFortune(fortune)

	l := &FortuneLog{}

	l.UserId = userId
	l.Gold = int(gold)
	l.CurGold = fortune.Gold
	l.Diamond = 0
	l.CurDiamond = fortune.Diamond
	l.Reason = reason
	SaveEarnFortuneLog(l)

	return fortune.Gold, true
}

func (m *UserFortuneManager) ConsumeCharm(userId string, charm int) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if fortune.Charm < charm {
		return false
	}

	fortune.Charm -= charm
	m.fortune[userId] = fortune

	return true
}

/*func (m *UserFortuneManager) LoginSafeBox(userId string, pwd string) *pb.Msg_LoginSafeBoxRes_ResCode {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return pb.Msg_LoginSafeBoxRes_FAILED.Enum()
	}

	if fortune.SafeBox.Pwd != pwd {
		return pb.Msg_LoginSafeBoxRes_PWDERR.Enum()
	}

	return pb.Msg_LoginSafeBoxRes_OK.Enum()
}*/

func (m *UserFortuneManager) ChangePwdSafeBox(userId, oldpwd, newpwd string) (int, string) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return -1, "获取信息失败"
	}

	if fortune.SafeBox.Pwd != oldpwd {
		return -2, "保管箱密码错误"
	}

	fortune.SafeBox.Pwd = newpwd
	m.fortune[userId] = fortune

	return 0, ""
}

func (m *UserFortuneManager) ResetPwdSafeBox(userId, newpwd string) int {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return -1
	}

	if !fortune.SafeBox.IsSetPwd {
		fortune.SafeBox.IsSetPwd = true
	}

	glog.Info("ResetPwdSafeBox in,newpwd=", newpwd)
	fortune.SafeBox.Pwd = newpwd
	m.fortune[userId] = fortune
	glog.Info("m.fortune[userId].SafeBox.Pwd=", m.fortune[userId].SafeBox.Pwd)

	return 0
}

func (m *UserFortuneManager) OpenSafeBox(userId string, pwd string) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if fortune.SafeBox.IsOpen == true {
		return true
	}

	fortune.SafeBox.IsOpen = true
	fortune.SafeBox.Pwd = pwd
	m.fortune[userId] = fortune

	return true
}

func (m *UserFortuneManager) UpdateSavings(userId, pwd string, gold int64, reason string, toUid string) (int, string, int64, int64, *BoxLog) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return -1, "用户信息错误", 0, 0, nil
	}

	if reason != "取款" && !fortune.SafeBox.IsOpen {
		return -2, "保管箱未激活，请到商场购买", 0, 0, nil
	}

	if reason == "取款" || reason == "赠送扣款" {
		if fortune.SafeBox.Pwd != pwd {
			return -3, "保管箱密码错误", 0, 0, nil
		}
	}

	if fortune.SafeBox.Savings+gold < 0 {
		return -4, "您的存款不足", 0, 0, nil
	}

	fortune.SafeBox.Savings += gold
	log := BoxLog{reason, gold, 0, fortune.SafeBox.Savings, time.Now().Local(), toUid, "0"}
	if fortune.SafeBox.BoxLogs == nil {
		fortune.SafeBox.BoxLogs = map[string]BoxLog{}
	}
	logId := strconv.Itoa(len(fortune.SafeBox.BoxLogs) + 1)
	log.LogId = logId
	fortune.SafeBox.BoxLogs[logId] = log
	m.fortune[userId] = fortune

	return 0, "", fortune.SafeBox.Savings, fortune.Gold, &log
}

//存保管箱离线
func (m *UserFortuneManager) UpdateSavingsAddOffLine(fromUserId string, toUserId string, gold int64) bool {
	m.checkUserFortune(toUserId)
	m.Lock()
	defer m.Unlock()
	fortune := m.fortune[toUserId]
	if fortune == nil {
		glog.Error("赠送用户存保管箱离线失败usreId :", toUserId)
		return false
	}

	fortune.SafeBox.Savings += gold
	fromFromUserId := fromUserId
	log := BoxLog{"赠送", gold, 0, fortune.SafeBox.Savings, time.Now().Local(), fromFromUserId, "0"}
	if fortune.SafeBox.BoxLogs == nil {
		fortune.SafeBox.BoxLogs = map[string]BoxLog{}
	}
	logId := strconv.Itoa(len(fortune.SafeBox.BoxLogs) + 1)
	log.LogId = logId
	fortune.SafeBox.BoxLogs[logId] = log
	m.fortune[toUserId] = fortune
	err := SaveFortune(fortune)
	if err != nil {
		glog.Error("赠送用户存保管箱离线失败userId:", toUserId, " f:", fortune)
		return false
	}
	return true
}

//存保险保管箱在线
func (m *UserFortuneManager) UpdateSavingsAdd(fromUserId string, toUserId string, gold int64) bool {
	m.Lock()
	defer m.Unlock()
	fortune := m.fortune[toUserId]

	if fortune == nil {
		glog.Error("赠送用户财富信息失败usreId :", toUserId)
		return false
	}

	fortune.SafeBox.Savings += gold
	fromFromUserId := fromUserId
	log := BoxLog{"赠送", gold, 0, fortune.SafeBox.Savings, time.Now().Local(), fromFromUserId, "0"}
	if fortune.SafeBox.BoxLogs == nil {
		fortune.SafeBox.BoxLogs = map[string]BoxLog{}
	}
	logId := strconv.Itoa(len(fortune.SafeBox.BoxLogs) + 1)
	log.LogId = logId
	fortune.SafeBox.BoxLogs[logId] = log
	m.fortune[toUserId] = fortune
	err := SaveFortune(fortune)
	if err != nil {
		glog.Error("赠送用户财富信息失败userId:", toUserId, " f:", fortune)
		return false
	}

	return true
}

func (m *UserFortuneManager) UpdateUserFortune(userId string) {
	m.Lock()
	defer m.Unlock()

	m.updateUserFortune(userId, 0)
}

func (m *UserFortuneManager) UpdateUserFortune2(userId string, rechargeDiamond int) {
	m.Lock()
	defer m.Unlock()
	glog.Info("UpdateUserFortune2 in,userId=", userId, "|rechargeDiamond=", rechargeDiamond)
	m.updateUserFortune(userId, rechargeDiamond)
}

func (m *UserFortuneManager) updateUserFortune(userId string, rechargeDiamond int) {
	f := m.fortune[userId]
	if f == nil {
		return
	}

	msg := &pb.NUpdateGold{}
	msg.Gold = f.Gold
	msg.Diamond = f.Diamond
	GetPlayerManager().SendClientMsgNew(userId, int32(pb.MessageId_N_UPDATE_GOLD_NOTIFY), msg)
}

func (m *UserFortuneManager) EarnGold(userId string, gold int64, reason string) (int64, bool) {
	ok := m.EarnFortune(userId, gold, 0, 0, false, reason)
	if !ok {
		return 0, false
	}

	f, _ := m.GetUserFortune(userId)
	m.updateUserFortune(userId, 0)
	return f.Gold, true
}

func (m *UserFortuneManager) ApiChangeGold(userId string, gold int64, channel string) (int64, bool) {
	m.checkUserFortune(userId)

	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		glog.Error("ApiChangeGold failed userId:", userId, " gold:", gold, " channel:", channel)
		return 0, false
	}

	if gold < 0 {
		if (fortune.Gold + gold) < 0 {
			return fortune.Gold, false
		}
	}

	fortune.Gold += gold

	l := &FortuneLog{}
	l.UserId = userId
	l.Gold = int(gold)
	l.CurGold = fortune.Gold
	l.Diamond = 0
	l.CurDiamond = fortune.Diamond
	l.Reason = channel
	go SaveEarnFortuneLog(l)
	go SaveGold(userId, fortune.Gold, false)

	return fortune.Gold, true
}

func (m *UserFortuneManager) EarnGoldNoMsg(userId string, gold int64, reason string) (int64, bool) {
	m.checkUserFortune(userId)

	m.Lock()
	defer m.Unlock()

	glog.Info("===>EarnGoldNoMsg userId:", userId, " gold:", gold, " reason:", reason)

	fortune := m.fortune[userId]
	if fortune == nil {
		glog.Error("EarnGoldNoMsg failed userId:", userId, " gold:", gold, " reason:", reason)
		return 0, false
	}

	fortune.Gold += int64(gold)

	if fortune.Gold < 0 {
		fortune.Gold = 0
	}

	l := &FortuneLog{}
	l.UserId = userId
	l.Gold = int(gold)
	l.CurGold = fortune.Gold
	l.Diamond = 0
	l.CurDiamond = fortune.Diamond
	l.Reason = reason
	SaveEarnFortuneLog(l)

	if gold != 0 {
		go SaveGold(userId, fortune.Gold, fortune.IsRobot)
	}

	return fortune.Gold, true
}

// add by wangsq start
func (m *UserFortuneManager) EarnCharm(userId string, charm int) (int, bool) {
	glog.Info("EarnCharm in.", charm)
	m.checkUserFortune(userId)

	f, _ := m.fortune[userId]
	f.Charm += charm
	SaveCharm(userId, f.Charm, false)
	m.fortune[userId] = f

	return f.Charm, true
}

func (m *UserFortuneManager) EarnHorn(userId string, horn int) (int, bool) {
	glog.Info("EarnHorn in.", horn)
	m.checkUserFortune(userId)

	f, _ := m.fortune[userId]
	f.Horn += horn
	SaveHorn(userId, f.Horn, false)
	m.fortune[userId] = f

	return f.Horn, true
}

func (m *UserFortuneManager) GetCharmExchangeInfo(userId string, itemId int) int {
	glog.Info("GetCharmExchangeInfo in.userId=", userId, "|itemId=", itemId)
	f, _ := m.fortune[userId]

	if f.CharmExchangeInfo == nil {
		f.CharmExchangeInfo = map[string]int{}
	}
	itemId_str := fmt.Sprintf("%d", itemId)
	_, infok := f.CharmExchangeInfo[itemId_str]
	if !infok {
		f.CharmExchangeInfo[itemId_str] = 0
	}

	m.fortune[userId] = f
	return f.CharmExchangeInfo[itemId_str]
}

func (m *UserFortuneManager) UpdateCharmExchangeInfo(userId string, itemId int, count int) {
	glog.Info("UpdateCharmExchangeInfo in.userId=", userId, "|itemId=", itemId)
	f, _ := m.fortune[userId]

	if f.CharmExchangeInfo == nil {
		f.CharmExchangeInfo = map[string]int{}
	}

	itemId_str := fmt.Sprintf("%d", itemId)
	_, infok := f.CharmExchangeInfo[itemId_str]
	if !infok {
		f.CharmExchangeInfo[itemId_str] = 0
	}

	f.CharmExchangeInfo[itemId_str] += count
	m.fortune[userId] = f
}

// add by wangsq end

func (m *UserFortuneManager) ExchangeGold(userId string, diamond int) (bool, int64, int) {
	glog.V(2).Info("ExchangeGold userId:", userId, " diamond:", diamond)

	if diamond <= 0 {
		return false, 0, 0
	}

	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		glog.V(2).Info("===>userId:", userId, " fortune nil")
		return false, 0, 0
	}

	if fortune.Diamond < diamond {
		glog.V(2).Info("==>fortune diamond:", fortune.Diamond, " diamond:", diamond, " <")
		return false, 0, 0
	}

	fortune.Diamond -= diamond
	spGold := 0
	if diamond == 50 {
		spGold = 30000
	} else if diamond == 100 {
		spGold = 100000
	} else if diamond == 300 {
		spGold = 500000
	} else if diamond == 500 {
		spGold = 1000000
	} else if diamond == 1000 {
		spGold = 3000000
	} else if diamond == 108 {
		spGold = 120000
	} else if diamond == 298 {
		spGold = 520000
	} else if diamond == 518 {
		spGold = 1020000
	} else if diamond == 998 {
		spGold = 3020000
	}
	fortune.Gold += int64(diamond*exchangeGoldRate + spGold)

	if !fortune.IsRobot {
		l := &FortuneLog{}
		l.UserId = userId
		l.Gold = 0
		l.CurGold = fortune.Gold
		l.Diamond = diamond
		l.CurDiamond = fortune.Diamond
		l.Reason = "钻石兑换金币"
		SaveConsumeFortuneLog(l)

		l.Gold = diamond*exchangeGoldRate + spGold
		l.Diamond = 0
		SaveEarnFortuneLog(l)
	}

	m.updateGoldInGame(userId)

	SaveDiamond(userId, fortune.Diamond, fortune.IsRobot)
	SaveGold(userId, fortune.Gold, fortune.IsRobot)

	return true, int64(fortune.Gold), fortune.Diamond
}

//获取用户首冲表
func (m *UserFortuneManager) GetUserFirstChargeTable(userIdStr string) (string, error) {
	userIdInt, err := strconv.Atoi(userIdStr)
	if err != nil {

		glog.Info("GetUserPromoterTable err:", err)
		return "", err
	}

	index := userIdInt/30000 + 1

	tableName := "user_firstcharge_" + fmt.Sprintf("%v", index)
	return tableName, nil
}

//根据用户获取月度首冲
func (m *UserFortuneManager) GetUserFirstCharge(userid string) (*UserFirstCharge, error) {
	glog.Info("GetUserFirstCharge userid=", userid)
	userFirstChargeTableC, er := m.GetUserFirstChargeTable(userid)
	if er != nil {
		return nil, er
	}

	userFirstCharge := &UserFirstCharge{}
	err := util.WithUserCollection(userFirstChargeTableC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userid}).One(userFirstCharge)
	})

	return userFirstCharge, err
}

//存储当前用户的月度首充
func (m *UserFortuneManager) SaveUserFirstCharge(userFirstCharge *UserFirstCharge) error {
	userId := userFirstCharge.UserId
	userFirstChargeTableC, er := m.GetUserFirstChargeTable(userId)
	if er != nil {
		return er
	}

	return util.WithUserCollection(userFirstChargeTableC, func(c *mgo.Collection) error {
		return c.Insert(userFirstCharge)
	})

	return nil

}

//修改当前用户的月度首充
func (m *UserFortuneManager) UpdateUserFirstCharge(userFirstCharge *UserFirstCharge) error {
	userId := userFirstCharge.UserId
	userFirstChargeTableC, er := m.GetUserFirstChargeTable(userId)
	if er != nil {
		return er
	}

	glog.Info("userId:", userId)

	return util.WithUserCollection(userFirstChargeTableC, func(c *mgo.Collection) error {
		return c.Update(bson.M{"userId": userId}, bson.M{"$set": bson.M{"productFirstStates": userFirstCharge.ProductFirstStates}})
	})

	return nil

}

func (m *UserFortuneManager) ResetDailyGiftBag(userId string) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return
	}

	now := time.Now()
	if !util.CompareDate(fortune.DailyGiftBagUseTime, now) {
		fortune.BuyDailyGiftBag = false
		fortune.DailyGiftBagUseTime = now
	}
}

func (m *UserFortuneManager) SetBuyDailyGiftBag(userId string) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if util.CompareDate(fortune.DailyGiftBagUseTime, time.Now()) {
		return false
	}

	fortune.BuyDailyGiftBag = true

	return true
}

func (m *UserFortuneManager) ManagerChangeGameTypeNotify(userId string, date int) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	fortune.ChangeGameTypeNotifyDay = date

	return true
}
