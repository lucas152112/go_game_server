package user

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"game/pb"
	"time"
	"game/util"
)

const (
	global_info_c            = "global_info"
	global_daily_info_c      = "global_daily_info"
	user_daily_info_log_c    = "user_daily_info_log"
	user_jqk_change_log_c    = "user_jqk_change_log"
	user_jqk_fz_change_log_c = "user_jqk_fz_change_log"
)

//全局每日数据
type Global_Daily_Info struct {
	Date          string  `bson:"date"`
	MaxJQKPrice   float32 `bson:"maxJQKPrice"`   //今日JQK成交最高价
	MinJQKPrice   float32 `bson:"minJQKPrice"`   //今日JQK成交最低价
	NewClubNum    int     `bson:"newClubNum"`    //今日新增俱乐部数量
	JQKBonus      float32 `bson:"jQKBonus"`      //今日JQK分红
	JQKTotal      float32 `bson:"jQKTotal"`      //分红时所有玩家的JQK之和
	JQK10kBonus   float32 `bson:"jQK10kBonus"`   //今日JQK万币分红
	RobotCoinsChg int     `bson:"robotCoinsChg"` //今日机器人金币输赢
}

func AddGlobalDailyInfoRobotCoinsChg(robotCoinsChg int) error {
	date := util.GetCurrentDate()
	return util.WithGameCollection(global_daily_info_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"date": date}, bson.M{"$inc": bson.M{"robotCoinsChg": robotCoinsChg}})
		return err
	})
}

func AddGlobalDailyInfoNewClubNum() error {
	date := util.GetCurrentDate()
	return util.WithGameCollection(global_daily_info_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"date": date}, bson.M{"$inc": bson.M{"newClubNum": 1}})
		return err
	})
}

func UpdateGlobalDailyInfoJQKBonus(date string, JQKBonus float32, JQKTotal float32, JQK10kBonus float32) error {
	return util.WithGameCollection(global_daily_info_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"date": date}, bson.M{"$set": bson.M{"jQKBonus": JQKBonus, "jQKTotal": JQKTotal, "jQK10kBonus": JQK10kBonus}})
		return err
	})
}

func UpdateGlobalDailyInfoMaxJQKPrice(maxJQKPrice float32) error {
	date := util.GetCurrentDate()
	return util.WithGameCollection(global_daily_info_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"date": date}, bson.M{"$set": bson.M{"maxJQKPrice": maxJQKPrice}})
		return err
	})
}

func UpdateGlobalDailyInfoMinJQKPrice(minJQKPrice float32) error {
	date := util.GetCurrentDate()
	return util.WithGameCollection(global_daily_info_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"date": date}, bson.M{"$set": bson.M{"minJQKPrice": minJQKPrice}})
		return err
	})
}

//获取全局每日信息
func GetGlobalDailyInfo(date string) (*Global_Daily_Info, error) {
	item := &Global_Daily_Info{}
	err := util.WithGameCollection(global_daily_info_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"date": date}).One(item)
	})
	return item, err
}

//全局数据
type Global_Info struct {
	PriKey           string    `bson:"priKey"`
	BetJQK           float32   `bson:"betJQK"`           //下注挖矿JQK数量
	InviteJQK        float32   `bson:"inviteJQK"`        //邀请挖矿JQK数量
	JQKBonusPool     float32   `bson:"jQKBonusPool"`     //JQK分红奖池
	TotalJQKBonus    float32   `bson:"totalJQKBonus"`    //总共JQK分红
	LastJQKMineTime  time.Time `bson:"lastJQKMineTime"`  //最近一次JQK挖矿时间
	LastJQKBonusTime time.Time `bson:"lastJQKBonusTime"` //最近一次JQK分红时间
	MaxJQK10kBonus   float32   `bson:"maxJQK10kBonus"`   //JQK最高万币分红
	InSurePool       pb.Currency   `bson:"inSurePool"`       //保险池
	ActiveJQK        float32   `bson:"activeJQK"`        //活动挖矿JQK数量
}

//更新最高万币分红
func UpdateGlobalMaxJQK10kBonus(maxJQK10kBonus float32) error {
	return util.WithGameCollection(global_info_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"priKey": "globalinfo"}, bson.M{"$set": bson.M{"maxJQK10kBonus": maxJQK10kBonus}})
		return err
	})
}

//更新JQK分红
func UpdateGlobalJQKBonus(jQKBonusChg float32) error {
	return util.WithGameCollection(global_info_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"priKey": "globalinfo"}, bson.M{"$set": bson.M{"lastJQKBonusTime": time.Now()}, "$inc": bson.M{"totalJQKBonus": jQKBonusChg}})
		return err
	})
}

