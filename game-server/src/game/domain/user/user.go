package user

import (
	"errors"
	"fmt"
	"game/domain/core"
	"game/pb"
	"game/util"
	"time"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	gorm.Model
	UserId            string      `bson:"userId"`
	UserName          string      `bson:"userName"`
	Password          string      `bson:"password"`
	Nickname          string      `bson:"nickname"`
	Gender            int         `bson:"gender"`
	Signiture         string      `bson:"signiture"`
	PhotoUrl          string      `bson:"photoUrl"`
	IsBind            bool        `bson:"isBind"`
	CreateTime        string      `bson:"createTime"`
	ChannelId         string      `bson:"channelId"`
	ChannelInt        int         `bson:"channelInt"`
	IsChangedNickname bool        `bson:"isChangedNickname"`
	DeviceModel       string      `bson:"deviceModel"` //定义model
	IsLocked          bool        `bson:"isLocked"`
	VersionName       string      `bson:"versionName"`
	Imei              string      `bson:"imei"`
	JingDu            float64     `gorm:"-"`
	WeiDu             float64     `gorm:"-"`
	IP                string      `gorm:"-"`
	Balance           pb.Currency `bson:"balance"` //钻石
	Coins             pb.Currency `bson:"coins"`   //金币
	InviteUserId      string      `bson:"inviteUserId"`
	WinGames          int         `bson:"winGames"`    //赢的局数
	TotalGames        int         `bson:"totalGames"`  //总局数
	IdType            int         `bson:"idType"`      //0玩家 20机器人
	FirstClubId       int         `bson:"firstClubId"` //第一次加入的俱乐部id
	AdminRole         int         `bson:"adminRole"`   //管理角色，0普通，100运营帐户,200 僵尸
	IsGuest           bool        `bson:"isGuest"`
	BGuest            bool        `gorm:"-"` //临时变量
	LastLogin         string      `bson:"lastLogin"`
	Language          int         `gorm:"default:0"` //语言
	FanCount          int         `gorm:"default:0"` //粉丝数量
	FollowCount       int         `gorm:"default:0"` //关注数量
	PushUserName      string      `gorm:"size:64"`
	TotalWinCoins     int64       `gorm:"default:0"`
	PlayCount         int         `gorm:"default:0"` //总手数
	InPool            int         `gorm:"default:0"` //主动入池
	PreAdd            int         `gorm:"default:0"` //翻前加注
	IsDel             int         `gorm:"default:0"` //
	Bean              pb.Currency `gorm:"default:0"` //金豆
}

//积分
func (u *User) ChangeScore(score pb.Currency, reason string) bool {
	if score < 0 {
		temp := u.Balance + score
		if temp < 0 {
			return false
		}
	}
	old := u.Balance
	u.Balance += score
	temp := u.Balance
	changeScore(u.ID, temp)

	logNew := &core.BalanceChangeLog{}
	logNew.UserID = u.UserId
	logNew.Reason = reason
	logNew.Old = int64(old)
	logNew.Value = int64(score)
	logNew.Result = int64(u.Balance)

	go SaveBalanceChangeLogNew(logNew)

	return true
}

func (u *User) ResetScore(score pb.Currency, reason string) bool {
	old := u.Balance
	u.Balance = score

	temp := u.Balance
	changeScore(u.ID, temp)

	logNew := &core.BalanceChangeLog{}
	logNew.UserID = u.UserId
	logNew.Reason = reason
	logNew.Old = int64(old)
	logNew.Value = int64(score)
	logNew.Result = int64(u.Balance)

	go SaveBalanceChangeLogNew(logNew)

	return true
}

func (u *User) ChangeCoinsForcely(score pb.Currency, reason string) bool {
	if u.IdType == int(util.IdType_Robot) {
		return true
	}
	old := u.Coins
	u.Coins += score
	temp := u.Coins
	changeCoins(u.ID, temp)

	logNew := &core.CoinChangeLog{}
	logNew.UserID = u.UserId
	logNew.Reason = reason
	logNew.Old = int64(old)
	logNew.Value = int64(score)
	logNew.Result = int64(u.Coins)

	go SaveCoinChangeLog(logNew)

	if u.IdType != int(util.IdType_Robot) {
		go AddRank(RankingBalance, u.UserId, int(score))
	}

	return true
}

