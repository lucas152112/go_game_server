package user

import (
	"fmt"
	"game/util"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserLog struct {
	UserId             string    `bson:"userId"`
	UserName           string    `bson:"userName"`
	TotalOnlineSeconds int       `bson:"totalOnlineSeconds"`
	MatchTimes         int       `bson:"matchTimes"`
	CreateTime         time.Time `bson:"createTime"`
	Channel            string    `bson:"channel"`
	Model              string    `bson:"model"`
}

type LoginRecord struct {
	UserId     string    `bson:"userId"`
	UserName   string    `bson:"userName"`
	Channel    string    `bson:"channel"`
	LoginTime  time.Time `bson:"loginTime"`
	LogoutTime time.Time `bson:"logoutTime"`
	LoginIP    string    `bson:"loginIP"`
	DeviceId   string    `bson:"deviceId"`
}

type SlowMsgRecord struct {
	UserId     string    `bson:"userId"`
	MsgId      string    `bson:"msgId"`
	StartTime  time.Time `bson:"startTime"`
	ElapseTime string    `bson:"elapseTime"`
}

const (
	userLogC       = "user_log"
	loginRecordC   = "login_record"
	slowMsgRecordC = "slow_msg"
)

func FindUserLog(userId string) (*UserLog, error) {
	l := &UserLog{}
	l.UserId = userId
	err := util.WithLogCollection(userLogC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(l)
	})
	return l, err
}

func SaveUserLog(l *UserLog) error {
	return util.WithLogCollection(userLogC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": l.UserId}, l)
		return err
	})
}

func InsertLoginRecord(r *LoginRecord) error {
	cur_C := loginRecordC + "_" + util.GetCurrentDate()
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(r)
	})
}

func SaveSlowMsg(userId string, msgId string, startTime time.Time, elapseTime string) error {
	msg := &SlowMsgRecord{}
	msg.UserId = userId
	msg.MsgId = msgId
	msg.StartTime = startTime
	msg.ElapseTime = elapseTime

	return util.WithLogCollection(slowMsgRecordC, func(c *mgo.Collection) error {
		return c.Insert(msg)
	})
}

type UserLastLoginInfo struct {
	UserId    string    `bson:"id"`
	LoginTime time.Time `bson:"time"`
	LoginDate string    `bson:"date"`
}

const (
	userLastLoginLogC = "user_Last_login_log"
)

func GetUserLastLoginInfo(userId string) (string, error) {
	l := &UserLastLoginInfo{}
	l.UserId = userId
	err := util.WithLogCollection(userLastLoginLogC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"id": userId}).One(l)
	})
	return l.LoginDate, err
}

func SaveUserLastLoginInfo(userId string) error {
	l := &UserLastLoginInfo{}
	l.UserId = userId
	now := time.Now()
	l.LoginTime = now
	l.LoginDate = now.Format("2006/01/02 15:04:05")
	return util.WithLogCollection(userLastLoginLogC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"id": l.UserId}, l)
		return err
	})
}

type UserFortuneInfo struct {
	UserId  string    `bson:"id"`
	Time    time.Time `bson:"time"`
	Gold    int64     `bson:"gold"`
	Save    int64     `bson:"save"`
	Channel string    `bson:"channel"`
	Charm   int       `bson:"charm"`
}

const (
	userFortuneLogC           = "user_Fortune_log"
	userFirstLoginFortuneLogC = "first_login_Fortune_log"
)

func SaveUserFirstLoginInfo(userId string, gold int64, save int64, channel string, charm int) error {
	l := &UserFortuneInfo{}
	l.UserId = userId
	now := time.Now()
	l.Time = now
	l.Gold = gold
	l.Save = save
	l.Channel = channel
	l.Charm = charm

	cur_C := userFirstLoginFortuneLogC + "_" + util.GetCurrentDate()

	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(l)
	})
}

func SaveUserFortuneInfo(userId string, gold int64, save int64, channel string, charm int) error {
	date, err := GetUserLastLoginInfo(userId)
	if err != nil && err != mgo.ErrNotFound {
		return err
	} else {
		if err == mgo.ErrNotFound {
			go SaveUserLastLoginInfo(userId)
			go SaveUserFirstLoginInfo(userId, gold, save, channel, charm)
		} else {
			now := time.Now()
			curDate := fmt.Sprintf("%v-%v-%v", now.Year(), now.Month(), now.Day())
			if date != curDate {
				go SaveUserLastLoginInfo(userId)
				go SaveUserFirstLoginInfo(userId, gold, save, channel, charm)
			}
		}
	}

	l := &UserFortuneInfo{}
	l.UserId = userId
	now := time.Now()
	l.Time = now
	l.Gold = gold
	l.Save = save
	l.Channel = channel
	l.Charm = charm

	cur_C := userFortuneLogC + "_" + util.GetCurrentDate()

	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(l)
	})
}

type UserLoginLog struct {
	UserId  string    `bson:"id"`
	Time    time.Time `bson:"time"`
	Channel string    `bson:"channel"`
}

const (
	userLoginLogC = "user_login_log"
)

func SaveUserLoginLog(userId string, channel string) error {
	l := &UserLoginLog{}
	l.UserId = userId
	now := time.Now()
	l.Time = now
	l.Channel = channel

	cur_C := userLoginLogC + "_" + util.GetCurrentDate()

	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(l)
	})
}
