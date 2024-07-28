package core

import (
	"github.com/jinzhu/gorm"
)

//UserGift ...
type UserGift struct {
	gorm.Model
	UserID     string //ID
	GiftAmount int64  //礼物数量
}

//UserTakeOutGiftLog ...
type UserTakeOutGiftLog struct {
	gorm.Model
	UserID     string //ID
	Status     int    //状态， 1:申请， 2：完成， 3：拒绝
	Amount     int64  //数量
	PhoneNum   string //手机
	ApplyTime  string //申请时间
	ManageTime string //处理时间
}

//UserPayInfoLog ...
type UserPayInfoLog struct {
	gorm.Model
	UserID   string //ID
	Amount   int    //支付金额
	PayTimes int    //支付次数
	TimeStr  string //
}

//UserPayLog ...
type UserPayLog struct {
	gorm.Model
	UserID        string `gorm:"size:32"`   //ID
	Amount        int    `gorm:"default:0"` //支付金额
	TimeStr       string `gorm:"size:32"`   //
	Channel       string `gorm:"size:32"`   //渠道
	OrderID       string `gorm:"size:128"`  //我们的订单
	ProductID     int    `gorm:"default:0"` //商品ID
	ThreeOrder    string `gorm:"size:128"`  //第三方订单
	PurchaseToken string `gorm:"size:1024"` //
}

//UserRegisterLog ...
type UserRegisterLog struct {
	gorm.Model
	UserID  string //
	TimeStr string //
	Channel string //
	IsTest  bool
}

//UserLoginLog ...
type UserLoginLog struct {
	gorm.Model
	UserID  string //
	TimeStr string //
	Channel string //
}

//UserActivityLog ...
type UserActivityLog struct {
	gorm.Model
	UserID   string //
	Channel  string //
	AddDay   string //
	OneDay   int    //
	ThreeDay int    //
	SevenDay int    //
}

//UserChannelActiveLog ...
type UserChannelActiveLog struct {
	gorm.Model
	Channel   string //
	TimeStr   string //
	UserCount int    //
}

//UserFeebackLog ...
type UserFeebackLog struct {
	gorm.Model
	UserID  string //
	Message string //
	Contact string //
	TimeStr string //
}

//UserData ...
type UserData struct {
	gorm.Model
	UserID           string  //
	LastSign         int64   //
	Sign             int     //
	Reward           string  //
	BuyCheap         int     `gorm:"default:0"` //0.99美金特惠
	MsgID            uint    `gorm:"default:0"` //已经发送的消息ID
	InviteRewardRate float64 `gorm:"default:0"`
	LoginDays        int     //累计登录天数
}

//UserTaskData ...
type UserTaskData struct {
	gorm.Model
	UserID        string //
	ResetTime     int64  //
	Task          int64  `gorm:"default:0"`
	GetStatus     int64  `gorm:"default:0"`
	ShareGetTimes int    `gorm:"default:0"` //分享领取奖励的次数
}

//SysMessage ...
type SysMessage struct {
	gorm.Model
	Title     string `gorm:"size:128"`
	Content   string `gorm:"size:512"`
	BeginTime string `gorm:"size:32"`
	EndTime   string `gorm:"size:32"`
	UserID    string `gorm:"size:32"`
}

//InviteLog ...
type InviteLog struct {
	gorm.Model
	InvitorID  string  `gorm:"size:32"`                  //代理，或者叫邀请人
	UserID     string  `gorm:"size:32"`                  //被邀请的用户ID
	Nickname   string  `gorm:"size:128" json:"nickname"` //被邀请的用户昵称
	Reward     float32 `json:"reward"`                   //奖励数量
	RewardType int     `json:"rewardType"`               //奖励类型 1 金币， 2 钻石
	Status     int     `json:"status"`                   //状态， 1 等待领取， 2 已经发放
	Reason     int     `json:"reason"`                   //奖励原因， 1 邀请新用户注册绑定， 2 充值分红
	TimeStr    string  `json:"time"`                     //活得奖励时间
	PayAmount  int     `gorm:"default:0" json:"payAmount"`
}