func (u *User) ChangeCoins(score pb.Currency, reason string) bool {
	if u.IdType == int(util.IdType_Robot) {
		return true
	}
	if score < 0 {
		temp := u.Coins + score
		if temp < 0 {
			return false
		}
	}
	old := u.Coins
	u.Coins += score

	temp := u.Coins
	changeCoins(u.ID, temp)

	logNew := &core.CoinChangeLog{}
	logNew.UserID = u.UserId
	logNew.Reason = reason
	logNew.Old = int64(old)
	logNew.Value = int64(score)
	logNew.Result = int64(u.Coins)

	go SaveCoinChangeLog(logNew)
	if u.IdType != int(util.IdType_Robot) {
		go AddRank(RankingBalance, u.UserId, int(score))
	}

	return true
}

func (u *User) ResetCoins(score pb.Currency, reason string) bool {
	old := u.Coins
	u.Coins = score
	temp := u.Coins
	changeCoins(u.ID, temp)

	logNew := &core.CoinChangeLog{}
	logNew.UserID = u.UserId
	logNew.Reason = reason
	logNew.Old = int64(old)
	logNew.Value = int64(score)
	logNew.Result = int64(u.Coins)

	go SaveCoinChangeLog(logNew)
	if u.IdType != int(util.IdType_Robot) {
		go AddRank(RankingBalance, u.UserId, int(score))
	}

	return true
}

//ChangeBeans ...
func (u *User) ChangeBeans(beans pb.Currency, reason string) bool {
	if u.IdType == int(util.IdType_Robot) {
		return true
	}
	if beans < 0 {
		temp := u.Bean + beans
		if temp < 0 {
			return false
		}
	}
	old := u.Bean
	u.Bean += beans

	temp := u.Bean
	changeBeans(u.ID, temp)

	log := &core.BeanChangeLog{}
	log.UserID = u.UserId
	log.Reason = reason
	log.Old = int64(old)
	log.Value = int64(beans)
	log.Result = int64(u.Bean)

	go SaveBeanChangeLog(log)

	return true
}

func (u *User) ChangeFirstClubId(clubId int) bool {
	u.FirstClubId = clubId

	changeFirstClub(u.ID, clubId)

	return true
}

func (u *User) BuildMessage() *pb.DZUserInfoAck {
	ack := &pb.DZUserInfoAck{}
	ack.UserId = u.UserId
	ack.UserName = u.UserName
	ack.Nickname = u.Nickname
	ack.PhotoUrl = u.PhotoUrl
	ack.IsBind = u.IsBind
	ack.Gender = u.Gender
	ack.Signiture = u.Signiture
	ack.Score = u.Balance
	ack.SafeToken = ""
	ack.InviteUserId = u.InviteUserId
	ack.Coins = u.Coins
	ack.Balance = u.Balance
	ack.AdminRole = u.AdminRole
	ack.FanCount = u.FanCount
	ack.FollowCount = u.FollowCount
	ack.Beans = u.Bean
	return ack
}

const (
	userNameIdC     = "user_name_id"
	userC           = "user"
	phoneUserC      = "phone_user"
	userNicknameIdC = "user_nickname_id"
	admin_user_id_c = "admin_user_id"
)

type UserNameId struct {
	UserName string `bson:"userName"`
	UserId   string `bson:"userId"`
}

// NickName
type UserNicknameId struct {
	Nickname  string `bson:"nickname"`
	UserId    string `bson:"userId"`
	ModifyNum int    `json:"modifyNum"`
}

func SaveUserNicknameId(userId string, nickname string, modifyNum int) error {
	u := &UserNicknameId{}
	u.UserId = userId
	u.Nickname = nickname
	u.ModifyNum = modifyNum
	return util.WithSafeUserCollection(userNicknameIdC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"nickname": u.Nickname}, u)
		return err
	})
}

