package user

import (
	"game/util"
	"sort"
	"sync"
	"time"

	"game/domain/hall"

	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	liver_c             = "liver_c"
	liver_game_config_c = "liver_game_config_c"
)

type Liver struct {
	LiverID string `bson:"liverID"` //主播ID
	Cover   string `bson:"cover"`   //封面
	Status  int    `bson:"status"`  //状态，2：下播，1：在播，3：关停, 4:永久关停
	Notice  string `bson:"notice"`  //公告
}

func SaveLiver(liver *Liver) error {
	return util.WithSafeUserCollection(liver_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"liverID": liver.LiverID}, liver)
		return err
	})
}

func SetLiverStatus(liverID string, status int) error {
	err := util.WithSafeUserCollection(liver_c, func(c *mgo.Collection) error {
		return c.Update(bson.M{"liverID": liverID}, bson.M{"$set": bson.M{"status": status}})
	})
	return err
}

//SetLiverNotice ...
func SetLiverNotice(liverID string, notice string) error {
	err := util.WithSafeUserCollection(liver_c, func(c *mgo.Collection) error {
		return c.Update(bson.M{"liverID": liverID}, bson.M{"$set": bson.M{"notice": notice}})
	})
	return err
}

func InitLiverStatus() {
	items := []*Liver{}
	util.WithSafeUserCollection(liver_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"status": 1}).All(&items)
	})

	for i := 0; i < len(items); i++ {
		util.WithSafeUserCollection(liver_c, func(c *mgo.Collection) error {
			return c.Update(bson.M{"liverID": items[i].LiverID}, bson.M{"$set": bson.M{"status": 2}})
		})
	}

	return
}

func GetLivers(page int, step int) (error, []*Liver) {
	items := []*Liver{}
	if page > 1 {
		page -= 1
	}
	err := util.WithSafeUserCollection(liver_c, func(c *mgo.Collection) error {
		return c.Find(nil).Sort("status").Skip(page * step).Limit(step).All(&items)
		//return c.Find(nil).All(&items)
	})

	glog.Info("GetLivers err:", err, " len:", len(items))

	return err, items
}

func GetAllActiveLivers() (error, []*Liver) {
	items := []*Liver{}
	err := util.WithSafeUserCollection(liver_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"$or": []bson.M{bson.M{"status": 1}, bson.M{"status": 2}}}).All(&items)
	})

	return err, items
}

func GetLiver(liverID string) (error, *Liver) {
	item := &Liver{}
	err := util.WithSafeUserCollection(liver_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"liverID": liverID}).One(item)
	})

	return err, item
}

func RemoveLever(liverID string) error {
	err := util.WithSafeUserCollection(liver_c, func(c *mgo.Collection) error {
		return c.Remove(bson.M{"liverID": liverID})
	})

	return err
}

func RemoveLeverGameConfig(liverID string) error {
	err := util.WithSafeUserCollection(liver_game_config_c, func(c *mgo.Collection) error {
		return c.Remove(bson.M{"liverID": liverID})
	})

	return err
}

type GameConfig struct {
	LiverID       string `bson:"liverID"`       //主播ID
	GameType      int    `bson:"gameType"`      //游戏类型
	BaseScore     int    `bson:"baseScore"`     //小盲注
	BeforeScore   int    `bson:"beforeScore"`   //前注
	PlayerCount   int    `bson:"playerCount"`   //座位数
	GameName      string `bson:"gameName"`      //游戏名称
	IsDUserDouble bool   `bson:"isDUserDouble"` //庄家两倍前注
}

func SaveLiverGameConfig(config *GameConfig) error {
	return util.WithSafeUserCollection(liver_game_config_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"liverID": config.LiverID}, config)
		return err
	})
}

func GetLiverGameConfigs(liverID string) (error, []*GameConfig) {
	items := []*GameConfig{}
	err := util.WithSafeUserCollection(liver_game_config_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"liverID": liverID}).All(&items)
	})

	return err, items
}