//ActionLivingLog ...
type ActionLivingLog struct {
	gorm.Model
	UserID  string `gorm:"size:32"` //用户ID
	LiverID string //主播ID
	Action  int    //用户行为
	Tax     int64  `gorm:"default:0"` //税收
}

//LivingRoomState ...
type LivingRoomState struct {
	gorm.Model
	TimeStr         string //日期
	LiverID         string //主播ID
	TimeLength      int64  //开播时长
	EnterUser       int    //进入人数
	SeatUser        int    //坐下人数
	GameUser        int    //游戏人数
	ChatUser        int    //聊天人数
	GiftUser        int    //赠礼物人数
	DiamondConsum   int64  //钻石消耗
	GameCount       int    //游戏手数
	GoldConsum      int64  //金币消耗
	RemainLength    int64  //总的停留时长
	OneDayRemain    int    //次日留存
	OneDayRate      int    //次日留存率
	GameLength      int64  //玩家游戏时长
	NewEnter        int    `gorm:"default:0"` //新用户进入直播间人数
	NewBeat         int    `gorm:"default:0"` //新用户下注人数
	NewGift         int    `gorm:"default:0"` //新用户送礼物人数
	NewGiftDiamond  int64  `gorm:"default:0"` //新用户送的钻石
	ChatDiamond     int64  `gorm:"default:0"` //聊点钻石
	DelayDiamond    int64  `gorm:"default:0"` //延迟钻石
	LookCardDiamond int64  `gorm:"default:0"` //看牌钻石
	NewFollow       int    `gorm:"default:0"` //新用户关注
	Follow          int    `gorm:"default:0"` //粉丝数
	AddFollow       int    `gorm:"default:0"` //新粉丝
	BeanConsum      int64  `gorm:"default:0"` //金豆消耗
	BeanGameUser    int    `gorm:"default:0"` //金豆游戏人数
}

//AllLivingRoomState ...
type AllLivingRoomState struct {
	gorm.Model
	TimeStr         string //日期
	LiverCount      int    //主播数量
	TimeLength      int64  //开播时长
	EnterUser       int    //进入人数
	SeatUser        int    //坐下人数
	GameUser        int    //游戏人数
	ChatUser        int    //聊天人数
	GiftUser        int    //赠礼物人数
	DiamondConsum   int64  //钻石消耗
	GameCount       int    //游戏手数
	GoldConsum      int64  //金币消耗
	RemainLength    int64  //总的停留时长
	OneDayRemain    int    //次日留存
	OneDayRate      int    //次日留存率
	GameLength      int64  //玩家游戏时长
	Diamond2Gold    int64  //钻石兑换成金币
	NewEnter        int    `gorm:"default:0"` //新用户进入直播间人数
	NewBeat         int    `gorm:"default:0"` //新用户下注人数
	NewGift         int    `gorm:"default:0"` //新用户送礼物人数
	NewGiftDiamond  int64  `gorm:"default:0"` //新用户送的钻石
	ChatDiamond     int64  `gorm:"default:0"` //聊点钻石
	DelayDiamond    int64  `gorm:"default:0"` //延迟钻石
	LookCardDiamond int64  `gorm:"default:0"` //看牌钻石
	NewFollow       int    `gorm:"default:0"` //新用户关注
	Follow          int    `gorm:"default:0"` //粉丝数
	AddFollow       int    `gorm:"default:0"` //新粉丝
	BeanConsum      int64  `gorm:"default:0"` //金豆消耗
	BeanGameUser    int    `gorm:"default:0"` //金豆游戏人数
}

//UserFollow ...
type UserFollow struct {
	gorm.Model
	UserID  string `gorm:"size:32"`
	LiverID string `gorm:"size:32"`
}

