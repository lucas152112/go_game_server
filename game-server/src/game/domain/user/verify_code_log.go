package user

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"game/util"
)

//验证码
type VerifyCodeLogInfo struct {
	Email      string    `json:"email" bson:"email"`
	ExpireAt   time.Time `json:"expireAt" bson:"expireAt"` //过期时间
	VerifyCode int64     `json:"verifyCode" bson:"verifyCode"`
}

const (
	verify_code_log_c = "verify_code_log"
)

//索引，过期时间：db.verify_code_log.createIndex( { "expireAt": 1 }, { expireAfterSeconds: 0 } )
//玩家战绩ID
func SaveVerifyCodeLogInfo(email string, verifyCode int64) error {
	log := &VerifyCodeLogInfo{}
	log.Email = email
	log.VerifyCode = verifyCode
	log.ExpireAt = time.Now().UTC().Add(2 * 60 * time.Second)

	return util.WithLogCollection(verify_code_log_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"email": log.Email}, log)
		return err
	})
}

func GetVerifyCodeLogInfoByEmail(email string) (*VerifyCodeLogInfo, error) {
	result := &VerifyCodeLogInfo{}
	err := util.WithLogCollection(verify_code_log_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"email": email}).One(result)
	})
	return result, err
}