//GetUserNicknameByID ...
func GetUserNicknameByID(userID string) (*UserNicknameId, error) {
	item := UserNicknameId{}
	err := util.WithUserCollection(userNicknameIdC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userID}).One(&item)
	})
	return &item, err
}

//FindByUserId ...
func FindByUserId(userID string) (*User, error) {
	db := DB
	var user User

	if db.Where("user_id = ?", userID).First(&user).RecordNotFound() {
		return &user, errors.New("not find")
	}

	return &user, nil
}

//FindByUserName ...
func FindByUserName(username string) (*User, error) {
	db := DB
	var user User

	if db.Where("user_name = ?", username).First(&user).RecordNotFound() {
		return &user, errors.New("not find")
	}

	return &user, nil
}

//FindByNickname ...
func FindByNickname(nickname string) []*User {
	db := DB
	var users []*User

	db.Where("nickname = ?", nickname).Find(&users)

	return users
}

//FindOneByNickname ...
func FindOneByNickname(nickname string) *User {
	db := DB
	var user User

	if db.Where("nickname = ?", nickname).First(&user).RecordNotFound() {
		return nil
	}

	return &user
}

//SaveUser ...保存用户
func SaveUser(u *User) error {
	return saveUser(u)
}

type PhoneUser struct {
	Phone    string    `bson:"phone"`
	UserId   string    `bson:"userId"`
	BindTime time.Time `bson:"bindTime"`
}

func GetPhoneIsBind(phone string) (string, error) {
	userTemp := PhoneUser{"", "", time.Now()}
	err := util.WithUserCollection(phoneUserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"phone": phone}).One(&userTemp)
	})
	return userTemp.UserId, err
}

func SavePhoneUser(phone string, userId string) error {
	u := PhoneUser{phone, userId, time.Now()}
	err := util.WithUserCollection(phoneUserC, func(c *mgo.Collection) error {
		return c.Insert(&u)
	})
	return err
}

func GetBindPhone(userID string) (string, error) {
	userTemp := PhoneUser{"", "", time.Now()}
	err := util.WithUserCollection(phoneUserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userID}).One(&userTemp)
	})
	return userTemp.Phone, err
}

const (
	user_rand_id_c = "user_rand_id"
)

type UserRandId_T struct {
	Id string `bson:"id"`
}

func SaveUserRandId(itemDB *UserRandId_T) error {
	return util.WithUserCollection(user_rand_id_c, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"id": itemDB.Id}, itemDB)
		return err
	})
}

func DeleteUserRandId(userId string) error {
	return util.WithUserCollection(user_rand_id_c, func(c *mgo.Collection) error {
		err := c.Remove(bson.M{"id": userId})
		return err
	})
}

func GetUserRandId() (error, string) {
	item := &UserRandId_T{}
	err := util.WithUserCollection(user_rand_id_c, func(c *mgo.Collection) error {
		return c.Find(nil).One(item)
	})

	if err == nil {
		userId := item.Id
		DeleteUserRandId(userId)
		return nil, userId
	} else {
		return err, ""
	}
}

func ChangeScore(userId string, scoreChg pb.Currency, reason string) bool {
	p := GetPlayerManager().FindPlayerById(userId)
	if p != nil {
		return p.ChangeScore(scoreChg, reason)
	}
	user, err := FindByUserId(userId)
	if err == nil {
		return user.ChangeScore(scoreChg, reason)
	}
	return false
}

func ResetScore(userId string, scoreChg pb.Currency, reason string) bool {
	p := GetPlayerManager().FindPlayerById(userId)
	if p != nil {
		return p.ResetScore(scoreChg, reason)
	}
	user, err := FindByUserId(userId)
	if err == nil {
		return user.ResetScore(scoreChg, reason)
	}
	return false
}

