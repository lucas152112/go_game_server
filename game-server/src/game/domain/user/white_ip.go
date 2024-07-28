package user

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"game/util"
)

const (
	White_IP_C = "white_ip"
)

type White_IP_INFO struct {
	Ip string `bson:"ip"`
}

func GetWhiteIP() []*White_IP_INFO {
	var ips []*White_IP_INFO
	util.WithGameCollection(White_IP_C, func(c *mgo.Collection) error {
		return c.Find(nil).All(&ips)
	})

	return ips
}

func SetWhiteIP(ip string) error {
	return util.WithGameCollection(White_IP_C, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"ip": ip}, bson.M{"$set": bson.M{"ip": ip}})
		return err
	})
}
