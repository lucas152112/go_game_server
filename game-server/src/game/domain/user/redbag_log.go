package user

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"game/util"
)

type RedBagLog struct {
	UserId      string    `bson:"userId"`
	OtherUserId string    `bson:"otherUserId"`
	Count       int       `bson:"count"`
	TimeStr     string    `bson:"timeStr"`
	Time        time.Time `bson:"time"`
	Status      int       `bson:"status"` //1发送， 2 接收
}

const (
	red_bag_log_c = "red_bag_log"
)

func SaveRedBagGiveLog(otherUserId string, userId string, count int) error {
	log := &RedBagLog{}
	log.OtherUserId = otherUserId
	log.UserId = userId
	log.Count = count
	log.TimeStr = util.GetDbTime()
	log.Time = time.Now()
	log.Status = 1
	return util.WithLogCollection(red_bag_log_c, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}

func SaveRedBagReceiveLog(otherUserId string, userId string, count int) error {
	log := &RedBagLog{}
	log.OtherUserId = otherUserId
	log.UserId = userId
	log.Count = count
	log.TimeStr = util.GetDbTime()
	log.Time = time.Now()
	log.Status = 2
	return util.WithLogCollection(red_bag_log_c, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}

func GetRedBagLog(userId string) (error, []*RedBagLog) {
	items := []*RedBagLog{}
	err := util.WithLogCollection(red_bag_log_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).Sort("-timeStr").All(&items)
	})

	return err, items
}

func GetRedBagLogByDate(dateStart string, dateEnd string) (error, []*RedBagLog) {
	items := []*RedBagLog{}
	err := util.WithLogCollection(red_bag_log_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"status": 1, "timeStr": bson.M{"$gte": dateStart, "$lt": dateEnd}}).All(&items)
	})
	return err, items
}