func ChangeCoinsForcely(userId string, scoreChg pb.Currency, reason string) bool {
	p := GetPlayerManager().FindPlayerById(userId)
	if p != nil {
		return p.ChangeCoinsForcely(scoreChg, reason)
	}
	user, err := FindByUserId(userId)
	if err == nil {
		return user.ChangeCoinsForcely(scoreChg, reason)
	}
	return false
}

func ChangeCoins(userId string, scoreChg pb.Currency, reason string) bool {
	p := GetPlayerManager().FindPlayerById(userId)
	if p != nil {
		glog.Info("Player add ")
		return p.ChangeCoins(scoreChg, reason)
	}
	user, err := FindByUserId(userId)
	if err == nil {
		glog.Info("User add ")
		return user.ChangeCoins(scoreChg, reason)
	}
	return false
}

func ResetCoins(userId string, scoreChg pb.Currency, reason string) bool {
	p := GetPlayerManager().FindPlayerById(userId)
	if p != nil {
		return p.ResetCoins(scoreChg, reason)
	}
	user, err := FindByUserId(userId)
	if err == nil {
		return user.ResetCoins(scoreChg, reason)
	}
	return false
}

//ChangeBeans ...
func ChangeBeans(userID string, beans pb.Currency, reason string) bool {
	p := GetPlayerManager().FindPlayerById(userID)
	if p != nil {
		return p.ChangeBeans(beans, reason)
	}

	user, err := FindByUserId(userID)
	if err == nil {
		return user.ChangeBeans(beans, reason)
	}
	return false
}

type AdminUserIDST struct {
	UserId  string    `bson:"userId"`
	AddTime time.Time `bson:"addTime"`
}

func GetAdminUser() *User {
	item := &AdminUserIDST{}
	err := util.WithUserCollection(admin_user_id_c, func(c *mgo.Collection) error {
		return c.Find(nil).One(item)
	})

	if err == mgo.ErrNotFound {
		user := &User{}
		user.AdminRole = 100
		user.UserId = GetNewUserId()
		user.UserName = fmt.Sprintf("systemadmin%v", user.UserId)
		user.Password = user.UserName
		user.CreateTime = util.GetDbTime()
		user.ChannelId = "100"
		user.DeviceModel = "admin100"
		user.VersionName = "100"
		user.ChannelInt = 100
		user.Imei = "systemadmin"
		user.Balance = 0
		user.Coins = 0

		SaveUser(user)

		item.UserId = user.UserId
		item.AddTime = time.Now()

		util.WithUserCollection(admin_user_id_c, func(c *mgo.Collection) error {
			_, err := c.Upsert(bson.M{}, item)
			return err
		})

		return user
	}

	user, errUser := FindByUserId(item.UserId)
	if errUser == nil {
		return user
	}
	return nil
}

//saveUser ...
func saveUser(u *User) error {
	if u.BGuest {
		userID := u.UserId
		item := &UserRandId_T{}
		item.Id = userID
		go SaveUserRandId(item)
		return nil
	}

	DB.Save(u)
	return nil
}

func changeScore(ID uint, score pb.Currency) {
	DB.Table("users").Where("id = ?", ID).Update("balance", score)
}

func changeCoins(ID uint, score pb.Currency) {
	DB.Table("users").Where("id = ?", ID).Update("coins", score)
}

func changeBeans(id uint, score pb.Currency) {
	DB.Table("users").Where("id = ?", id).Update("bean", score)
}

func changeGameTimes(ID uint, winTimes int, totalTimes int, coins int64) {
	DB.Table("users").Where("id = ?", ID).Update(map[string]interface{}{"win_games": winTimes, "total_games": totalTimes, "total_win_coins": coins})
}

func changeFirstClub(ID uint, clubID int) {
	DB.Table("users").Where("id = ?", ID).Update("first_club_id", clubID)
}

func ChangePassword(ID uint, password string) {
	DB.Table("users").Where("id = ?", ID).Update("password", password)
}

func ChangeNickname(ID uint, nickname string) {
	DB.Table("users").Where("id = ?", ID).Update("nickname", nickname)
}

func ChangePhotoURL(ID uint, photoURL string) {
	DB.Table("users").Where("id = ?", ID).Update("photo_url", photoURL)
}