func GetLiverGameConfig(liverID string) (error, *GameConfig) {
	item := &GameConfig{}
	err := util.WithSafeUserCollection(liver_game_config_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"liverID": liverID}).One(&item)
	})

	return err, item
}

func insureIndexLiverLiverID() error {
	err := util.WithUserCollection(liver_c, func(c *mgo.Collection) error {
		index := mgo.Index{
			Key:        []string{"liverID"},
			Unique:     true,
			DropDups:   true,
			Background: true, // See notes.
		}
		e := c.EnsureIndex(index)
		return e
	})

	return err
}

func insureIndexLiverStaus() error {
	err := util.WithUserCollection(liver_c, func(c *mgo.Collection) error {
		index := mgo.Index{
			Key:        []string{"status"},
			Background: true, // See notes.
		}
		e := c.EnsureIndex(index)
		return e
	})

	return err
}

func insureIndexLiverGameConfigLiverID() error {
	err := util.WithUserCollection(liver_game_config_c, func(c *mgo.Collection) error {
		index := mgo.Index{
			Key:        []string{"liverID"},
			Unique:     false,
			DropDups:   false,
			Background: true, // See notes.
		}
		e := c.EnsureIndex(index)
		return e
	})

	return err
}

func indexLiverDrop() error {
	err := util.WithUserCollection(liver_c, func(c *mgo.Collection) error {
		indexes, e := c.Indexes()
		if e != nil {
			return e
		}
		for _, index := range indexes {
			glog.Info("indexGetRecordLog:", index.Key)
			c.DropIndex(index.Key...)
		}

		return e
	})

	return err
}
func init() {
	gifterManager = &GifterManager{}
	gifterManager.LiverGifter = make(map[string]*liverOfGifter)
	forbidManager = &ForbidManager{}
	forbidManager.liverForbidList = make(map[string]*liverForbid)

}

//Gifter ...
type Gifter struct {
	UserID   string
	HeadURL  string
	Nickname string
	Diamond  int64
	Status   int //1 online 2 offline
}

type liverOfGifter struct {
	sync.RWMutex
	LiverID string
	Gifters []*Gifter
	Lookers map[string]*Gifter
}

func (item *liverOfGifter) getGifterList(countIn int, page int, typeGet int) []*Gifter {
	item.RLock()
	defer item.RUnlock()
	back := []*Gifter{}
	count := 0
	skipCount := countIn * page
	for i := 0; i < len(item.Gifters); i++ {
		if typeGet == 2 {
			if item.Gifters[i].Status == 2 {
				continue
			}
		}

		if i < skipCount {
			continue
		}
		it := Gifter{}
		it.UserID = item.Gifters[i].UserID
		it.Diamond = item.Gifters[i].Diamond
		it.HeadURL = item.Gifters[i].HeadURL
		it.Nickname = item.Gifters[i].Nickname

		back = append(back, &it)
		count++

		if count >= countIn {
			break
		}
	}

	return back
}

func (item *liverOfGifter) getLookerList(countIn int, page int) []*Gifter {
	item.RLock()
	defer item.RUnlock()
	back := []*Gifter{}
	count := 0
	i := 0
	skipCount := countIn * page
	for _, info := range item.Lookers {
		if i < skipCount {
			i++
			continue
		}
		it := Gifter{}
		it.UserID = info.UserID
		it.Diamond = info.Diamond
		it.HeadURL = info.HeadURL
		it.Nickname = info.Nickname

		back = append(back, &it)
		count++

		if count >= countIn {
			break
		}
	}

	return back
}

func (item *liverOfGifter) getLookerCount() int {
	item.RLock()
	defer item.RUnlock()
	return len(item.Lookers)
}

func (item *liverOfGifter) playerOnline(userID string) {
	item.RLock()
	defer item.RUnlock()
	for i := 0; i < len(item.Gifters); i++ {
		if item.Gifters[i].UserID == userID {
			item.Gifters[i].Status = 1
			return
		}
	}
}

func (item *liverOfGifter) liverOffLine() {
	item.Lock()
	item.Gifters = []*Gifter{}
	for it := range item.Lookers {
		delete(item.Lookers, it)
	}
	item.Unlock()
}

