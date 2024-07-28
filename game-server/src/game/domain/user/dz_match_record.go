package user

import (
	"game/pb"
	"game/util"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DZMatchRecord struct {
	UserId    string `bson:"userId"`
	Game      int    `bson:"game"`      //总局数(好友场)
	Play      int    `bson:"play"`      //总手数(好友场、金币场)
	Win       int    `bson:"win"`       //赢手数
	NoMang    int    `bson:"noMang"`    //不是大小盲手数
	NoMangBet int    `bson:"noMangBet"` //不是大小盲(且翻牌前加注或跟注)手数
	AllIn     int    `bson:"allIn"`     //全下手数
	AllInWin  int    `bson:"allInWin"`  //全下赢手数
	RaisePre  int    `bson:"raisePre"`  //翻牌前加注手数
	Flop      int    `bson:"flop"`      //flop的手数(翻牌)
	Bet3      int    `bson:"bet3"`      //3-Bet的手数(翻牌前除去大小盲后的第二次加注)
	BetC      int    `bson:"betC"`      //C-Bet手数(持续加注，翻牌前加注，翻牌后也加注)
}

func FindDZMatchRecord(userId string) *DZMatchRecord {
	record := &DZMatchRecord{}
	record.UserId = userId

	matchRecord_dz_c := GetUserDZMatchRecordTable(userId)

	util.WithGameCollection(matchRecord_dz_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(record)
	})

	return record
}

func SaveDZMatchRecord(record *DZMatchRecord) error {
	if record == nil {
		return nil
	}
	matchRecord_dz_c := GetUserDZMatchRecordTable(record.UserId)

	return util.WithGameCollection(matchRecord_dz_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": record.UserId}, record)
		return err
	})
}

func (r *DZMatchRecord) BuildMessage() *pb.DZMatchRecordDef {
	msg := &pb.DZMatchRecordDef{}
	msg.GameTimes = r.Game
	msg.PlayTimes = r.Play
	if r.NoMang != 0 {
		f := int(float32(r.NoMangBet) / float32(r.NoMang) * 100)
		msg.VPIP = f
	}

	if r.Play != 0 {
		f := int(float32(r.Win) / float32(r.Play) * 100)
		msg.WinRate = f
	}

	return msg
}