func changeIsGuest(ID uint, isGuest bool) {
	DB.Table("users").Where("id = ?", ID).Update("is_guest", isGuest)
}

func ChangeUserLocked(userID string, isLock bool) {
	DB.Table("users").Where("user_id = ?", userID).Update("is_locked", isLock)
}

func changeUserRole(userID string, role int) {
	DB.Table("users").Where("user_id = ?", userID).Update("admin_role", role)
}

func SetRole(userId string, role int) error {
	changeUserRole(userId, role)
	return nil
}

//UpdateLastLogin ...
func UpdateLastLogin(ID uint) {
	DB.Table("users").Where("id = ?", ID).Update("last_login", util.GetDbTime())
}

//UpdateInviteUser ...
func UpdateInviteUser(ID uint, userID string) {
	DB.Table("users").Where("id = ?", ID).Update("invite_user_id", userID)
}

//UpdateUserLanguage ...
func UpdateUserLanguage(ID uint, language int) {
	DB.Table("users").Where("id = ?", ID).Update("language", language)
}

//GetUserInvitor ...
func GetUserInvitor(userID string) string {
	u, err := FindByUserId(userID)
	if err == nil {
		return u.InviteUserId
	}

	return ""
}

//UpdateUserFanCount 更新
func UpdateUserFanCount(userID string, count int) {
	player := GetPlayerManager().FindPlayerById(userID)
	if player != nil {
		player.User.FanCount += count
	}

	DB.Table("users").Where("user_id = ?", userID).Update("fan_count", gorm.Expr("fan_count+ ?", count))
}

//UpdateUserFollowCount ...
func UpdateUserFollowCount(userID string, count int) {
	DB.Table("users").Where("user_id = ?", userID).Update("follow_count", gorm.Expr("follow_count+ ?", count))
}

//UpdatePushUserName ...
func UpdatePushUserName(ID uint, userName string) {
	DB.Table("users").Where("id = ?", ID).Update("push_user_name", userName)
}

//UpdateUserPlayData ...
func UpdateUserPlayData(userID string, inPool bool, preAdd bool) {
	p := GetPlayerManager().FindPlayerById(userID)
	if inPool {
		if p != nil && p.User != nil {
			p.User.InPool++
		}
		go DB.Table("users").Where("user_id = ?", userID).Update("in_pool", gorm.Expr("in_pool+ ?", 1))
	}

	if preAdd {
		if p != nil && p.User != nil {
			p.User.PreAdd++
		}
		go DB.Table("users").Where("user_id = ?", userID).Update("pre_add", gorm.Expr("pre_add+ ?", 1))
	}
}

//UpdateUserPlayCount ...
func UpdateUserPlayCount(userID string) {
	go DB.Table("users").Where("user_id = ?", userID).Update("play_count", gorm.Expr("play_count+ ?", 1))
}

//DeleteAccount ...
func DeleteAccount(id uint) {
	DB.Table("users").Where("id = ?", id).Update("is_del", int(time.Now().Unix()))
}

//CancelDeleteAccount ...
func CancelDeleteAccount(id uint) {
	DB.Table("users").Where("id = ?", id).Update("is_del", 0)
}

//UpdateUserInfo ...
func UpdateUserInfo(id uint, nickname string, coins int64, balance int64, userName string) {
	glog.Info("UpdateUserInfo id:", id)
	DB.Table("users").Where("id = ?", id).Update(map[string]interface{}{"nickname": nickname, "coins": coins, "balance": balance, "user_name": userName})
}

//DeleteUserNickNameID ...
func DeleteUserNickNameID(userID string) error {
	return util.WithUserCollection(userNicknameIdC, func(c *mgo.Collection) error {
		err := c.Remove(bson.M{"userId": userID})
		return err
	})
}

//DeletePhoneUser ...
func DeletePhoneUser(userID string) error {
	return util.WithUserCollection(phoneUserC, func(c *mgo.Collection) error {
		err := c.Remove(bson.M{"userId": userID})
		return err
	})
}
