package user

import (
	"game/pb"
	"game/util"
	"time"

	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MatchRecord struct {
	UserId                string    `bson:"userId"`
	WinTimes              int       `bson:"winTimes"`
	LoseTimes             int       `bson:"loseTimes"`
	CurDayEarnGold        int       `bson:"curDayEarnGold"`
	CurWeekEarnGold       int       `bson:"curWeekEarnGold"`
	MaxCards              []int     `bson:"maxCards"`
	DayEarnGoldResetTime  time.Time `bson:"dayEarnGoldResetTime"`
	WeekEarnGoldResetTime time.Time `bson:"weekEarnGoldResetTime"`
	CurMonWin             int       `bson:"curMonWin"`
	MonWinRetTime         time.Time `bson:"monWinRetTime"`
}

func (r *MatchRecord) BuildMessage() *pb.NUserMatchRecordDef {
	r.resetEarnGold()

	msg := &pb.NUserMatchRecordDef{}
	msg.PlayWinCount = r.WinTimes
	msg.PlayTotalCount = r.WinTimes + r.LoseTimes
	msg.TheDayWinGold = r.CurDayEarnGold
	msg.TheWeekWinGold = r.CurWeekEarnGold
	msg.MaxCards = r.MaxCards

	return msg
}

func (r *MatchRecord) resetEarnGold() {
	now := time.Now()
	if !util.CompareDate(now, r.DayEarnGoldResetTime) {
		r.CurDayEarnGold = 0
		r.DayEarnGoldResetTime = now
	}

	if !util.CompareDate(now, r.WeekEarnGoldResetTime) && now.Weekday() == time.Monday {
		r.CurWeekEarnGold = 0
		r.WeekEarnGoldResetTime = now
	}
}

const (
	matchRecordC = "match_record"
)

func FindMatchRecord(userId string) *MatchRecord {
	r := &MatchRecord{}
	r.UserId = userId
	matchRecordTableC, er := GetUserMatchRecordTable(userId)
	if er != nil {
		glog.Error("FindMatchRecord GetUserMatchRecordTable err")
	}
	err := util.WithUserCollection(matchRecordTableC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(r)
	})
	if err != nil {
		glog.Error("加载玩家比赛记录失败err:", err)
	}
	return r
}

func SaveMatchRecord(r *MatchRecord) error {
	matchRecordTableC, er := GetUserMatchRecordTable(r.UserId)
	if er != nil {
		return er
	}

	return util.WithSafeUserCollection(matchRecordTableC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": r.UserId}, r)
		return err
	})
}
