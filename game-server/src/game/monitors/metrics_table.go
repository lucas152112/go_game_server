package monitors

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

var tableGauge  *prometheus.GaugeVec
/**
 gameType      800
 currencyType  2    3    2为积分场、3为金币场
 roomId        1, 2, 3  1是密码局，2是俱乐部，3是直播局
 */


var tableCount prometheus.Gauge

func init()  {
	tableGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:"game_table",
		Help:"当前系统中所有的牌桌信息",
	},[]string{
			"matchType",      // 800
			"roomType",       //
			"currencyType",
	})

	tableCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:"game_table_count",
		Help:"当前游戏中所创建的牌桌总数",
	})

}

func GetTableGauge() *prometheus.GaugeVec  {
	return tableGauge
}

func GetTableCount() prometheus.Gauge {
	return tableCount
}

func GameTableInc( matchType ,currencyType,roomType int )  {
	tableGauge.With(prometheus.Labels{
		"matchType":strconv.Itoa(matchType),
		"currencyType":strconv.Itoa(currencyType),
		"roomType":strconv.Itoa(roomType),
	}).Inc()
}

func GameTableDec( matchType ,currencyType,roomType int )  {
	tableGauge.With(prometheus.Labels{
		"matchType":strconv.Itoa(matchType),
		"currencyType":strconv.Itoa(currencyType),
		"roomType":strconv.Itoa(roomType),
	}).Dec()
}