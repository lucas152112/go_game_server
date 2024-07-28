package core

import (
	"fmt"

	"github.com/jinzhu/gorm"

	// Enable PostgreSQL support
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// Enable MySQL support

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	// DB is the global datbase object
	DB *gorm.DB
)

// InitDB initialize database
func InitDB(dbtype string, param string) {
	fmt.Println("Initializing database:", dbtype, " param:", param)
	var errDb error
	//db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	DB, errDb = gorm.Open(dbtype, param)
	if errDb != nil {
		fmt.Println("open db error:", errDb)
	} else {
		DB.AutoMigrate(&UserGift{})
		DB.AutoMigrate(&UserTakeOutGiftLog{})
		DB.AutoMigrate(&UserPayInfoLog{})
		DB.AutoMigrate(&UserPayLog{})
		DB.AutoMigrate(&UserRegisterLog{}, &UserLoginLog{}, &UserActivityLog{}, &UserChannelActiveLog{})
		DB.AutoMigrate(&UserFeebackLog{}, &UserData{}, &UserTaskData{})
		DB.AutoMigrate(&SysMessage{}, &InviteLog{}, &ActionLivingLog{}, &LivingRoomState{}, &AllLivingRoomState{})
		DB.AutoMigrate(&UserFollow{}, &ChannelState{}, &LivingTask{}, &LivingSignReward{}, &ChannelList{})
		DB.AutoMigrate(&CoinChangeLog{}, &BeanChangeLog{}, &BalanceChangeLog{}, &GiftLog{})
	}
}