func (item *liverOfGifter) playerOffline(userID string) {
	item.RLock()
	for i := 0; i < len(item.Gifters); i++ {
		if item.Gifters[i].UserID == userID {
			item.Gifters[i].Status = 2
			break
		}
	}

	item.RUnlock()

	item.Lock()
	delete(item.Lookers, userID)
	item.Unlock()
}

func (item *liverOfGifter) myRank(userID string) (int, *Gifter) {
	item.RLock()
	for i := 0; i < len(item.Gifters); i++ {
		if item.Gifters[i].UserID == userID {
			it := &Gifter{}
			it.UserID = item.Gifters[i].UserID
			it.Diamond = item.Gifters[i].Diamond
			it.HeadURL = item.Gifters[i].HeadURL
			it.Nickname = item.Gifters[i].Nickname
			item.RUnlock()
			return i + 1, it
		}
	}

	item.RUnlock()

	return -1, nil
}

func (item *liverOfGifter) addLooker(userID string, headURL string, nickname string) {
	item.RLock()
	_, okGifter := item.Lookers[userID]
	if okGifter {
		item.RUnlock()
		return
	}

	item.RUnlock()

	item.Lock()
	gifter := &Gifter{}
	gifter.Diamond = 0
	gifter.Nickname = nickname
	gifter.HeadURL = headURL
	gifter.UserID = userID
	gifter.Status = 1
	item.Lookers[userID] = gifter
	item.Unlock()

	return
}

func (item *liverOfGifter) addGift(userID string, headURL string, nickname string, diamond int64) {
	bFind := false

	item.RLock()
	_, okGifter := item.Lookers[userID]
	item.RUnlock()

	if !okGifter {
		item.addLooker(userID, headURL, nickname)
	}

	item.RLock()
	for i := 0; i < len(item.Gifters); i++ {
		if item.Gifters[i].UserID == userID {
			temp := item.Gifters[i].Diamond
			temp += diamond

			item.Gifters[i].Diamond = temp
			bFind = true
			break
		}
	}

	item.RUnlock()

	if !bFind {
		item.Lock()
		gifter := &Gifter{}
		gifter.Diamond = diamond
		gifter.Nickname = nickname
		gifter.HeadURL = headURL
		gifter.UserID = userID
		gifter.Status = 1
		item.Gifters = append(item.Gifters, gifter)
		item.Unlock()
	}
	item.Lock()
	sort.Sort(gifterComp(item.Gifters))
	item.Unlock()
}

//GifterManager ...
type GifterManager struct {
	sync.RWMutex
	LiverGifter map[string]*liverOfGifter
}

var gifterManager *GifterManager

//GetGifterManager ...
func GetGifterManager() *GifterManager {
	return gifterManager
}

func (m *GifterManager) getLiverGifter(liverID string) *liverOfGifter {
	m.RLock()
	item, ok := m.LiverGifter[liverID]
	m.RUnlock()
	if ok {
		return item
	}

	return nil
}

//SendGift ...
func SendGift(userID string, headURL string, nickname string, diamond int64, liverID string) {
	m := gifterManager
	item := m.getLiverGifter(liverID)
	if item == nil {
		item = &liverOfGifter{}
		item.LiverID = liverID
		item.Lookers = make(map[string]*Gifter)
		m.Lock()
		m.LiverGifter[liverID] = item
		m.Unlock()
	}

	if item != nil {
		item.addGift(userID, headURL, nickname, diamond)
	}
	return
}

//AddLooker ...
func AddLooker(userID string, headURL string, nickname string, liverID string) {
	m := gifterManager
	item := m.getLiverGifter(liverID)
	if item == nil {
		item = &liverOfGifter{}
		item.LiverID = liverID
		item.Lookers = make(map[string]*Gifter)
		m.Lock()
		m.LiverGifter[liverID] = item
		m.Unlock()
	}

	if item != nil {
		item.addLooker(userID, headURL, nickname)
	}
}

