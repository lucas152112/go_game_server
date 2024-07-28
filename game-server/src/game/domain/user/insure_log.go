package user

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"game/pb"
	"time"
	"game/util"
)

const (
	insure_log_c           = "insure_log"
	INSURE_CHG_TYPE_RESULT = 1
	INSURE_CHG_TYPE_OUT    = 2
)

//保险金额变化
type InSureScoreLogInfo struct {
	UserId     string    `json:"userId" bson:"userId"`         //用户id
	GameName   string    `json:"gameName" bson:"gameName"`     //牌局名称
	ScoreChg   pb.Currency   `json:"scoreChg" bson:"scoreChg"`     //变化金额
	NewScore   pb.Currency   `json:"newScore" bson:"newScore"`     //最新金额
	Date       string    `json:"date" bson:"date"`             //日期字符串
	ChgType    int       `json:"chgType" bson:"chgType"`       //变化类型
	CreateTime int64     `json:"createTime" bson:"createTime"` //创建时间
	ExpireAt   time.Time `json:"expireAt" bson:"expireAt"`     //过期时间
	TableId    int       `json:"tableId" bson:"tableId"`       //桌子id
}

func SaveInSureScoreLogInfo(userId string, tableId int, gameName string, scoreChg , newScore pb.Currency, chgType int) error {
	log := &InSureScoreLogInfo{}
	log.UserId = userId
	log.TableId = tableId
	log.GameName = gameName
	log.ScoreChg = scoreChg
	log.NewScore = newScore
	log.ChgType = chgType
	log.CreateTime = time.Now().Unix()
	log.Date = util.GetDbTime2()
	t, err := util.ParseTime3(log.Date)
	if err == nil {
		log.CreateTime = t.Unix()
	}
	log.ExpireAt = time.Now().UTC().Add(30 * 24 * 60 * 60 * time.Second)
	return util.WithLogCollection(insure_log_c, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}

func GetInSureScoreLogInfoList(userId string, strMaxDate string) ([]*InSureScoreLogInfo, error) {
	items := []*InSureScoreLogInfo{}
	timeStmp := time.Now().Unix()
	glog.Info("GetInSureScoreLogInfoList userId:", userId, ",timeStmp:", timeStmp, ",strMaxDate:", strMaxDate)
	t, err := util.ParseTime3(strMaxDate)
	if err == nil {
		timeStmp = t.Unix()
	}
	glog.Info("GetInSureScoreLogInfoList timeStmp:", timeStmp, t, err)
	err = util.WithLogCollection(insure_log_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId, "chgType": INSURE_CHG_TYPE_RESULT, "createTime": bson.M{"$lt": timeStmp}}).Sort("-createTime").Limit(10).All(&items)
	})
	return items, err
}
