package user

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"game/util"
)

type USER_WE_CHAT struct {
	OpenId       string `bson:"openId"` //openid
	AccessToken  string `bson:"accessToken"`
	RefreshToken string `bson:"refreshToken"`
	AccessTime   int64  `bson:"accessTime"`
	FreshTime    int64  `bson:"freshTime"`
}

const (
	user_we_chat_c = "user_we_chat"
)

func FindWeChatByOpenId(openId string) (*USER_WE_CHAT, error) {
	u := &USER_WE_CHAT{}
	err := util.WithUserCollection(user_we_chat_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"openId": openId}).One(u)
	})
	return u, err
}

func SaveWeChatByOpenId(openId string, access string, refresh string, accessTime int64, freshTime int64) error {
	u := &USER_WE_CHAT{}
	u.OpenId = openId
	u.AccessToken = access
	u.RefreshToken = refresh
	u.AccessTime = accessTime
	u.FreshTime = freshTime
	return util.WithSafeUserCollection(user_we_chat_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"openId": u.OpenId}, u)
		return err
	})
}
