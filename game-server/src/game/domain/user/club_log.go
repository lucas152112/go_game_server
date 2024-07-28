package user

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"game/util"
)

//俱乐部总数据中AdminId为-1
const (
	club_profit_log_c        = "club_profit_log"
	club_player_log_c        = "club_player_log"
	club_total_log_c         = "club_total_log"
	club_member_profit_log_c = "club_member_profit_log"
)

type CLUB_PLAYER_LOG struct {
	ClubId  int    `bson:"clubId"`
	AdminId string `bson:"adminId"`
	UserId  string `bson:"userId"`
}

func AddClubPlayerLog(clubId int, adminId string, userId string) error {
	playerLog := &CLUB_PLAYER_LOG{}
	playerLog.ClubId = clubId
	playerLog.AdminId = adminId
	playerLog.UserId = userId
	cur_C := club_player_log_c + "_" + util.GetCurrentDate()
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"clubId": clubId, "adminId": adminId, "userId": userId}, playerLog)
		return err
	})
}

func GetClubPlayerLogPlayerNum(clubId int, adminId string, date string) (int, error) {
	cur_C := club_player_log_c + "_" + date
	playerNum := 0
	err := util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		count, err := c.Find(bson.M{"clubId": clubId, "adminId": adminId}).Count()
		if err == nil {
			playerNum = count
		}
		return err
	})
	return playerNum, err
}

func GetActiveClubNum(date string) (int, error) {
	result := []int{}
	cur_C := club_player_log_c + "_" + date
	err := util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Find(bson.M{}).Distinct("clubId", &result)
	})
	return len(result), err
}

//俱乐部收益日志
type CLUB_PROFIT_LOG struct {
	ClubId       int     `bson:"clubId"`
	AdminId      string  `bson:"adminId"`
	GamesNum     int     `bson:"gamesNum"`
	TotalTakeUsd float32 `bson:"totalTakeUsd"`
	InSureProfit float32 `bson:"inSureProfit"`
	GameTax      float32 `bson:"gameTax"`
}

func GetClubAdminProfitLog(adminId string, date string) ([]*CLUB_PROFIT_LOG, error) {
	items := []*CLUB_PROFIT_LOG{}
	cur_C := club_profit_log_c + "_" + date
	err := util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Find(bson.M{"adminId": adminId}).All(&items)
	})
	return items, err
}

func GetClubAdminsProfitLog(clubId int, date string) ([]*CLUB_PROFIT_LOG, error) {
	items := []*CLUB_PROFIT_LOG{}
	cur_C := club_profit_log_c + "_" + date
	err := util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Find(bson.M{"clubId": clubId}).All(&items)
	})
	return items, err
}

func GetClubProfitLog(clubId int, adminId string, date string) (*CLUB_PROFIT_LOG, error) {
	item := &CLUB_PROFIT_LOG{}
	cur_C := club_profit_log_c + "_" + date
	err := util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Find(bson.M{"clubId": clubId, "adminId": adminId}).One(item)
	})
	return item, err
}

func AddClubProfitLogTotalTakeUsd(clubId int, adminId string, takeUsd int) error {
	cur_C := club_profit_log_c + "_" + util.GetCurrentDate()
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"clubId": clubId, "adminId": adminId}, bson.M{"$inc": bson.M{"totalTakeUsd": float32(takeUsd)}})
		return err
	})
}

func AddClubProfitLogGamesNum(clubId int, adminId string) error {
	cur_C := club_profit_log_c + "_" + util.GetCurrentDate()
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"clubId": clubId, "adminId": adminId}, bson.M{"$inc": bson.M{"gamesNum": 1}})
		return err
	})
}

func AddClubProfitLogBigGameResult(clubId int, adminId string, inSureProfitChg float32, gameTaxChg float32) error {
	cur_C := club_profit_log_c + "_" + util.GetCurrentDate()
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"clubId": clubId, "adminId": adminId}, bson.M{"$inc": bson.M{"inSureProfit": inSureProfitChg, "gameTax": gameTaxChg}})
		return err
	})
}

type CLUB_TOTAL_LOG struct {
	ClubId       int     `bson:"clubId"`
	GamesNum     int     `bson:"gamesNum"`
	InSureProfit float32 `bson:"inSureProfit"`
	GameTax      float32 `bson:"gameTax"`
	TotalProfit  float32 `bson:"totalProfit"`
}

func GetClubTotalLogSortGamesNum() ([]*CLUB_TOTAL_LOG, error) {
	items := []*CLUB_TOTAL_LOG{}
	err := util.WithLogCollection(club_total_log_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{}).Sort("-gamesNum").Limit(50).All(&items)
	})
	return items, err
}

func GetClubTotalLogSortTotalPrifit() ([]*CLUB_TOTAL_LOG, error) {
	items := []*CLUB_TOTAL_LOG{}
	err := util.WithLogCollection(club_total_log_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{}).Sort("-totalProfit").Limit(50).All(&items)
	})
	return items, err
}

func GetClubTotalLog(clubId int) (*CLUB_TOTAL_LOG, error) {
	item := &CLUB_TOTAL_LOG{}
	err := util.WithLogCollection(club_total_log_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"clubId": clubId}).One(item)
	})
	return item, err
}

func AddClubTotalLogLittleGameResult(clubId int) error {
	return util.WithLogCollection(club_total_log_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"clubId": clubId}, bson.M{"$inc": bson.M{"gamesNum": 1}})
		return err
	})
}

func AddClubTotalLogBigGameResult(clubId int, inSureProfitChg float32, gameTaxChg float32) error {
	return util.WithLogCollection(club_total_log_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"clubId": clubId}, bson.M{"$inc": bson.M{"inSureProfit": inSureProfitChg, "gameTax": gameTaxChg, "totalProfit": inSureProfitChg + gameTaxChg}})
		return err
	})
}

type CLUB_MEMBER_PROFIT_LOG struct {
	ClubId       int     `bson:"clubId"`
	AdminId      string  `bson:"adminId"`
	UserId       string  `bson:"userId"`
	InSureProfit float32 `bson:"inSureProfit"`
	GameTax      float32 `bson:"gameTax"`
	TotalUsdChg  float32 `bson:"totalUsdChg"`
}

func GetClubMemberProfitLog(clubId int, adminId string) ([]*CLUB_MEMBER_PROFIT_LOG, error) {
	items := []*CLUB_MEMBER_PROFIT_LOG{}
	err := util.WithLogCollection(club_member_profit_log_c, func(c *mgo.Collection) error {
		return c.Find(bson.M{"clubId": clubId, "adminId": adminId}).All(&items)
	})
	return items, err
}

func AddClubMemberProfitLogBigGameResult(clubId int, adminId string, userId string, usdChg float32, inSureProfitChg float32, gameTaxChg float32) error {
	return util.WithLogCollection(club_member_profit_log_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"clubId": clubId, "adminId": adminId, "userId": userId}, bson.M{"$inc": bson.M{"totalUsdChg": usdChg, "inSureProfit": inSureProfitChg, "gameTax": gameTaxChg}})
		return err
	})
}
