package user

import (
	"encoding/json"
	"game/pb"
	"game/server"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"

	"github.com/jinzhu/gorm"

	// Enable PostgreSQL support
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// Enable MySQL support

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type PlayerManager struct {
	sync.RWMutex
	items   map[string]*server.Session
	online  int
	UserIds []string
}

var playerManager *PlayerManager

func (m *PlayerManager) FindPlayerById(userId string) *GamePlayer {
	m.RLock()
	defer m.RUnlock()
	sess, ok := m.items[userId]
	if !ok {
		return nil
	}

	p := GetPlayer(sess.Data)
	if p != nil {
		if p.DZMatchRecord == nil {
			info := &DZMatchRecord{}
			info.UserId = userId
			p.DZMatchRecord = info
		}
	}
	return p
}

func init() {
	playerManager = &PlayerManager{}
	playerManager.items = make(map[string]*server.Session)
	playerManager.UserIds = []string{}
}

func GetPlayerManager() *PlayerManager {
	return playerManager
}

var (
	// DB is the global datbase object
	DB *gorm.DB
)

//Init ...
func (m *PlayerManager) Init(dbtype string, param string) {
	glog.Info("Initializing database:", dbtype, " param:", param)
	var errDb error
	//db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	DB, errDb = gorm.Open(dbtype, param)
	if errDb != nil {
		glog.Info("open db error:", errDb)
	} else {
		DB.AutoMigrate(&User{})

	}
	/*err, id := GetUserRandId()
	if err != nil && id == "" {
		for i := 118000; i < 119000; i++ {
			temp := fmt.Sprintf("%v", i)
			m.UserIds = append(m.UserIds, temp)
		}
		//
		sort.Sort(UserIdSlice(m.UserIds))
		//	sort.Sort(UserIdSlice(m.UserIds))
		//	sort.Sort(UserIdSlice(m.UserIds))
		//	sort.Sort(UserIdSlice(m.UserIds))
		//	sort.Sort(UserIdSlice(m.UserIds))
		//
		//	//190000
		for i := 0; i < 500; i++ {
			t := &UserRandId_T{m.UserIds[i]}
			SaveUserRandId(t)
		}
		//
		//	for i := 0; i < 2; i++ {
		//		id := m.GetRandUserId()
		//
		//		glog.Info("GetRandUserId userId = ", id)
		//	}
	}*/
}

func (this *PlayerManager) GetRandUserId() string {
	this.Lock()
	defer this.Unlock()

	err, userId := GetUserRandId()
	if err == nil {
		return userId
	}

	return ""
}

func (m *PlayerManager) AddUserItem(userId string, newSess *server.Session) {
	m.RLock()
	oldSess, ok := m.items[userId]
	m.RUnlock()
	if ok && oldSess != nil {
		glog.Info("AddUserItem userId = ", userId, ",newSess:", newSess, ",oldSess:", oldSess)
		oldSess.LoggedIn = false
		oldSess.OnLogout = nil
		oldSess.Data = nil
		m.Lock()
		delete(m.items, userId)
		m.Unlock()
		//删除此Session
		oldSess.ClearConn()
	}
	m.Lock()
	m.items[userId] = newSess
	m.Unlock()
}

func (m *PlayerManager) AddItem(userId string, isRobot bool, newSess *server.Session) bool {
	m.Lock()
	defer m.Unlock()

	oldSess, ok := m.items[userId]
	if ok {
		// 已在线
		if oldSess != nil {
			if oldSess.GetConn() == newSess.GetConn() && newSess == oldSess {
				oldSess.OnLogout = nil
				oldSess.ClearConn()
				return true
			} else {
				p := GetPlayer(oldSess.Data)
				if p != nil {
					p.Stop()
					p.OnLogoutFunc = nil
				}
				oldSess.Kickout()
				delete(m.items, userId)
				return false
			}
		}
	}

	m.items[userId] = newSess

	if !isRobot {
		m.online++
	}

	return true
}

func (m *PlayerManager) ChangeItem(userId string, sess *server.Session) bool {
	m.Lock()
	defer m.Unlock()

	if old, ok := m.items[userId]; ok {
		p := GetPlayer(old.Data)
		if p != nil {
			p.Stop()
		}
		m.items[userId] = sess
		return true
	}

	return false
}

func (m *PlayerManager) DelItem(userId string, isRobot bool) {
	m.RLock()
	old, ok := m.items[userId]
	m.RUnlock()

	if ok {
		p := GetPlayer(old.Data)
		if p != nil {
			p.Stop()
		}
	}
	m.Lock()
	delete(m.items, userId)
	m.Unlock()

	if !isRobot {
		m.online--
	}
}

//RemoveItem ...
func (m *PlayerManager) RemoveItem(userID string) {
	m.Lock()
	delete(m.items, userID)
	m.Unlock()
}

func (m *PlayerManager) Kickout(userId string) {
	sess := m.getSess(userId)
	if sess != nil {
		p := GetPlayer(sess.Data)
		if p != nil {
			p.OnLogoutFunc = nil
			p.Stop()
		}
		sess.Kickout()
	}

	m.Lock()
	delete(m.items, userId)
	m.Unlock()
}

func (m *PlayerManager) FindSessById(userId string) (*server.Session, bool) {
	m.RLock()
	defer m.RUnlock()

	sess, ok := m.items[userId]
	return sess, ok
}

func (m *PlayerManager) getSess(userId string) *server.Session {
	m.RLock()
	defer m.RUnlock()

	return m.items[userId]
}

func (m *PlayerManager) IsOnline(userId string) bool {
	return m.getSess(userId) != nil
}

func (m *PlayerManager) SendServerMsg(srcId string, dstIds []string, msgId int32, body interface{}) {
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return
		}

		m.SendServerMsg2New(srcId, dstIds, msgId, b)
		return
	}
	m.SendServerMsg2(srcId, dstIds, msgId, nil)
}

