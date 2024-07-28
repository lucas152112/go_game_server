package user

import (
	"github.com/golang/glog"
	"sync"
	"time"
)

type Ip_Info struct {
	IP        string
	UserId    string
	LoginTime int64
}

type IPManager struct {
	sync.RWMutex
	items       map[string]*Ip_Info
	whiteIpList map[string]int
	isLimit     bool
}

const (
	IP_LIMIT_TIME = 60 * 30
	IP_CHECK_TIME = 60 * 5
)

var ipManager *IPManager

func init() {
	glog.Info("IpManager init in")
	ipManager = &IPManager{}
	ipManager.items = make(map[string]*Ip_Info)
	ipManager.whiteIpList = make(map[string]int)
	ipManager.isLimit = false
	go ipHeartBeat()
}

func (m *IPManager) Init() {
	lists := GetWhiteIP()
	for _, v := range lists {
		ip := v.Ip
		m.whiteIpList[ip] = 1
	}
}

func GetIpManager() *IPManager {
	return ipManager
}

func (m *IPManager) checkLoginTime() {
	m.Lock()
	defer m.Unlock()

	cur := time.Now().Unix()

	for userId, v := range m.items {
		t := v.LoginTime
		if (cur - t) > int64(IP_LIMIT_TIME) {
			delete(m.items, userId)
		}
	}
}

func ipHeartBeat() {
	for {
		time.Sleep(IP_CHECK_TIME * time.Second)
		ipManager.checkLoginTime()
	}
}

func (m *IPManager) IsCanLogin(ip string, userId string) bool {
	if !m.isLimit {
		return true
	}

	m.RLock()
	defer m.RUnlock()

	_, okW := m.whiteIpList[ip]
	if okW {
		return true
	}

	cur := time.Now().Unix()
	v, ok := m.items[ip]
	if !ok {
		t := Ip_Info{}
		t.IP = ip
		t.UserId = userId
		t.LoginTime = cur
		m.items[ip] = &t
		return true
	} else {
		if v != nil {
			if (cur - v.LoginTime) > int64(IP_LIMIT_TIME) {
				v.LoginTime = cur
				return true
			} else {
				if v.UserId == userId {
					return true
				} else {
					return false
				}
			}
		} else {
			t := Ip_Info{}
			t.IP = ip
			t.UserId = userId
			t.LoginTime = cur
			m.items[ip] = &t
			return true
		}
	}
}

func (m *IPManager) AddWhiteIP(ip string) bool {
	m.RLock()
	defer m.RUnlock()

	m.whiteIpList[ip] = 1
	SetWhiteIP(ip)
	return true
}

func (m *IPManager) ChangeIPLimitStatus(isLimit int) {
	if isLimit == 1 {
		m.isLimit = true
	} else {
		m.isLimit = false
	}
}
