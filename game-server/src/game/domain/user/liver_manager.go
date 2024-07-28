package user

import (
	mRand "math/rand"
	"sync"

	"game/pb"
	"sort"

	"github.com/golang/glog"
)

type LiverST struct {
	Liver     *Liver
	GameID    int
	MatchType int
}

type liverInfoST struct {
	LiverID   string
	HeadURL   string
	Signiture string
	Nickname  string
}

type LiverManager struct {
	sync.RWMutex
	onlineLivers   map[string]*LiverST
	offLineLivers  map[string]*LiverST
	oneDayRateList []string
	remainList     []string
	diamondList    []string
	liverList      map[string]*liverInfoST
}

var liverManager *LiverManager

func init() {
	liverManager = &LiverManager{}
	liverManager.onlineLivers = make(map[string]*LiverST)
	liverManager.offLineLivers = make(map[string]*LiverST)
	liverManager.oneDayRateList = []string{}
	liverManager.remainList = []string{}
	liverManager.diamondList = []string{}
	liverManager.liverList = make(map[string]*liverInfoST)
}

func GetLiverManager() *LiverManager {
	return liverManager
}

func (m *LiverManager) Init() {
	m.Lock()
	errInit, livers := GetAllActiveLivers()
	if errInit == nil {
		for i := 0; i < len(livers); i++ {
			t := livers[i]
			temp := &LiverST{t, 0, 0}
			errConfig, configGame := GetLiverGameConfig(t.LiverID)
			if errConfig != nil {
				continue
			}
			temp.MatchType = configGame.GameType
			m.offLineLivers[t.LiverID] = temp

			u, errU := FindByUserId(t.LiverID)
			if errU == nil && u != nil {
				info := &liverInfoST{t.LiverID, u.PhotoUrl, u.Signiture, u.Nickname}
				m.liverList[t.LiverID] = info
			}
		}
	}

	m.Unlock()

	m.InitLiverState()

	glog.Info("LiverManager Init offline list len ", len(m.offLineLivers))
}

//主播上播
func (m *LiverManager) LiverOnline(liverID string, gameID int, matchType int) {
	glog.Info("LiverOnline liverID:", liverID, " gameID:", gameID)
	m.Lock()
	defer m.Unlock()

	stOnline, okOnline := m.onlineLivers[liverID]
	if okOnline {
		stOnline.Liver.Status = 1
		stOnline.GameID = gameID
		stOnline.MatchType = matchType
		return
	}

	st, ok := m.offLineLivers[liverID]
	if ok {
		st.Liver.Status = 1
		st.GameID = gameID
		st.MatchType = matchType
		m.onlineLivers[liverID] = st

		delete(m.offLineLivers, liverID)
	}

	return
}

//主播下播
func (m *LiverManager) LiverOffline(liverID string) {
	glog.Info("LiverOffline liverID:", liverID)
	m.Lock()
	defer m.Unlock()

	stOnline, okOnline := m.onlineLivers[liverID]
	if okOnline {
		stOnline.Liver.Status = 2
		stOnline.GameID = 0
		m.offLineLivers[liverID] = stOnline
		delete(m.onlineLivers, liverID)
	}

	go LiverOffLine(liverID)

	go RemoveLiverAllUser(liverID)

	return
}

//删除主播
func (m *LiverManager) LiverDelete(liverID string) {
	m.Lock()
	defer m.Unlock()

	_, okOnline := m.onlineLivers[liverID]
	if okOnline {
		delete(m.onlineLivers, liverID)
		return
	}

	_, ok := m.offLineLivers[liverID]
	if ok {
		delete(m.offLineLivers, liverID)
	}

	return
}

//添加主播
func (m *LiverManager) LiverAdd(liverID string, cover string) {
	m.Lock()
	defer m.Unlock()

	item, okOnline := m.onlineLivers[liverID]
	if okOnline {
		item.Liver.Cover = cover
		return
	}

	liver := &Liver{}
	liver.LiverID = liverID
	liver.Cover = cover
	liver.Status = 2

	liverST := &LiverST{}
	liverST.Liver = liver
	liverST.GameID = 0
	liverST.MatchType = 0

	m.offLineLivers[liverID] = liverST

	return
}

//SetLiverNotice ...
func (m *LiverManager) SetLiverNotice(liverID string, notice string) {
	m.RLock()
	defer m.RUnlock()

	item, okOnline := m.onlineLivers[liverID]
	if okOnline {
		item.Liver.Notice = notice
	} else {
		it, okOffLine := m.offLineLivers[liverID]
		if okOffLine {
			it.Liver.Notice = notice
		}
	}

	go SetLiverNotice(liverID, notice)
}

//GetLiverNotice ...
func (m *LiverManager) GetLiverNotice(liverID string) string {
	m.RLock()
	defer m.RUnlock()

	item, okOnline := m.onlineLivers[liverID]
	if okOnline {
		return item.Liver.Notice
	} else {
		it, okOffLine := m.offLineLivers[liverID]
		if okOffLine {
			return it.Liver.Notice
		}
	}

	return ""
}