//LiverOffLine ...
func LiverOffLine(liverID string) {
	m := gifterManager
	item := m.getLiverGifter(liverID)
	if item != nil {
		item.liverOffLine()
		m.Lock()
		delete(m.LiverGifter, liverID)
		m.Unlock()
	}
}

//UserOffline ...
func UserOffline(userID string, liverID string) {
	m := gifterManager
	item := m.getLiverGifter(liverID)
	if item != nil {
		item.playerOffline(userID)
	}
}

//UserOnline ...
func UserOnline(liverID string, userID string) {
	m := gifterManager
	item := m.getLiverGifter(liverID)
	if item != nil {
		item.playerOnline(userID)
	}
}

//GetLiverGifterList ...
func GetLiverGifterList(liverID string, count int, page int, getType int) []*Gifter {
	m := gifterManager
	item := m.getLiverGifter(liverID)
	if item != nil {
		return item.getGifterList(count, page, getType)
	}

	return []*Gifter{}
}

//GetLiverLooker ...
func GetLiverLooker(liverID string, count int, page int) []*Gifter {
	m := gifterManager
	item := m.getLiverGifter(liverID)
	if item != nil {
		return item.getLookerList(count, page)
	}

	return []*Gifter{}
}

//GetLiverLookerCount ...
func GetLiverLookerCount(liverID string) int {
	m := gifterManager
	item := m.getLiverGifter(liverID)
	if item != nil {
		return item.getLookerCount()
	}

	return 0
}

//GetLiverMyRank ...
func GetLiverMyRank(liverID string, userID string) (int, *Gifter) {
	m := gifterManager
	item := m.getLiverGifter(liverID)
	if item != nil {
		return item.myRank(userID)
	}
	return -1, nil
}

type gifterComp []*Gifter

func (p gifterComp) Len() int { return len(p) }
func (p gifterComp) Less(i, j int) bool {
	return p[i].Diamond > p[j].Diamond
}
func (p gifterComp) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type liverStateST struct {
	LiverID      string
	OneDayRate   int
	RemainLength int
	Diamond      int
}

type liverStateComp1 []*liverStateST

func (p liverStateComp1) Len() int { return len(p) }
func (p liverStateComp1) Less(i, j int) bool {
	return p[i].OneDayRate > p[j].OneDayRate
}
func (p liverStateComp1) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type liverStateComp2 []*liverStateST

func (p liverStateComp2) Len() int { return len(p) }
func (p liverStateComp2) Less(i, j int) bool {
	return p[i].RemainLength > p[j].RemainLength
}
func (p liverStateComp2) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type liverStateComp3 []*liverStateST

func (p liverStateComp3) Len() int { return len(p) }
func (p liverStateComp3) Less(i, j int) bool {
	return p[i].Diamond > p[j].Diamond
}
func (p liverStateComp3) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

//GetRecommondLiverByOneDayRate ...
func GetRecommondLiverByOneDayRate() ([]string, []string, []string) {
	list1 := []*liverStateST{}
	list2 := []*liverStateST{}
	list3 := []*liverStateST{}

	last1Day := util.GetLastDay()
	items := hall.GetLivingSummaryData2(last1Day, last1Day)
	for i := 0; i < len(items); i++ {
		t := liverStateST{}
		t.LiverID = items[i].LiverID
		t.OneDayRate = items[i].OneDayRate
		if items[i].EnterUser > 0 {
			t.RemainLength = int(items[i].RemainLength / int64(items[i].EnterUser))
		} else {
			t.RemainLength = 0
		}
		t.Diamond = int(items[i].DiamondConsum)

		list1 = append(list1, &t)
		list2 = append(list2, &t)
		list3 = append(list3, &t)
	}

	sort.Sort(liverStateComp1(list1))
	sort.Sort(liverStateComp2(list2))
	sort.Sort(liverStateComp3(list3))

	back1 := []string{}
	back2 := []string{}
	back3 := []string{}

	for i := 0; i < len(list1); i++ {
		back1 = append(back1, list1[i].LiverID)
	}

	for i := 0; i < len(list2); i++ {
		back2 = append(back2, list2[i].LiverID)
	}
	for i := 0; i < len(list3); i++ {
		back3 = append(back3, list3[i].LiverID)
	}
	return back1, back2, back3
}

