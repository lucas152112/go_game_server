package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

var config *Config
var confName string

func init() {
	flag.StringVar(&confName, "config", "config", "Project operating environment")
}

func LoadConfigFile() error {
	configFile := confName + "." + env + ".json"
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	//glog.Info("config:==>", string(raw))
	c := &Config{}
	err = json.Unmarshal(raw, c)
	if err != nil {
		return err
	}
	config = c
	return nil
}

/**

 */
func Get() *Config {
	return config
}