func (m *PlayerManager) SendServerMsg2(srcId string, dstIds []string, msgId int32, body []byte) {
	go func() {
		msg := &pb.ServerMsg{}
		msg.Client = proto.Bool(false)
		msg.SrcId = proto.String(srcId)
		msg.MsgId = proto.Int32(msgId)
		msg.MsgBody = body

		for _, dstId := range dstIds {
			sess := m.getSess(dstId)
			if sess != nil {
				//sess.SendMQ(msg)
			} else {
				return
			}
		}
	}()
}

func (m *PlayerManager) SendClientMsg(userId string, msgId int32, body proto.Message) {
	sess := m.getSess(userId)
	if sess == nil {
		return
	}
	if GetBackgroundUserManager().Filter(userId, msgId) {
		return
	}
	sess.SendToClient(server.BuildClientMsg(msgId, body))
}

func (m *PlayerManager) SendClientMsg2(dstIds []string, msgId int32, body proto.Message) {
	if len(dstIds) <= 0 {
		return
	}

	b := server.BuildClientMsg(msgId, body)

	for _, userId := range dstIds {
		if GetBackgroundUserManager().Filter(userId, msgId) {
			continue
		}
		sess := m.getSess(userId)
		if sess != nil {
			sess.SendToClient(b)
		}
	}
}

func (m *PlayerManager) Todo(cb func(string)) {
	defer m.RUnlock()
	m.RLock()
	for userId := range m.items {
		cb(userId)
	}
}

func (m *PlayerManager) BroadcastClientMsg(msgId int32, body interface{}) {
	go func() {
		m.RLock()
		items := []*server.Session{}
		for userId, item := range m.items {
			if GetBackgroundUserManager().Filter(userId, msgId) {
				continue
			}
			items = append(items, item)
		}
		m.RUnlock()

		b := server.BuildClientMsg(msgId, body)

		if server.GetServerInstance().IsRefuseService() {
			return
		}

		for _, item := range items {
			if item != nil {
				item.SendToClient(b)
			}
		}
	}()
}