var forbidManager *ForbidManager

//ForbidManager ...
type ForbidManager struct {
	sync.RWMutex
	liverForbidList map[string]*liverForbid
}

type forbidInfo struct {
	UserID     string
	Time       int64
	InviteTime int64
}

type liverForbid struct {
	liverID    string
	forBidList map[string]*forbidInfo
}

//AddLiverInviteUser ...
func AddLiverInviteUser(liverID, userID string) bool {
	m := forbidManager
	times := 300
	m.Lock()
	defer m.Unlock()
	item, ok := m.liverForbidList[liverID]
	if !ok {
		it := &liverForbid{}
		it.liverID = liverID
		it.forBidList = make(map[string]*forbidInfo)
		userInfo := &forbidInfo{}
		userInfo.UserID = userID
		userInfo.InviteTime = time.Now().Unix() + int64(times)
		userInfo.Time = 0

		m.liverForbidList[liverID] = it
		it.forBidList[userID] = userInfo
		glog.Info("AddLiverInviteUser userID:", userID, ",liverID:", liverID)
		return true
	}

	itUser, itUserOK := item.forBidList[userID]
	if !itUserOK {
		userInfo := &forbidInfo{}
		userInfo.UserID = userID
		userInfo.InviteTime = time.Now().Unix() + int64(times)
		userInfo.Time = 0

		item.forBidList[userID] = userInfo
		glog.Info("AddLiverInviteUser userID:", userID, ",liverID:", liverID)
		return true
	}

	if (time.Now().Unix() - itUser.InviteTime) < 0 {
		glog.Info("AddLiverInviteUser userID:", userID, ",liverID:", liverID)
		return false
	}

	itUser.InviteTime = time.Now().Unix() + int64(times)
	glog.Info("AddLiverInviteUser userID:", userID, ",liverID:", liverID)
	return true
}

//AddLiverForbidUser ...
func AddLiverForbidUser(liverID, userID string, times int) bool {
	m := forbidManager
	m.Lock()
	defer m.Unlock()
	item, ok := m.liverForbidList[liverID]
	if !ok {
		it := &liverForbid{}
		it.liverID = liverID
		it.forBidList = make(map[string]*forbidInfo)
		userInfo := &forbidInfo{}
		userInfo.UserID = userID
		userInfo.Time = time.Now().Unix() + int64(times)
		userInfo.InviteTime = 0

		it.forBidList[userID] = userInfo
		m.liverForbidList[liverID] = it

		return true
	}

	itUser, itUserOK := item.forBidList[userID]
	if !itUserOK {
		userInfo := &forbidInfo{}
		userInfo.UserID = userID
		userInfo.Time = time.Now().Unix() + int64(times)
		userInfo.InviteTime = 0

		item.forBidList[userID] = userInfo
		return true
	}

	if (time.Now().Unix() - itUser.Time) < 0 {
		return false
	}

	itUser.Time = time.Now().Unix() + int64(times)
	return true
}

//GetLiverForbidUserInfo ...
func GetLiverForbidUserInfo(liverID, userID string) (bool, int) {
	m := forbidManager
	m.RLock()

	item, ok := m.liverForbidList[liverID]
	if !ok {
		m.RUnlock()
		return false, 0
	}

	itUser, itUserOK := item.forBidList[userID]
	if !itUserOK {
		m.RUnlock()
		return false, 0
	}

	if (time.Now().Unix() - itUser.Time) > 0 {
		m.RUnlock()
		return false, 0
	}

	left := itUser.Time - time.Now().Unix()
	m.RUnlock()
	return true, int(left)
}

//RemoveLiverAllUser ...
func RemoveLiverAllUser(liverID string) {
	m := forbidManager
	m.Lock()
	defer m.Unlock()

	item, ok := m.liverForbidList[liverID]
	if !ok {
		return
	}

	for userID := range item.forBidList {
		delete(item.forBidList, userID)
	}

	delete(m.liverForbidList, liverID)

	return
}
