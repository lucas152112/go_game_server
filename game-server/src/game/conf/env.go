package conf

import "flag"

/**
 运行环境变量
 */

const (
	ENV_DEVELOP      = "develop"          //开发环境
	ENV_TESTING      = "testing"          //测试环境
	ENV_PRODUCTION   = "production"       //生产环境
)

var env string

func init()  {
	flag.StringVar(&env, "env", ENV_PRODUCTION, "Project operating environment")
}


func IsDevelopEnv() bool {
	return env == ENV_DEVELOP
}

func IsProductionEnv() bool  {
	return env == ENV_PRODUCTION
}

func IsTestingEnv() bool  {
	return env == ENV_TESTING
}