//UpdateLiverInfo ...
func (m *LiverManager) UpdateLiverInfo(liverID string, headURL string, signiture string, nickname string) {
	m.Lock()
	defer m.Unlock()
	info, okInfo := m.liverList[liverID]
	if okInfo {
		info.HeadURL = headURL
		info.LiverID = liverID
		info.Nickname = nickname
		info.Signiture = signiture

		return
	}

	infoST := &liverInfoST{liverID, headURL, signiture, nickname}
	m.liverList[liverID] = infoST
	return
}

//GetLiverInfo ...
func (m *LiverManager) GetLiverInfo(liverID string) *pb.FlowInfo {
	m.RLock()
	defer m.RUnlock()
	info, okInfo := m.liverList[liverID]
	if okInfo {
		item := &pb.FlowInfo{}
		item.LiverID = info.LiverID
		item.HeadURL = info.HeadURL
		item.Nickname = info.Nickname
		item.Signiture = info.Signiture

		return item
	}

	return nil
}

func (m *LiverManager) UpateLiverGameType(liverID string, gameType int) {
	itemOff, okOff := m.offLineLivers[liverID]
	if okOff {
		itemOff.MatchType = gameType
		return
	}

	item, okOnline := m.onlineLivers[liverID]
	if okOnline {
		item.MatchType = gameType
		return
	}
}

//获取一个在线的主播，当ID为空时，随机获取一个，当ID非空时，随机获取一个非ID的主播
func (m *LiverManager) GetOneOnlineLiver(gameID int) *LiverST {
	m.RLock()
	defer m.RUnlock()

	if len(m.onlineLivers) == 0 {
		glog.Info("GetOneOnlineLiver list len = 0")
		return nil
	}

	if gameID == 0 {
		index := 0
		ln := len(m.onlineLivers)
		if ln > 0 {
			index = mRand.Intn(ln)
		} else {
			glog.Info("GetOneOnlineLiver list len = 0")
			return nil
		}

		if index >= 0 {
			t := 0
			for _, item := range m.onlineLivers {
				if t == index {
					return item
				}
				t++
			}
		}
	} else {
		index := 0
		fIndex := -1
		for _, item := range m.onlineLivers {
			if gameID == item.GameID {
				fIndex = index + 1
				if fIndex == len(m.onlineLivers) {
					fIndex = 0
					break
				}
			}

			if fIndex == index {
				return item
			}

			index++
		}

		if fIndex == 0 {
			for _, item := range m.onlineLivers {
				return item
			}
		}
	}

	glog.Info("GetOneOnlineLiver out")

	return nil
}

//获取所有主播
func (m *LiverManager) GetLivers() ([]*LiverST, string) {
	m.RLock()
	defer m.RUnlock()

	list := []*LiverST{}
	count := 0
	for _, item := range m.onlineLivers {
		list = append(list, item)
		count++
		if count > 100 {
			break
		}
	}

	sort.Sort(LiverSTSlice(list))

	index := 100
	var recommondID string
	recommondID = ""
	for i := 0; i < len(list); i++ {
		id := list[i].Liver.LiverID
		for j := 0; j < len(m.oneDayRateList); j++ {
			if id == m.oneDayRateList[j] {
				if j < index {
					index = j
					break
				}
			}
		}
	}

	if index < len(m.oneDayRateList) {
		recommondID = m.oneDayRateList[index]
	} else {
		if len(list) > 0 {
			recommondID = list[0].Liver.LiverID
		} else {
			recommondID = ""
		}
	}

	if len(m.onlineLivers) == 0 {
		recommondID = ""
	}

	if count > 100 {
		return list, recommondID
	}

	tempList := []*LiverST{}
	for _, item := range m.offLineLivers {
		tempList = append(tempList, item)
		count++
		if count > 100 {
			break
		}
	}

	sort.Sort(LiverSTSlice(tempList))

	for i := 0; i < len(tempList); i++ {
		list = append(list, tempList[i])
	}

	return list, recommondID
}

type LiverSTSlice []*LiverST

func (p LiverSTSlice) Len() int           { return len(p) }
func (p LiverSTSlice) Less(i, j int) bool { return p[i].Liver.LiverID > p[j].Liver.LiverID }
func (p LiverSTSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

//InitLiverState ...
func (m *LiverManager) InitLiverState() {
	m.Lock()
	defer m.Unlock()
	l1, l2, l3 := GetRecommondLiverByOneDayRate()
	m.oneDayRateList = []string{}
	m.oneDayRateList = l1
	m.remainList = []string{}
	m.remainList = l2
	m.diamondList = []string{}
	m.diamondList = l3
}

//LiverIsOnline ...
func (m *LiverManager) LiverIsOnline(liverID string) bool {
	m.RLock()
	defer m.RUnlock()

	_, ok := m.onlineLivers[liverID]
	if ok {
		return true
	}
	return false
}
