package user

import (
	"fmt"

	//"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	//"time"
	"game/util"
)

//用户统计，根据phone_user 执行统计

/**

 */
func SummaryAllUser() (int, error) {
	Num := 0
	err := util.WithUserCollection(phoneUserC, func(c *mgo.Collection) error {
		n, e := c.Find(bson.M{}).Count()
		if e == nil {
			Num = n
		}
		return e
	})
	return Num, err
}

//获取当天的
func SummaryTodayUser() (int, error) {
	Num := 0
	tm := time.Now()
	date := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
	date_f, _ := time.Parse("2006-01-02 15:04:05", date+" 00:00:00")
	date_t, _ := time.Parse("2006-01-02 15:04:05", date+" 23:59:59")
	err := util.WithUserCollection(phoneUserC, func(c *mgo.Collection) error {
		n, e := c.Find(bson.M{"bindTime": bson.M{"$gt": date_f, "$lt": date_t}}).Count()
		if e == nil {
			Num = n
		}
		return e
	})
	if err != nil {
		print(err.Error())
	}
	return Num, err
}

/**

 */
/**
今日金币消耗
*/
func GetToday() (time.Time, time.Time) {
	tm := time.Now()
	date := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
	date_f, _ := time.Parse("2006-01-02 15:04:05", date+" 00:00:00")
	date_t, _ := time.Parse("2006-01-02 15:04:05", date+" 23:59:59")
	return date_f, date_t
}

func GetTodayStr1() (string, string) {
	Begin, end := GetToday()
	return Begin.Format("20060102150405"), end.Format("20060102150405")
}
func GetTodayStr2() (string, string) {
	Begin, end := GetToday()
	return Begin.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05")
}

func GetPvCoinConsume() int64 {
	return 0
}

func GetPvBalanceConsume() int64 {
	return 0
}