//ChannelState ...
type ChannelState struct {
	gorm.Model
	TimeStr      string `gorm:"size:32"` //日期
	Channel      string //渠道
	NewPayUser   int    `gorm:"default:0"` //新充值人数
	NewPayCount  int    `gorm:"default:0"` //新充值笔数
	NewPaySum    int    `gorm:"default:0"` //新充值金额
	NewBeat      int    `gorm:"default:0"` //新游戏
	NewPayRate   int    `gorm:"default:0"` //充值成功率
	PayUser      int    `gorm:"default:0"` //充值人数
	PayCount     int    `gorm:"default:0"` //充值笔数
	PaySum       int    `gorm:"default:0"` //充值金额
	LoginUser    int    `gorm:"default:0"` //登录人数
	PayRate      int    `gorm:"default:0"` //充值成功率
	EnterUser    int    `gorm:"default:0"` //进入人数
	GameUser     int    `gorm:"default:0"` //游戏人数
	DiamondUser  int    `gorm:"default:0"` //钻石消耗人数
	DiamondSum   int    `gorm:"default:0"` //钻石消耗总数
	RemainLength int64  `gorm:"default:0"` //总的停留时长
	OneDay       int    `gorm:"default:0"` //一日留
	ThreeDay     int    `gorm:"default:0"` //三日留
	SevenDay     int    `gorm:"default:0"` //七日留
}

//LivingTask ...
type LivingTask struct {
	gorm.Model
	UserID        string `gorm:"size:32"`   //用户ID
	TaskResetTime int64  `gorm:"default:0"` //重置的时间戳
	Task          int64  `gorm:"default:0"` //111, 110, 100, 11, 10, 1, 0
	TaskStatus    int64  `gorm:"default:0"` //同上
	LastSignTime  int64  `gorm:"default:0"` //最后签到时间
	Sign          int64  `gorm:"default:0"` //1111111
	RemainLength  int    `gorm:"default:0"` //观看时长
	IsGetReward   int    `gorm:"default:0"` //是否领取了签到总奖励
}

//LivingSignReward ...
type LivingSignReward struct {
	gorm.Model
	UserID  string `gorm:"size:32"`   //用户ID
	SignDay string `gorm:"size:64"`   //签到时间
	Phone   string `gorm:"size:32"`   //手机号
	Line    string `gorm:"size:128"`  //line号
	Status  int    `gorm:"default:0"` //状态， 0申请， 1同意， 2拒绝
}

//ChannelList ...
type ChannelList struct {
	Key  string
	Name string
}

//BeanChangeLog ...
type BeanChangeLog struct {
	gorm.Model
	UserID string `gorm:"size:32"`
	Reason string `gorm:"size:32"`
	Old    int64  `gorm:"default:0"`
	Value  int64  `gorm:"default:0"`
	Result int64  `gorm:"default:0"`
}

//CoinChangeLog ...
type CoinChangeLog struct {
	gorm.Model
	UserID string `gorm:"size:32"`
	Reason string `gorm:"size:128"`
	Old    int64  `gorm:"default:0"`
	Value  int64  `gorm:"default:0"`
	Result int64  `gorm:"default:0"`
}

//BalanceChangeLog ...
type BalanceChangeLog struct {
	gorm.Model
	UserID string `gorm:"size:32"`
	Reason string `gorm:"size:128"`
	Old    int64  `gorm:"default:0"`
	Value  int64  `gorm:"default:0"`
	Result int64  `gorm:"default:0"`
}

//GiftLog ...
type GiftLog struct {
	gorm.Model
	Sender     string `gorm:"size:32"`   //赠送者
	RecordType string `gorm:"size:128"`  //类型
	Receiver   string `gorm:"size:32"`   //接收者
	Currency   int    `gorm:"default:0"` //花费类型 8钻石, 6金豆
	Name       string `gorm:"size:64"`   //礼物名称
	Price      int    `gorm:"default:0"` //礼物价格
	Num        int    `gorm:"default:0"` //礼物数量
	Left       int64  `gorm:"default:0"` //余额
	Nickname   string `gorm:"size:64"`   //昵称
	TimeStr    string `gorm:"size:32"`   //
}
