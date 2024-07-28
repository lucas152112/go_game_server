package config

import (
	"errors"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"sync"
	"time"
	"game/util"
)

type CardConfig struct {
	GameType    int `bson:"gameType"`
	Single      int `bson:"single"`
	Double      int `bson:"double"`
	ShunZi      int `bson:"shuZi"`
	JinHua      int `bson:"jinHua"`
	ShunJin     int `bson:"shunJin"`
	BaoZi       int `bson:"baoZi"`
	Special235  int `bson:"special235"`
	WinGold     int `bson:"winGold"`
	LoseGold    int `bson:"loseGold"`
	WinRateHigh int `bson:"winRateHigh"`
	WinRateLow  int `bson:"winRateLow"`
	total       int
}

func (c *CardConfig) calcTotal() int {
	c.total = c.Single + c.Double + c.ShunZi + c.JinHua + c.ShunJin + c.BaoZi
	return c.total
}

const (
	cardConfigC = "card_config"
)

const (
	CARD_TYPE_SINGLE   = 1 // 单牌类型
	CARD_TYPE_DOUBLE   = 2 // 对子类型
	CARD_TYPE_SHUN_ZI  = 3 // 顺子类型
	CARD_TYPE_JIN_HUA  = 4 // 金花类型
	CARD_TYPE_SHUN_JIN = 5 // 同花顺类型(顺金)
	CARD_TYPE_BAO_ZI   = 6 // 豹子类型
	CARD_TYPE_SPECIAL  = 7 //235
)

func FindCardConfigs() ([]*CardConfig, error) {
	configs := []*CardConfig{}
	err := util.WithUserCollection(cardConfigC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&configs)
	})

	return configs, err
}

func SaveCardConfig(config *CardConfig) error {
	return util.WithSafeUserCollection(cardConfigC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"gameType": config.GameType}, config)
		return err
	})
}

type CardConfigManager struct {
	sync.RWMutex
	configs       map[int]*CardConfig
	defaultConfig *CardConfig
}

var cardConfigM *CardConfigManager

func init() {
	rand.Seed(time.Now().UnixNano())
	cardConfigM = &CardConfigManager{}
	cardConfigM.configs = make(map[int]*CardConfig)

	config := &CardConfig{}
	config.Single = 16380
	config.Double = 3744
	config.ShunZi = 720
	config.JinHua = 800
	config.ShunJin = 48
	config.BaoZi = 52
	config.Special235 = 60
	config.WinGold = 100000
	config.LoseGold = 100000
	config.WinRateHigh = 100
	config.WinRateLow = 0
	config.calcTotal()
	cardConfigM.defaultConfig = config
}

func (m *CardConfigManager) reset() {
	config := &CardConfig{}
	config.Single = 16380
	config.Double = 3744
	config.ShunZi = 720
	config.JinHua = 800
	config.ShunJin = 48
	config.BaoZi = 52
	config.Special235 = 60
	config.WinGold = 100000
	config.LoseGold = 100000
	config.WinRateHigh = 100
	config.WinRateLow = 0
	config.calcTotal()
	cardConfigM.defaultConfig = config
}

func GetCardConfigManager() *CardConfigManager {
	return cardConfigM
}

func (m *CardConfigManager) Init() bool {
	m.Lock()
	defer m.Unlock()

	configs, err := FindCardConfigs()
	if err != nil && err != mgo.ErrNotFound {
		glog.Error(err)
		return false
	}

	for _, config := range configs {
		glog.Info("==>card_config:", config)
		config.calcTotal()
		if config.total == 0 {
			glog.Error("牌型配置错误,概率总和为0 config:", config)
			continue
		}
		m.configs[config.GameType] = config
	}

	return true
}

func (m *CardConfigManager) GetRandCardType(gameType int) int {
	m.RLock()
	defer m.RUnlock()

	if gameType == 0 {
		panic(errors.New("====>赛事类型为0"))
	}

	/*config := m.configs[gameType]
	if config == nil {
		glog.V(2).Info("===>gameType:", gameType, " 找不到牌型配置，使用默认配置")
		config = m.defaultConfig
	}*/
	config := m.defaultConfig

	/*if config.calcTotal() < 133 {
		m.reset()
		config = m.defaultConfig
	}*/

	offset := 0
	r := rand.Int() % config.total
	if r >= offset && r < offset+config.Single {
		//if config.Single > 0 {
		//config.Single -= 1
		return CARD_TYPE_SINGLE
		//}
	}

	offset += config.Single
	if r >= offset && r < offset+config.Double {
		//if config.Double > 0 {
		//config.Double -= 1
		return CARD_TYPE_DOUBLE
		//}
	}

	offset += config.Double
	if r >= offset && r < offset+config.ShunZi {
		//if config.ShunZi > 0 {
		//config.ShunZi -= 1
		return CARD_TYPE_SHUN_ZI
		//}
	}

	offset += config.ShunZi
	if r >= offset && r < offset+config.JinHua {
		//if config.JinHua > 0 {
		//config.JinHua -= 1
		return CARD_TYPE_JIN_HUA
		//}
	}

	offset += config.JinHua
	if r >= offset && r < offset+config.ShunJin {
		//if config.ShunJin > 0 {
		//config.ShunJin -= 1
		return CARD_TYPE_SHUN_JIN
		//}
	}

	offset += config.ShunJin
	if r >= offset && r < offset+config.BaoZi {
		//if config.BaoZi > 0 {
		//config.BaoZi -= 1
		return CARD_TYPE_BAO_ZI
		//}
	}

	offset += config.BaoZi
	if r >= offset && r < offset+config.Special235 {
		//if config.Special235 > 0 {
		//config.Special235 -= 1
		return CARD_TYPE_SPECIAL
		//}
	}

	return CARD_TYPE_SINGLE
}

func (m *CardConfigManager) SetCardConfig(gameType, single, double, shunZi, jinHua, shunJin, baoZi, winGold, loseGold, winRateHigh, winRateLow int) {
	m.Lock()
	defer m.Unlock()
	glog.Info("++++++++++++++++++++++++++++++++SetCardConfig offset =")

	c := &CardConfig{}
	c.GameType = gameType
	c.Single = single
	c.Double = double
	c.ShunZi = shunZi
	c.JinHua = jinHua
	c.ShunJin = shunJin
	c.BaoZi = baoZi
	c.WinGold = winGold
	c.LoseGold = loseGold
	c.WinRateHigh = winRateHigh
	c.WinRateLow = winRateLow
	c.calcTotal()

	if c.total == 0 {
		return
	}

	m.configs[gameType] = c
	SaveCardConfig(c)
	glog.Info("===>保存牌型配置:", c)
}

func (m *CardConfigManager) GetWinRate(gameType int) (int, int) {
	m.RLock()
	defer m.RUnlock()

	c := m.configs[gameType]
	if c == nil {
		glog.V(2).Info("===>gameType:", gameType, " 找不到牌型配置，使用默认配置")
		c = m.defaultConfig
	}

	return c.WinRateHigh, c.WinRateLow
}

func (m *CardConfigManager) GetTimeConfig(gameType int) (int, int) {
	m.RLock()
	defer m.RUnlock()

	c := m.configs[gameType]
	if c == nil {
		c = m.defaultConfig
	}

	return c.WinGold, c.LoseGold
}
