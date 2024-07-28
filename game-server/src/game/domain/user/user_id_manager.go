package user

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

const (
	userIdC = "user_id_index"
)

type sequence struct {
	Name string `bson:"name"`
	Val  int64  `bson:"seq"`
}

func GetNewUserId() string {
	newId := GetPlayerManager().GetRandUserId()
	return newId
}

func GetUserFortuneTable(userIdStr string) (string, error) {
	userIdInt, err := strconv.Atoi(userIdStr)
	if err != nil {
		glog.Info("GetUserFortuneTable err:", err)
		return "", err
	}

	index := userIdInt/30000 + 1
	tableName := "user_fortune_" + fmt.Sprintf("%v", index)
	return tableName, nil
}

func GetUserMatchRecordTable(userIdStr string) (string, error) {
	userIdInt, err := strconv.Atoi(userIdStr)
	if err != nil {
		glog.Info("GetUserMatchRecordTable err:", err)
		return "", err
	}

	index := userIdInt/100000 + 1
	tableName := "match_record_" + fmt.Sprintf("%v", index)
	return tableName, nil
}

func GetUserDZMatchRecordTable(userIdStr string) string {
	userIdInt, err := strconv.Atoi(userIdStr)
	if err != nil {
		glog.Info("GetUserDZMatchRecordTable err:", err)
		return ""
	}

	index := userIdInt/100000 + 1
	tableName := "dz_match_record_" + fmt.Sprintf("%v", index)
	return tableName
}