//更新JQK挖矿
func UpdateGlobalJQK(betJQKChg float32, inviteJQKChg float32) error {
	return util.WithGameCollection(global_info_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"priKey": "globalinfo"}, bson.M{"$inc": bson.M{"betJQK": betJQKChg, "inviteJQK": inviteJQKChg}, "$set": bson.M{"lastJQKMineTime": time.Now()}})
		return err
	})
}

//保险池
func AddGlobalInSurePool(inSurePoolChg pb.Currency) error {
	return util.WithGameCollection(global_info_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"priKey": "globalinfo"}, bson.M{"$inc": bson.M{"inSurePool": inSurePoolChg}})
		return err
	})
}

//活动JQK挖矿
func AddGlobalActiveJQK(activeJQKChg float32) error {
	return util.WithGameCollection(global_info_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"priKey": "globalinfo"}, bson.M{"$inc": bson.M{"activeJQK": activeJQKChg}})
		return err
	})
}

//获取全局信息
func GetGlobalInfo() (*Global_Info, error) {
	item := &Global_Info{}
	err := util.WithGameCollection(global_info_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"priKey": "globalinfo"}).One(item)
	})
	return item, err
}

//每日用户信息
type User_Daily_Info_Log struct {
	UserId          string  `bson:"userId"`
	BetValue        pb.Currency `bson:"betValue"`        //下注额
	MineJQK         float32 `bson:"mineJQK"`         //挖矿JQK数量
	JQKBonus        float32 `bson:"jQKBonus"`        //JQK分红
	FreeChargeTimes int     `bson:"freeChargeTimes"` //破产补助次数
}

//更新下注额
func AddUserBet(userId string, betValueChg pb.Currency) error {
	cur_C := user_daily_info_log_c + "_" + util.GetCurrentDate()
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": userId}, bson.M{"$inc": bson.M{"betValue": betValueChg}})
		return err
	})
}

//加减挖矿JQK
func AddUserMineJQK(userId string, mineJQKChg float32) error {
	cur_C := user_daily_info_log_c + "_" + util.GetLastDay2()
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": userId}, bson.M{"$inc": bson.M{"mineJQK": mineJQKChg}})
		return err
	})
}

//加减JQK分红
func AddUserJQKBonus(userId string, jQKBonusChg float32) error {
	cur_C := user_daily_info_log_c + "_" + util.GetLastDay2()
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": userId}, bson.M{"$inc": bson.M{"jQKBonus": jQKBonusChg}})
		return err
	})
}

func AddUserFreeChargeTimes(userId string, freeChargeTimesChg int) error {
	cur_C := user_daily_info_log_c + "_" + util.GetCurDay()
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": userId}, bson.M{"$inc": bson.M{"freeChargeTimes": freeChargeTimesChg}})
		return err
	})
}

//押注升序排列获取所有的数据
func GetUsersDailyInfo(date string) ([]*User_Daily_Info_Log, error) {
	cur_C := user_daily_info_log_c + "_" + date
	items := []*User_Daily_Info_Log{}
	err := util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Find(bson.M{"betValue": bson.M{"$gt": 0}}).Sort("betValue").All(&items)
	})
	return items, err
}

func GetUserDailyInfo(userId string, date string) (*User_Daily_Info_Log, error) {
	cur_C := user_daily_info_log_c + "_" + date
	item := &User_Daily_Info_Log{}
	err := util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(item)
	})
	return item, err
}

//权益币JQK变化记录
type JQKChangeLog struct {
	UserId  string    `bson:"userId"`
	Reason  string    `bson:"reason"`
	Value   float32   `bson:"value"`
	Result  float32   `bson:"result"`
	TimeStr string    `bson:"timeStr"`
	Time    time.Time `bson:"time"`
}

func SaveJQKChangeLog(log *JQKChangeLog) error {
	return util.WithLogCollection(user_jqk_change_log_c, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}

//冻结权益币JQKFrozen变化记录
type JQKFrozenChangeLog struct {
	UserId  string    `bson:"userId"`
	Reason  string    `bson:"reason"`
	Value   float32   `bson:"value"`
	Result  float32   `bson:"result"`
	TimeStr string    `bson:"timeStr"`
	Time    time.Time `bson:"time"`
}

func SaveJQKFrozenChangeLog(log *JQKFrozenChangeLog) error {
	return util.WithLogCollection(user_jqk_fz_change_log_c, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}
