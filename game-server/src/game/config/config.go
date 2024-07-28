package config

import (
	"game/conf"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"sync"
	"game/util"
)

const (
	ControlKey = "f1b9df1ed816c76ecfb2acf1c65b2a0d"
	RobotKey   = "11e6f1e3c1f07d0ce71eb59d6391d9f2"
	ActiveKey  = "57fd936746db585bae644e5ff997c74b"

	TurnOnTask   = true     // 开启任务
	TurnOnReplay = true     // 开启回放

	REPLAY_ACCESS_KEY = "1f1bc0d7af166b6f3a41c275217d3954"
	REVIEWS_ACCESS_KEY = "d5632ae29692bb18f7b445f14553a5ef"

	DUMMY_LOOKER_UID = "000000"
)

type Config struct {
	ConfigId   int    `bson:"configId"`
	CurVersion string `bson:"curVersion"`
}

const (
	configC = "config"
)

/**

 */
func GetVersion() string   {
	if conf.Get().Version == "" {
		return "v0.0.0"
	}else{
		return conf.Get().Version
	}

}

func FindConfig() (*Config, error) {
	config := &Config{}
	err := util.WithUserCollection(configC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"configId": 1}).One(config)
	})

	return config, err
}

func SaveConfig(config *Config) error {
	config.ConfigId = 1
	return util.WithSafeUserCollection(configC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"configId": config.ConfigId}, config)
		return err
	})
}

type ConfigManager struct {
	sync.RWMutex
	config *Config
}

var configM *ConfigManager

func init() {
	configM = &ConfigManager{}
}

func GetConfigManager() *ConfigManager {
	return configM
}

func (m *ConfigManager) Init() {
	m.Lock()
	defer m.Unlock()

	config, err := FindConfig()
	if err != nil && err != mgo.ErrNotFound {
		glog.Error("加载服务器配置失败:", err)
	}

	m.config = config
	glog.Info("===>config:", config)
}

func (m *ConfigManager) SetCurVersion(version string) {
	m.Lock()
	defer m.Unlock()

	m.config.CurVersion = strings.Trim(version, " ")

	SaveConfig(m.config)
}

func (m *ConfigManager) GetCurVersion() string {
	m.RLock()
	defer m.RUnlock()

	return m.config.CurVersion
}
