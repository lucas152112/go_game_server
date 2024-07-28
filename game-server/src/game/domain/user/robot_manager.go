package user

import (
	"fmt"
	"game/pb"
	"game/util"
	"math/rand"
	"sync"

	"github.com/golang/glog"
)

//RobotControlMsg ...
/*type RobotControlMsg struct {
	ControlID int32  //1:获取，2:释放
	GameType  int32  //1:比赛，2:普通游戏
	RobotID   string //要释放的机器人ID
	GameID    int32  //比赛ID或者游戏ID
}*/

//RobotManager ...
type RobotManager struct {
	sync.RWMutex
	items     []string //map[string]int
	robotList map[string]int
	//mq             chan *RobotControlMsg
	robotNameIndex int
}

var robotManager *RobotManager

func init() {
	robotManager = &RobotManager{}
	robotManager.items = []string{} //make(map[string]int)
	robotManager.robotList = make(map[string]int)
	robotManager.robotNameIndex = 1
	//go robotManager.Run()
}

//GetRobotManager ...
func GetRobotManager() *RobotManager {
	return robotManager
}

//AddMsg ...
/*func (m *RobotManager) AddMsg(controlID int32, gameType int32, robotID string) {
	msg := &RobotControlMsg{}
	msg.ControlID = controlID
	msg.GameType = gameType
	msg.RobotID = robotID

	m.mq <- msg

	return
}*/

//Run ...
/*func (m *RobotManager) Run() {
	for {
		select {
		case msg, ok := <-m.mq:
			if !ok {
				glog.Info("<-s.mq != ok")
				return
			}

			go m.ProcessMsg(msg)
		}
	}
}*/

//ProcessMsg ...
/*func (m *RobotManager) ProcessMsg(msg *RobotControlMsg) {
	glog.Info("ProcessMsg::msg:", msg)
	if msg.ControlID == 1 {
		if msg.GameType == 1 {

		} else {
			glog.Info("ProcessMsg::gameID:", msg.GameID)
		}
	} else {
		if msg.RobotID != "" {
			m.AddRobot(msg.RobotID)
		} else {
			glog.Error("ProcessMsg::robotID=null")
		}
	}
	return
}*/

//AddRobot ...
func (m *RobotManager) AddRobot(robotID string) {
	m.Lock()
	defer m.Unlock()
	m.items = append(m.items, robotID)

	return
}

//GetRobot ...
func (m *RobotManager) GetRobot() (string, int) {
	if m.robotNameIndex > 325 {
		return "", 0
	}

	userID := ""
	tempID := 0
	m.Lock()
	if len(m.items) > 100 {
		userID = m.items[0]
		m.items = append(m.items[:0], m.items[1:]...)

		index, ok := m.robotList[userID]
		if ok {
			m.Unlock()
			return userID, index
		}
	}

	m.robotNameIndex++
	tempID = m.robotNameIndex

	glog.Info("GetRobot::index:", tempID)
	m.Unlock()

	userID = ProcessRobotLogin(tempID)
	m.Lock()
	m.robotList[userID] = tempID
	m.Unlock()
	return userID, int(tempID)
}

//GetRobotInfo ...
func GetRobotInfo(robotID string) *pb.DZDeskUserDef {
	user, err := FindByUserId(robotID)
	if err != nil {
		glog.Info("GetRobotInfo::FindByUserId err userID:", robotID)
		return nil
	}
	userInfo := &pb.DZDeskUserDef{}
	userInfo.UserId = user.UserId
	userInfo.Nickname = user.Nickname
	userInfo.Gender = user.Gender
	userInfo.Signiture = user.Signiture
	userInfo.PhotoUrl = user.PhotoUrl
	userInfo.JingDu = user.JingDu
	userInfo.WeiDu = user.WeiDu
	userInfo.IP = user.IP
	userInfo.IdType = user.IdType

	record := FindDZMatchRecord(robotID)
	if record != nil {
		userInfo.MatchInfo = *record.BuildMessage()
		userInfo.MatchInfo.VPIP = rand.Int()%20 + 20
		userInfo.MatchInfo.GameTimes = rand.Int() % 50
		userInfo.MatchInfo.PlayTimes = rand.Int() % 200
		userInfo.MatchInfo.WinRate = rand.Int()%20 + 10

	} else {
		userInfo.MatchInfo = pb.DZMatchRecordDef{}
	}

	return userInfo
}

//ProcessRobotLogin ...
func ProcessRobotLogin(robotNameIndex int) string {
	userName := fmt.Sprintf("robotdz%v", robotNameIndex)
	channel := "robot"
	u, err := FindByUserName(userName)
	if err != nil {
		if err.Error() == "not find" {
			u.UserId = GetNewUserId()
			u.UserName = userName
			u.Nickname = fmt.Sprintf("G%v", u.UserId)
			u.Password = userName
			u.CreateTime = util.GetDbTime()
			u.ChannelId = channel
			u.IsBind = false
			u.IsGuest = true
			u.BGuest = false
			u.Coins = 50000
			u.Balance = 0
			u.IdType = int(util.IdType_Robot)

			if rand.Float64() < 0.5 {
				u.Gender = 1
				u.PhotoUrl = fmt.Sprintf("%v", rand.Int()%4)
			} else {
				u.Gender = 2
				u.PhotoUrl = fmt.Sprintf("%v", 3+rand.Int()%4)
			}

			err = SaveUser(u)
			if err != nil {
				glog.Info("ProcessRobotLogin saveUser error")
				return ""
			}
		} else {
			glog.Info("ProcessRobotLogin FindByUserName error:", err)
			return ""
		}
	}

	return u.UserId
}