func (m *PlayerManager) GetOnlineCount() int {
	m.RLock()
	defer m.RUnlock()

	return m.online
}

func (m *PlayerManager) GetOnlineCountWithRobot() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.items) * 5
}

//获取所有在线用户
func (m *PlayerManager) GetOnlineUser() []string {

	glog.Info("BroadcastClientMsg")

	m.RLock()
	users := []string{}
	for userId, _ := range m.items {

		users = append(users, userId)
	}

	m.RUnlock()
	return users
}

func (m *PlayerManager) SendClientMsgNew(userId string, msgId int32, body interface{}) {
	sess := m.getSess(userId)
	if sess == nil {
		return
	}

	b := server.BuildClientMsg(msgId, body)
	sess.SendToClient(b)
	//glog.Info("SendClientMsgNew userId:", userId, ",msgId:", msgId)
}

func (m *PlayerManager) SendClientMsg2New(dstIds []string, msgId int32, body interface{}) {
	if len(dstIds) <= 0 {
		return
	}

	b := server.BuildClientMsg(msgId, body)

	for _, userId := range dstIds {
		//glog.Info("SendClientMsg2New userId:",userId)
		//过滤
		if GetBackgroundUserManager().Filter(userId, msgId) {
			continue
		}
		sess := m.getSess(userId)
		if sess != nil {
			sess.SendToClient(b)
			//glog.Info("SendClientMsg2New userId:",userId," success")
		}
	}
}

func (m *PlayerManager) SendServerMsgNew(srcId string, dstIds []string, msgId int32, body interface{}) {
	if body != nil {
		//b := server.BuildClientMsg(msgId, body)
		b, err := json.Marshal(body)
		if err != nil {
			panic(err)
			return
		}

		m.SendServerMsg2New(srcId, dstIds, msgId, b)
		return
	}
	m.SendServerMsg2(srcId, dstIds, msgId, nil)
}

func (m *PlayerManager) SendServerMsg2New(srcId string, dstIds []string, msgId int32, body []byte) {
	go func() {
		msg := &server.ClientMsg{}
		msg.MsgId = msgId
		msg.MsgBody = string(body)

		for _, dstId := range dstIds {
			sess := m.getSess(dstId)
			if sess != nil {
				sess.SendMQ(msg)
			} else {
				return
			}
		}
	}()
}

func (m *PlayerManager) BroadcastClientMsgEx(msgId int32, times int, interval int, body interface{}) {
	go func() {
		for i := 0; i < times; i++ {
			m.RLock()
			items := []*server.Session{}
			for _, item := range m.items {
				items = append(items, item)
			}
			m.RUnlock()

			b := server.BuildClientMsg(msgId, body)

			for _, item := range items {
				if item != nil {
					item.SendToClient(b)
				}
			}

			time.Sleep(time.Minute * time.Duration(interval))
		}
	}()
}

func (m *PlayerManager) GetPlayerCount() int {
	m.Lock()
	defer m.Unlock()

	return len(m.items)
}

func (m *PlayerManager) StopServer() {
	for _, player := range m.items {
		//修复这个Bug
		if player.OnLogout != nil {
			player.OnLogout()
		}

	}
}

type UserIdSlice []string

func (p UserIdSlice) Len() int           { return len(p) }
func (p UserIdSlice) Less(i, j int) bool { return rand.Float32() < 0.5 }
func (p UserIdSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (m *PlayerManager) BroadcastToRobotClientMsg(msgId int32, body interface{}) {
	go func() {
		m.RLock()
		items := []*server.Session{}
		for userId, item := range m.items {
			if GetBackgroundUserManager().Filter(userId, msgId) {
				continue
			}
			player := GetPlayer(item.Data)
			if player != nil {
				if player.User.IdType != 20 {
					continue
				}
			}
			items = append(items, item)
		}
		m.RUnlock()

		b := server.BuildClientMsg(msgId, body)

		for _, item := range items {
			if item != nil {
				item.SendToClient(b)
			}
		}
	}()
}
