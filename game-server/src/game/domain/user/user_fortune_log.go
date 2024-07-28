package user

import (
	mgo "gopkg.in/mgo.v2"
	"time"
	"game/util"
)

type FortuneLog struct {
	UserId     string    `bson:"userId"`
	Gold       int       `bson:"gold"`
	CurGold    int64     `bson:"curGold"`
	Diamond    int       `bson:"diamond"`
	CurDiamond int       `bson:"curDiamond"`
	ThirdCurrency int	 `bson:"thirdCurrency"`
	CurThirdCurrency int `bson:"curThirdCurrency"`
	Score      int       `bson:"score"`
	CurScore   int       `bson:"curScore"`
	Reason     string    `bson:"reason"`
	DBTime     string    `bson:"dbTime"`
	Time       time.Time `bson:"time"`
}

const (
	earnFortuneLogC    = "earn_fortune_log"
	consumeFortuneLogC = "consume_fortune_log"
)

func SaveEarnFortuneLog(log *FortuneLog) error {
	cur_C := earnFortuneLogC + "_" + util.GetCurrentDate()
	log.DBTime = util.GetDbTime()
	log.Time = time.Now()
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}

func SaveConsumeFortuneLog(log *FortuneLog) error {
	cur_C := consumeFortuneLogC + "_" + util.GetCurrentDate()
	log.DBTime = util.GetDbTime()
	log.Time = time.Now()
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}
