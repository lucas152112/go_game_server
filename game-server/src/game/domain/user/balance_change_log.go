package user

import (
	"game/domain/core"
	"game/pb"
	"game/util"
	"time"

	mgo "gopkg.in/mgo.v2"
)

type BalanceChangeInfo struct {
	UserId  string      `bson:"userId"`
	Reason  string      `bson:"reason"`
	Value   pb.Currency `bson:"value"`
	Result  pb.Currency `bson:"result"`
	TimeStr string      `bson:"timeStr"`
	Time    time.Time   `bson:"time"`
}

const (
	user_balance_change_log_c = "user_balance_change_log"
)

func SaveBalanceChangeLog(log *BalanceChangeInfo) error {
	return util.WithLogCollection(user_balance_change_log_c, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}

type UserClubCoinsChangeInfo struct {
	ClubId  int         `bson:"clubId"`
	Reason  string      `bson:"reason"`
	Value   pb.Currency `bson:"value"`
	Result  pb.Currency `bson:"result"`
	TimeStr string      `bson:"timeStr"`
	Time    time.Time   `bson:"time"`
}

func SaveUserClubCoinsChangeLog(log *UserClubCoinsChangeInfo) error {
	return util.WithLogCollection("user_club_coins_change_log", func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}

//SaveBeanChangeLog ...
func SaveBeanChangeLog(log *core.BeanChangeLog) error {
	core.DB.Save(log)
	return nil
}

//SaveCoinChangeLog ...
func SaveCoinChangeLog(log *core.CoinChangeLog) error {
	core.DB.Save(log)
	return nil
}

//SaveBalanceChangeLogNew ...
func SaveBalanceChangeLogNew(log *core.BalanceChangeLog) error {
	core.DB.Save(log)
	return nil
}

//SaveGiftLog ...
func SaveGiftLog(log *core.GiftLog) error {
	core.DB.Save(log)
	return nil
}

//GetCoinChangeLog ...
func GetCoinChangeLog(userID, begin, end string, page, pageCount int) (int, []*core.CoinChangeLog) {
	db := core.DB

	var count int
	var logs []*core.CoinChangeLog

	if userID != "" {
		db.Model(&core.CoinChangeLog{}).Where("user_id = ? and created_at > ? and created_at < ?", userID, begin, end).Count(&count)

		if count > 0 {
			db.Order("id desc").Where("user_id = ? and created_at > ? and created_at < ?", userID, begin, end).Limit(pageCount).Offset(page * pageCount).Find(&logs)
		}
	} else {
		db.Model(&core.CoinChangeLog{}).Where("created_at > ? and created_at < ?", begin, end).Count(&count)

		if count > 0 {
			db.Order("id desc").Where("created_at > ? and created_at < ?", begin, end).Limit(pageCount).Offset(page * pageCount).Find(&logs)
		}
	}

	return count, logs
}

//GetBalanceChangeLog ...
func GetBalanceChangeLog(userID, begin, end string, page, pageCount int) (int, []*core.BalanceChangeLog) {
	db := core.DB

	var count int
	var logs []*core.BalanceChangeLog

	if userID != "" {
		db.Model(&core.BalanceChangeLog{}).Where("user_id = ? and created_at > ? and created_at < ?", userID, begin, end).Count(&count)

		if count > 0 {
			db.Order("id desc").Where("user_id = ? and created_at > ? and created_at < ?", userID, begin, end).Limit(pageCount).Offset(page * pageCount).Find(&logs)
		}
	} else {
		db.Model(&core.BalanceChangeLog{}).Where("created_at > ? and created_at < ?", begin, end).Count(&count)

		if count > 0 {
			db.Order("id desc").Where("created_at > ? and created_at < ?", begin, end).Limit(pageCount).Offset(page * pageCount).Find(&logs)
		}
	}

	return count, logs
}

//GetBeanChangeLog ...
func GetBeanChangeLog(userID, begin, end string, page, pageCount int) (int, []*core.BeanChangeLog) {
	db := core.DB

	var count int
	var logs []*core.BeanChangeLog
	if userID != "" {
		db.Model(&core.BeanChangeLog{}).Where("user_id = ? and created_at > ? and created_at < ?", userID, begin, end).Count(&count)

		if count > 0 {
			db.Order("id desc").Where("user_id = ? and created_at > ? and created_at < ?", userID, begin, end).Limit(pageCount).Offset(page * pageCount).Find(&logs)
		}
	} else {
		db.Model(&core.BeanChangeLog{}).Where("created_at > ? and created_at < ?", begin, end).Count(&count)

		if count > 0 {
			db.Order("id desc").Where("created_at > ? and created_at < ?", begin, end).Limit(pageCount).Offset(page * pageCount).Find(&logs)
		}
	}

	return count, logs
}

//GetGiftLog ...
func GetGiftLog(liverID, begin, end string, page, pageCount int) (int, []*core.GiftLog) {
	db := core.DB

	var count int
	db.Model(&core.GiftLog{}).Where("receiver = ? and created_at > ? and created_at < ?", liverID, begin, end).Count(&count)

	var logs []*core.GiftLog

	if count > 0 {
		db.Order("id desc").Where("receiver = ? and created_at > ? and created_at < ?", liverID, begin, end).Limit(pageCount).Offset(page * pageCount).Find(&logs)
	}

	return count, logs
}
