package main

import (
	"flag"
	"fmt"
	"game/conf"
	"game/config"
	domainClub "game/domain/dzclub"
	domainDZ "game/domain/dzgame"
	domainHall "game/domain/hall"
	domainLog "game/domain/logServer"
	"game/domain/rankingList"
	"game/domain/stats"
	"game/domain/user"
	domainUser "game/domain/user"
	"game/handlers"
	"game/monitors"
	"game/server"
	"game/util"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"game/domain/core"
	"game/domain/hall"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	flagResult := flag.Lookup("redis")
	glog.Info("redis master addr:", flagResult.Value)

	//加入文件配置
	if err := conf.LoadConfigFile(); err != nil {
		glog.Info(err.Error())
		return
	}

	glog.Info("===>启动游戏服务器")
	loadConfig()
	initialize()

	runtime.GOMAXPROCS(runtime.NumCPU())

	go stopHandler(server.GetServerInstance().GetSigChan())

	//加入监控
	go monitors.Server()

	go util.CleanupVisitors()
	go timerSummary()

	/*go func() {
		for {
			time.Sleep(time.Second)
			readGameOperateToRedis()
		}
	}()*/

	server.GetServerInstance().StartServer(handlers.GetMsgRegistry())
}

func stopHandler(c chan os.Signal) {
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGINT)
	for {
		sig := <-c
		glog.Infof("stopHandler Received signal %s", sig)
		if sig != syscall.SIGINT {
			break
		}
		glog.Infof("===== ALL STACKTRACE =====\n%s\n", util.StackTrace(true))
	}

	domainDZ.GetGameManager().StopServer()
	rankingList.GetRankingList().Save()

	server.GetServerInstance().SetRefuseService()
	server.GetServerInstance().WaitStopServer()
	domainUser.GetPlayerManager().StopServer()
	domainDZ.GetMatchManager().StopServer()
	glog.Flush()
	os.Exit(0)
}

func initialize() {
	rand.Seed(time.Now().UnixNano())
	config.GetConfigManager().Init()
	config.GetCardConfigManager().Init()

	stats.GetMatchLogManager().Init()
	domainUser.GetIpManager().Init()
	user.GetUserFortuneManager().UpdateGoldInGameFunc = nil
	user.GetPlayerManager().Init(viper.GetString("dbtype"), viper.GetString("dbparam"))
	core.InitDB(viper.GetString("dbtype"), viper.GetString("dbparam"))
	domainUser.InitLiverStatus()

	domainClub.GetClubManager().Init()
	domainClub.GetAllianceManager().Init()
	rankingList.GetRankingList().Init()
	domainDZ.GetGameManager().Init()
	domainUser.InitRank()
	domainUser.GetLiverManager().Init()
	domainHall.GetChannelManager().Init()
}

var G_Last_Day = ""

func timerSummary() {
	for {
		if G_Last_Day == "" {
			domainLog.InitTable()
		}
		G_Last_Day = util.GetCurrentDate2()

		time.Sleep(5 * time.Minute)

		if G_Last_Day != util.GetCurrentDate2() {
			domainLog.UpdateTotalData()
			domainLog.DeleteTotalData()
			domainLog.InitTable()
			hall.DayLivingSummary(G_Last_Day)
			domainUser.GetLiverManager().InitLiverState()
		}

		go domainHall.GetChannelManager().Init()
	}
}

func loadConfig() {
	viper.SetDefault("debug", false)
	viper.SetDefault("dbtype", "sqlite3")
	viper.SetDefault("dbparam", "feidao.db")
	viper.SetDefault("gameport", ":9011")
	viper.SetDefault("coinrobot", false)

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.ReadInConfig()

	fmt.Println("debug:", viper.GetBool("debug"))
	fmt.Println("dbconfig dbtype ", viper.GetString("dbtype"))
	fmt.Println("dbconfig dbparam ", viper.GetString("dbparam"))
	fmt.Println("dbconfig gameport ", viper.GetString("gameport"))
	fmt.Println("robotconfig coinrobot ", viper.GetBool("coinrobot"))

	return
}

func readGameOperateToRedis() {
	cli := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6377",
		Password: "",
		DB:       0,
	})

	defer cli.Close()

	val, err := cli.Get("last_operate").Result()
	if err == nil {
		glog.Info("readGameOperateToRedis msg:", val)
	}
}
