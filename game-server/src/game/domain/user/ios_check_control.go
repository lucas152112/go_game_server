package user

import (
	mgo "gopkg.in/mgo.v2"
	"game/util"
)

const (
	ios_check_control_c = "ios_check_control"
)

type Ios_Check_Status struct {
	Ios    string `bson:"ios"`
	Status string `bson:"status"`
}

func GetIosCheckStatus() error {
	data := &Ios_Check_Status{}
	err := util.WithUserCollection(ios_check_control_c, func(c *mgo.Collection) error {
		return c.Find(nil).One(&data)
	})
	return err
}
