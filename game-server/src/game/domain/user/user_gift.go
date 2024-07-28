package user

import (
	"errors"
	"game/domain/core"
	logServer "game/domain/logServer"
	"game/util"

	"github.com/jinzhu/gorm"
)

func GetUserGift(userID string) int64 {
	var userGift core.UserGift
	db := core.DB
	if db.Where("user_id = ?", userID).First(&userGift).RecordNotFound() {
		userGift.UserID = userID
		userGift.GiftAmount = 0
	}

	return userGift.GiftAmount
}

func UpdateUserGift(userID string, amount int64) (error, int64) {
	var userGift core.UserGift
	db := core.DB
	if db.Where("user_id = ?", userID).First(&userGift).RecordNotFound() {
		userGift.UserID = userID
		userGift.GiftAmount = 0
	}
	oldAmount := userGift.GiftAmount

	temp := userGift.GiftAmount + amount
	if temp < 0 {
		return errors.New("not enough amount"), userGift.GiftAmount
	}

	userGift.GiftAmount += amount
	go db.Save(&userGift)

	var takeOutLog core.UserTakeOutGiftLog
	takeOutLog.Amount = oldAmount
	takeOutLog.UserID = userID
	takeOutLog.PhoneNum = ""
	takeOutLog.Status = 1
	takeOutLog.ApplyTime = util.GetDbTime()

	go db.Save(&takeOutLog)

	go logServer.ChangeTotalDiamond(-int64(amount)) //
	go logServer.DiamdondDel(int64(amount))         //

	return nil, temp
}

func SaveUserTakeOutGiftLog(userID string, phoneNum string) error {
	var userGift core.UserGift
	db := core.DB
	if db.Where("user_id = ?", userID).First(&userGift).RecordNotFound() {
		userGift.UserID = userID
		userGift.GiftAmount = 0
		return errors.New("GIFT AMOUNT IS 0")
	}

	if userGift.GiftAmount <= 0 {
		return errors.New("GIFT AMOUNT IS 0")
	}

	oldAmout := userGift.GiftAmount
	userGift.GiftAmount = 0

	db.Save(&userGift)

	var takeOutLog core.UserTakeOutGiftLog
	takeOutLog.Amount = oldAmout
	takeOutLog.UserID = userID
	takeOutLog.PhoneNum = phoneNum
	takeOutLog.Status = 1
	takeOutLog.ApplyTime = util.GetDbTime()

	db.Save(&takeOutLog)

	go logServer.ChangeTotalDiamond(-int64(oldAmout)) //
	go logServer.DiamdondDel(int64(oldAmout))         //
	return nil
}

func WebGetUserTakeOutGifgLog(beginTime string, endTime string, page int, pageSize int) (int, []*core.UserTakeOutGiftLog) {
	var logs []*core.UserTakeOutGiftLog

	var log core.UserTakeOutGiftLog

	db := core.DB

	//totalPage := 0
	t := 0
	db.Model(&log).Where("apply_time >= ? and apply_time <= ?", beginTime, endTime).Count(&t)

	//totalPage = t / pageSize
	//left := t % pageSize
	//if left != 0 {
	//	totalPage += 1
	//}

	db.Where("apply_time >= ? and apply_time <= ?", beginTime, endTime).Scopes(Paginate(page, pageSize)).Find(&logs)
	return t, logs
}

func GetUserTakeOutGifgLog(userID string, page int, pageSize int) (int, []*core.UserTakeOutGiftLog) {
	var logs []*core.UserTakeOutGiftLog

	var log core.UserTakeOutGiftLog

	db := core.DB

	totalCount := 0
	t := 0
	db.Model(&log).Where("user_id = ?", userID).Count(&t)

	totalCount = t / pageSize
	left := t % pageSize
	if left != 0 {
		totalCount += 1
	}

	db.Where("user_id = ?", userID).Scopes(Paginate(page, pageSize)).Find(&logs)
	return totalCount, logs
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 0 {
			page = 0
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := page * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func GetUserGiftAmount(userID string) int64 {
	db := core.DB

	var userGift core.UserGift
	if db.Where("user_id = ?", userID).First(&userGift).RecordNotFound() {
		userGift.UserID = userID
		userGift.GiftAmount = 0

		go db.Save(&userGift)
	}

	return userGift.GiftAmount
}

func ManageTakeOutApply(id uint, status int) (error, string) {
	db := core.DB
	var log core.UserTakeOutGiftLog
	if db.Where("id = ?", id).First(&log).RecordNotFound() {
		return errors.New("not this log"), ""
	}

	if log.Status != 1 {
		return errors.New("duplicate manage this apply"), ""
	}

	log.ManageTime = util.GetDbTime()
	log.Status = status

	db.Save(&log)

	if status == 3 {
		go UpdateUserGift(log.UserID, log.Amount)
	}

	return nil, log.ManageTime
}
