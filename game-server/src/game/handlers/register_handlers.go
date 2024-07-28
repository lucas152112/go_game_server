package handlers

import (
	"game/handlers/admin"
	"game/handlers/cm"
	"game/handlers/collection"
	"game/handlers/config"
	"game/handlers/debug"
	"game/handlers/dzclub"
	"game/handlers/dzgame"
	"game/handlers/dzuser"
	"game/handlers/gift"
	"game/handlers/hall"

	"game/handlers/pro/home"
	"game/handlers/task"

	"game/handlers/live"
	"game/handlers/logServer"
	"game/handlers/rankinglist"
	"game/handlers/shop"
	"game/handlers/user"
	"game/pb"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

func registerHandlers(r *mux.Router) {
	glog.V(2).Info("register handlers")
	registerHttpHandlers(r)
	registerUserHandlers()
}

func registerHttpHandlers(r *mux.Router) {
	r.HandleFunc("/online", admin.GetOnlineCountHandler)      //统计 在线数
	r.HandleFunc("/onlineByType", admin.GetOnlineTypeHandler) //统计 游戏类型在线人数
	r.HandleFunc("/add_diamond", admin.NAddDiamondHandler)    //
	r.HandleFunc("/add_horns", admin.NAddHornsHandler)        //

	r.HandleFunc("/setCardConfig", admin.SetCardConfigDataHandler) //配置

	r.HandleFunc("/sendSysMsg", admin.WebBroadcastHandler) //通知广播？系统消息

	r.HandleFunc("/goldLimitUserCount", admin.GetUserCountGoldLimitHandler) //查询范围金币的用户
	r.HandleFunc("/getAllGold", admin.GetAllGoldHandler)                    //查询金币所有用户
	r.HandleFunc("/setCurVersion", admin.SetCurVersionHandler)              //设置配置版本？

	r.HandleFunc("/log/getSysTipLog", admin.GetSystemTipLogHandler)
	r.HandleFunc("/log/getSlotPoolLog", admin.GetSlotPoolLogHandler)
	r.HandleFunc("/log/getGameFeeLog", admin.GetGameFeeLogHandler)
	r.HandleFunc("/log/getCharmLog", admin.GetCharmLogHandler)
	r.HandleFunc("/forbit_user_ip", admin.AddUserForbidIPHandler)        //禁用用户IP
	r.HandleFunc("/white_ip", admin.AddWhiteIPHandler)                   //白名单
	r.HandleFunc("/change_ip_limit", admin.ChangeIPLimitHandler)         //IP限制开发
	r.HandleFunc("/change_advert", admin.ChangeAdvertHandler)            // advert 配置 Bannar
	r.HandleFunc("/admin/stop", admin.DZStopServerHandler)               //游戏退出？
	r.HandleFunc("/log/getRoomsOnline", logServer.GetRoomsOnlineHandler) //获取在线房间

	//根据需求新出的接口
	r.HandleFunc("/v2/getUserInfo", logServer.GetUserInfoHandler)                   //获取用户基本信息
	r.HandleFunc("/v2/changeCoins", logServer.AddCoinsHandler)                      //修改用户账户金币
	r.HandleFunc("/v2/changeBalance", logServer.AddBalanceHandler)                  //修改用户账户钻石
	r.HandleFunc("/v2/changeBeans", logServer.AddBeanHandler)                       //修改用户账户金豆
	r.HandleFunc("/v2/getClubInfo", logServer.GetClubInfoHandler)                   //获取俱乐部的信息
	r.HandleFunc("/v2/getClubUserScore", logServer.GetClubUserScoresHandler)        //获取用户俱乐部的总积分
	r.HandleFunc("/v2/getClubUserScoreList", logServer.GetClubUserListScoreHandler) //获取用户俱乐部积分列表
	r.HandleFunc("/v2/getClubTotalScore", logServer.GetClubTotalScoresHandler)
	r.HandleFunc("/lockUser", admin.LockUserHandler) // 锁定用户 剔除

	//指标统计
	r.HandleFunc("/v2/SummaryPlayer", logServer.GetPlayersSummary) //统计上面用户的信息

	r.HandleFunc("/v/version", config.GetVersionHandler)  //服务的版本号
	r.HandleFunc("/v/notify", config.NotifyVersionChange) //通知版本更新

	r.HandleFunc("/cm", cm.HttpHandler) // ClubManage 后台的接口
	// 服务的版本号

	//我加几个统计 总人数
	//俱乐部
	r.HandleFunc("/log/getUserClubs", logServer.GetUserClubsHandler)       //获取用户俱乐部信息
	r.HandleFunc("/log/getClubUserInfo", logServer.GetClubUserInfoHandler) //获取俱乐部玩家
	r.HandleFunc("/log/getClubProfit", logServer.GetClubProfitHandler)     //俱乐部利润

	r.HandleFunc("/log/getClubAdminsProfit", logServer.GetClubAdminsProfitHandler)
	r.HandleFunc("/log/getClubAdminProfit", logServer.GetClubAdminProfitHandler)
	r.HandleFunc("/log/getClubSummaryLog", logServer.GetClubSummaryLogHandler) //俱乐部统计信息
	r.HandleFunc("/log/getClubInfo", logServer.GetClubInfoHandler)

	r.HandleFunc("/log/getClubRank", logServer.GetClubRankHandler) //排行

	r.HandleFunc("/log/checkAdmin", logServer.CheckClubAdminHandler) //检查用户的俱乐部管理
	r.HandleFunc("/log/getUserRank", logServer.GetUserRankHandler)   //用户排行

	r.HandleFunc("/log/getClubMemberProfit", logServer.GetClubMemberHandler)
	r.HandleFunc("/log/getClubLog", logServer.GetClubLogHandler)
	r.HandleFunc("/log/history", admin.WebHistoryHandler)

	r.HandleFunc("/club/userCount", admin.ChangeClubUserCountHandler) //更新俱乐部用户数量
	r.HandleFunc("/club/count", admin.ChangeClubCountHandler)         //
	r.HandleFunc("/club/amount", admin.ChangeClubDiamond)             //修改俱乐部钻石
	r.HandleFunc("/dz/alliance_clubs", admin.GetAllianceClubsHandler)
	r.HandleFunc("/client/page", user.ClientPageHandler)
	r.HandleFunc("/dz/club_tax_type", admin.ChangeClubTaxingTypeHandler)

	r.HandleFunc("/dz/get_club_tax_type", admin.GetClubTaxingTypeHandler)
	r.HandleFunc("/dz/ios_check", admin.IosCheckHandler)
	r.HandleFunc("/defaultclub/set", admin.SetDefaultClubHandler)       //设置默认俱乐部
	r.HandleFunc("/defaultclub/delete", admin.DeleteDefaultClubHandler) //删除默认俱乐部
	r.HandleFunc("/defaultclub/get", admin.GetRecommendClubsHandler)    //获取默认俱乐部列表

	r.HandleFunc("/game/infos", live.DumpGameInfos) //获取gameInfo
	r.HandleFunc("/backyard", live.Backyard)
	r.HandleFunc("/clubyard", live.Clubyard)

	r.HandleFunc("/change_activity", config.Change_activity)       //修改活动
	r.HandleFunc("/debug/getVerifyCode", debug.GetPhotoVerifyCode) //获取验证码
	r.HandleFunc("/build_create_club", user.BuildUserClub)         //创建账户 创建俱乐部
	r.HandleFunc("/build_join_club", user.BuildSetJoinClub)        //创建账户 加入俱乐部
	r.HandleFunc("/replay", live.ReplayFetch)                      //回放
	r.HandleFunc("/reviews", live.Reviews)                         //本局牌谱
	r.HandleFunc("/winnerrank", live.WinnerRank)                   //本局牌谱

	r.HandleFunc("/club/ChangeClubId", dzclub.ChangeClubId) //修改俱乐部ID
	//主播相关
	r.HandleFunc("/internal/liver/add_liver", admin.AdminAddLiver)              //添加主播
	r.HandleFunc("/internal/liver/set_liver_status", admin.AdminSetLiverStatus) //更改主播状态
	r.HandleFunc("/internal/liver/list_liver", admin.AdminGetLiverList)         //主播列表
	r.HandleFunc("/internal/liver/add_game_config", admin.AdminAddGameConfig)   //添加游戏配置
	r.HandleFunc("/internal/liver/gift_list", admin.GetLiverGiftListHandler)    //是主播获得礼物纪录
	r.HandleFunc("/internal/summary", admin.DataSummaryHandler)                 //每日统计
	r.HandleFunc("/internal/user_remain_rate", admin.DataRemainRateHandler)     //用户留存
	r.HandleFunc("/internal/liver_summary", admin.LiverSummaryHandler)          //主播每日统计

	//比赛相关
	r.HandleFunc("/internal/mtt/create", admin.AdminCreateMttHandler)    //创建比赛
	r.HandleFunc("/internal/mtt/dismiss", admin.AdminDismissMttHandler)  //解散比赛
	r.HandleFunc("/mtt/day_sum", admin.GetMTTDaySumHandler)              //获取比赛每天的统计数据
	r.HandleFunc("/mtt/match", admin.GetMTTMatchDataHandler)             //获取比赛数据
	r.HandleFunc("/mtt/user_fee", admin.GetMTTMatchUserFeeHandler)       //获取俱乐部比赛收入
	r.HandleFunc("/mtt/sys_sum", admin.GetMTTSysDaySumHandler)           //获取系统比赛的收入情况
	r.HandleFunc("/mtt/match_rank_list", admin.GetMTTRankingListHandler) //获取系统比赛用户排名列表
	r.HandleFunc("/mtt/sys_match_list", admin.GetMTTMatchSysFeeHandler)  //获取系统比赛结算列表

	r.HandleFunc("/gift/take_out_list", gift.WebGetApplyTakeOutListHandler)                //获取提现申请记录
	r.HandleFunc("/gift/manage_take_out_apply", gift.WebManageTakeOutHandler)              //获取提现申请记录
	r.HandleFunc("/robot/get_pool_coins", admin.GetRobotPoolCoinsHandler)                  //机器人金币池
	r.HandleFunc("/robot/get_pool_beans", admin.GetRobotPoolBeansHandler)                  //机器人金豆池
	r.HandleFunc("/robot/set_pool_coins", admin.SetRobotCoinsPoolHandler)                  //设置机器人金币奖励池
	r.HandleFunc("/robot/set_pool_beans", admin.SetRobotBeansPoolHandler)                  //设置机器人金豆奖励池
	r.HandleFunc("/pay/get_pay_log", admin.GetPayLogHandler)                               //获取支付记录
	r.HandleFunc("/hall/feedback", admin.GetFeedbackLogHandler)                            //获取反馈记录
	r.HandleFunc("/user/invite_reward_rate", admin.SetInviteRewardRate)                    //设置分红比例
	r.HandleFunc("/user/get_invite_reward_log", admin.GetFenHongLogHandler)                //获取分红记录
	r.HandleFunc("/sys/add_message", admin.AddSysMessageHandler)                           //添加系统消息
	r.HandleFunc("/sys/get_messages", admin.GetSysMessagesHandler)                         //获取系统消息
	r.HandleFunc("/sys/del_message", admin.DelSysMessageHandler)                           //删除系统消息
	r.HandleFunc("/living/liver_log", admin.GetLivingLiverLogHandler)                      //直播间主播统计
	r.HandleFunc("/living/all_log", admin.GetAllLivingLogHandler)                          //直播间汇总统计
	r.HandleFunc("/user/check_icon", admin.CheckHeadIconHandler)                           //审核用户头像
	r.HandleFunc("/channel/channel_log", admin.GetChannelLogHandler)                       //渠道统计
	r.HandleFunc("/living/get_sign_reward_log", admin.GetLivingSignRewardLogHandler)       //获取直播间签到提现记录
	r.HandleFunc("/living/update_sign_reward_log", admin.UpdateLivingSignRewardLogHandler) //更新直播间签到提现记录状态
	r.HandleFunc("/user/balance_change_log", admin.BalanceChangeLogHandler)                //用户资产变化记录
	r.HandleFunc("/living/gift_log", admin.GiftLogHandler)                                 //主播礼物日志
	r.HandleFunc("/living/gift_amount", admin.GiftAmountHandler)                           //主播礼物数量
	r.HandleFunc("/living/gift_update", admin.UpdateUserGiftHandler)                       //修改主播礼物数量
}

func registerUserHandlers() {
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_DZ_LOGIN))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_DZ_VERIFY_CODE))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_DZ_ACCESS_TOKEN_LOGIN))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_DZ_GET_VERIFY_CODE))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_DZ_HEART_BEAT))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_DZ_REGISTER)) //注册
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_DZ_RESET_PASSWORD))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_N_LOGIN))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_DZ_CODE_LOGIN))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_DZ_GET_VERIFY_CODE_BY_EMAIL))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_DZ_RESET_PASSWORD_BY_EMAIL))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_DZ_REGISTER_BY_EMAIL))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_DZ_ROBOT_LOGIN))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MESSAGEID_Disconnection_Report))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageIDRegisterGuest))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageIDAppleRegister))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageIDRegisterLine))

	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_N_LOGIN), user.NLoginHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_LOGIN), dzuser.DZLoginHandler) //游客登录
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_VERIFY_CODE), dzuser.VerifyCode)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MODIFY_USER_AVATAR), dzuser.ModifyUserAvatar)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MODIFY_USER_NICKNAME), dzuser.ModifyUserNickName)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MODIFY_USER_NICKNAME_Request), dzuser.ModifyUserNickNameRequest)

	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_HISTORY_DZ_JIFEN), dzgame.DZHistoryGameResult_JiFen)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_HISTORY_DZ_GLOD), dzgame.DZHistoryGameResult_Glod)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_HISTORY_DZ_ACTIVE), dzgame.DZHistoryGameResult_Active)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_HISTORY_DZ_LIVE), dzgame.DZHistoryGameResult_Live)

	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_SHOP_DISMOND_PRODUCT), shop.Product_Dismond_List)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_SHOP_GLOD_PRODUCT), shop.Product_Glob_List)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_SHOP_BEAN_PRODUCT), shop.ProductBeanList)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_SHOP_PAY), shop.Product_Pay)

	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_GET_USER_INFO), dzuser.DZGetUserInfoHandler)           //获取用户信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ROOM_LIST), dzgame.DZGameListHandler)                  //获取游戏房间列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CREATE_ROOM), dzgame.CreateRoomHandler)                //创建房间
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ENTER_ROOM), dzgame.DZEnterGameHandler)                //进入房间
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_APPLY_SIT_DOWN), dzgame.DZApplySitDownHandler)         //申请坐下
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_LEAVE_ROOM), dzgame.DZLeaveRoomHandler)                //离开房间
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_STAND_UP), dzgame.DZStandUpHandler)                    //站起
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_HOST_START_GAME), dzgame.DZStartGameHandler)           //开始游戏
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_APPLY_TAKE), dzgame.DZApplyTakeHandler)                //申请带入
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_APPLY_TAKE_HOST_ACK), dzgame.DZApplyTakeAnswerHandler) //申请带入房主应答
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_HOST_SETTING), dzgame.DZHostSettingHandler)            //房主修改设置
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_PLAYER_OPERATE_CARDS), dzgame.DZGameOperateHandler)    //玩家操作牌
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CODE_LOGIN), dzuser.WeChatLoginHandler)                //微信登录
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ACCESS_TOKEN_LOGIN), dzuser.AccessTokenLoginHandler)   //微信accessToken登录
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_GET_VERIFY_CODE), dzuser.DZGetVerifyCodeHandler)       //获取验证码 11
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_HEART_BEAT), dzuser.DZHeartBeatHandler)                //心跳
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_REGISTER), dzuser.DZRegisterHandler)                   //注册
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_RESET_PASSWORD), dzuser.DZResetPasswordHandler)        //找回密码
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_dz_user_change_passwd), dzuser.DZChangeUserPassword)      //修改用户密码 新
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CUR_GAME_RESULT), dzgame.CurGameResultHandler)         //实时战绩
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_LAST_GAME_RESULT), dzgame.LastGameResultHandler)       //上手回顾
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_GAME_RESULT_LIST), dzgame.GameResultListHandler)       //历史战绩
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_COMPLETE_USER_INFO), dzuser.DZCompleteUserInfoHandler) //完善用户信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_RESULT), dzgame.DZMatchResultHandler)            //游戏记录信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_HISTORY_ADD_UP), dzgame.HistoryAddUpHandler)           //
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_HISTORY_LIST), dzgame.HistoryListHandler)              //战绩列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_HISTORY_RESULT), dzgame.HistoryResultHandler)          //
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_HISTORY_TABLE_RESULT), dzgame.HistoryTableResultHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ROOM_RESULT_DETAIL), dzgame.DZRoomResultDetail)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ROOM_RESULT_PRE), dzgame.DZShowPreGameResultHandler) //上手 本局牌谱
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_IS_GAMING), dzgame.DZIsGamingHandler)                //
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_BALANCE), dzuser.GetBalanceHandler)                  //获取资产信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CHAT), dzgame.DZChatHandler)                         //聊天请求
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_TABLE_COUNT), dzgame.DZTableCountHandler)            //获取游戏桌子的数量
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_GET_SUBSIDY), dzuser.DZGetSubsidyHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_REGISTER_BY_EMAIL), dzuser.DZRegisterByEmailHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_GET_MATCH_INFO), dzuser.GetMatchRecordHandler)          //游戏中获取游戏统计信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CREATE_CLUB), dzclub.DZCreateClubHandler)               //创建俱乐部
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_LSIT_ADD), dzclub.DZAddClubListHandler)            //加入的俱乐部列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_LSIT_CREATE), dzclub.DZCreateClubListHandler)      //创建的俱乐部列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_APPLY), dzclub.DZApplyAddHandler)                  //申请加入俱乐部
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_APPLY_NOTIFY_ANSWER), dzclub.DZApplyAnswerHandler) //俱乐部主处理申请请求
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_ADD_INFO), dzclub.DZAddClubInfoHandler)            //加入的俱乐部详细信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_CREATE_INFO), dzclub.DZCreateClubInfoHandler)      //创建的俱乐部详细信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_FIND), dzclub.DZFindClubHandler)                   //查找俱乐部
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_DELETE_MEMBER), dzclub.DZDeleteClubMemberHandler)  //删除会员或者是退出俱乐部
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_GAME_LIST), dzclub.DZClubGameListHandler)          //俱乐部牌局列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_SHOW_CARDS), dzgame.DZShowCardsHandler)                 //秀牌
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_PRE_GAME_RESULT), dzgame.DZPreGameResultHandler)        //上手回顾
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_LOOK_LEFT_CENTER_CARD), dzgame.DZLookLeftCenterCardHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_UPGRADE_ACCOUNT_BY_EMAIL), dzgame.DZUpgradeAccountByEmailHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_APPLY_TAKE_MSGS), dzclub.DZGetApplyTakeMsgHandler) //俱乐部获取申请带入的消息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_APPLY_ADD_MSGS), dzclub.DZGetApplyAddMsgHandler)   //俱乐部获取申请加入的消息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_MEMBERS), dzclub.DZGetClubUserHandler)             //俱乐部获取获取成员
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_ADMINS), dzclub.DZGetClubAdminUserHandler)         //俱乐部获取获取管理员
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_OPRATE_ADMINS), dzclub.DZOperateAdminHandler)      //俱乐部添加删除管理员
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_MODIFY_SETTING), dzclub.DZSettingChangeHandler)    //俱乐部设置修改
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_USER_FIND), dzclub.DZFindClubUserHandler)          //俱乐部查找用户
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_GET_ADVERT), dzuser.DZGetAdvertHandler)                 //获取广告页
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_INSURE_BUY_SCORE_REQ), dzgame.DZInSureBuyScoreHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_GET_CLUB_HAVE_MESSAGES), dzclub.DZIsHaveMessageHandler) //获取是否有申请加入的消息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_LANGUAGE_SETTING), dzuser.DZLanguageSettingHandler)     //语言设置
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_PUSH_MSG), dzclub.DZClubPushMsgHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_PUSH_SWITCH), dzclub.DZClubPushSwitchHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ALLIANCE_CREATE), dzclub.AllianceCreateHandler)                  //创建联盟
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ALLIANCE_FIND), dzclub.AllianceFindHandler)                      //查找联盟
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ALLIANCE_ADD_APPLY), dzclub.AllianceApplyHandler)                //加入联盟申请
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ALLIANCE_INFO), dzclub.AllianceInfoHandler)                      //联盟信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ALLIANCE_ADD_APPLY_MSG), dzclub.AllianceApplyMsgHandler)         //联盟申请加入的消息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ALLIANCE_ADD_APPLY_MSG_ANSWER), dzclub.AllianceApplyDealHandler) //联盟申请加入的消息处理
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ALLIANCE_MEMBER_DELETE), dzclub.AllianceMemberDeleteHandler)     //联盟成员删除或者是解散
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ALLIANCE_APPLY_CLOSE), dzclub.AllianceCloseApplyHandler)         //修改联盟申请开关控制
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_REDBAG_HISTORY_LIST), dzuser.DZRedbagHistoryHandler)             //红包历史记录
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_RANKING_LIST), dzuser.DZGetRankingListHandler)                   //获取排名
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_TAKE_OUT), dzgame.DZTakeOutHandler)                              //短牌申请带出积分牌
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_RECOMMEND_CLUB_LIST), dzclub.DZRecommendListHandler)             //推荐俱乐部列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_BUY_COINS), dzuser.DZBuyCoinsHandler)                            //购买金币
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ROBOT_LOGIN), dzuser.DZRobotLoginHandler)                        //机器人登陆
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_PUSH_MSG_TEST), dzclub.DZPushMsgTestHandler)                     //推送测试
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_CREATE), dzgame.CreateMatchHandler)                        //创建比赛
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_ENTER), dzgame.DZEnterMatchHandler)                        //进入比赛
	GetMsgRegistry().RegisterMsg(int32(pb.MessageID_DZ_MATCH_ENTER_GAME), dzgame.DZEnterMatchGameHandler)               //进入比赛游戏桌旁观
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_JOIN), dzgame.DZJoinMatchHandler)                          //报名比赛
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_BACK), dzgame.DZBackMatchHandler)                          //返回比赛
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_BACK_HALL), dzgame.DZBackMatchHallHandler)                 //返回大厅
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_LIST), dzgame.DZMatchListHandler)                          //比赛列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_CANCEL), dzgame.DZCancelJoinMatchHandler)                  //取消报名
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_MANUAL_START), dzgame.DZManualStartMatchHandler)           //手动开始比赛
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_DISMISS), dzgame.DZDismissMatchHandler)                    //解散比赛
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_CANCEL_AUTO_PALY), dzgame.DZMatchCancelAutoPlayHandler)    //比赛取消托管
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_RESULT_SUMMARY), dzgame.MatchResultSummaryHandler)         //比赛战绩统计
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MATCH_REWARD_BAG), dzgame.DZMatchBagListHandler)                 //获取比赛背包
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_LIVE_GEN_TOKEN), live.GenerateTokenHandler)                      //生成声网Token
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_LIVE_ROOM_LIST), live.LiveRoomListHandler)                       //获取主播场里的直播间列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_LIVE_SHOW_REQUEST), live.LiveShowRequestHandler)                 //申请上播
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_LIVE_SHOW_REQUEST_DECIDED), live.LiveShowRequestDecidedHandler)  //管理审核上播
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_LIVE_SHOW_CLOSE), live.LiveShowCloseHandler)                     //主动下播
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_LIVE_SHOW_BREAK), live.LiveShowBreakHandler)                     //管理踢下播
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_LIVING_SWITCH), live.LivingSwitchHandler)                        //直播/非直播切换
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_MICPHONE_SWITCH), live.MicphoneSwitchHandler)                    //麦克打开/关闭切换
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_SEND_GIFT), live.SendGiftHandler)                                //送礼物
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ACTIVITY_LIST), config.Activity_list)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_ANCHOR_SETTING), live.AnchorSettingHandler)                          //主播设置
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_DISABLE_MICPHONE), live.DisableMicphoneHandler)                      //禁麦
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_DISABLE_ALL_MICPHONE_IN_GAME), live.DisableAllMicphoneInGameHandler) //全员禁麦
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_DISABLE_ALL_DEFAULT), live.DisableAllDefaultHandler)                 //默认全体禁麦
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_REMOVE_LIVE_COIN_GAME), live.RemoveLiveCoinGameHandler)              //删除主播金币场牌局
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CHANGE_LIVE_CIN_GAME), live.ChangeLiveCoinGameHandler)               //修改主播金币场牌局
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_PERSONAL_CARD), dzclub.PersonalCard)                                    //个人卡片

	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_EDIT_LAMP), dzclub.EditClubLampHandler)              //俱乐部跑马灯编辑
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_CLUB_MEMBER_REMARK), dzclub.ClubMemberRemarkHandler) //俱乐部添加玩家备注
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_CLUB_GRANT_SCORE), dzclub.GrantScore)                   //俱乐部 发送积分
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_CLUB_RECLAIM_SCORE), dzclub.ReclaimScore)               //俱乐部 回收积分
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_CLUB_FUND_DETAIL), dzclub.FundDetail)                   //俱乐部 资产明细 收发明细
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_CLUB_DIAMOND_TO_COIN), dzclub.ClubDiamondToCoin)        //俱乐部兑换钻石到金币
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_CLUB_GET_REDDOT), dzclub.ClubGetRedDot)                 //俱乐部红点
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_CLUB_SET_REDDOT), dzclub.ClubSetRedDot)                 //俱乐部红点
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_CLUB_OWNER_DATA), dzclub.OwnerData)                     //俱乐部 群主数据
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_CLUB_MEMBER_COUNT), dzclub.ClubMemberCount)             //俱乐部 成员数据

	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_DZ_ONLINE_MUSIC_OP), live.LivingMusicOp)            //音效
	GetMsgRegistry().RegisterMsg(int32(pb.MessageID_DZ_MATCH_DESK_INFO), dzgame.DZMatchDeskInfoHandler) //比赛场牌桌信
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_SHORT_CART_SET), dzgame.ShortCardSetHandler)     //短牌的设置

	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_TASK_QUERY_REWARD), task.QueryReward)            //任务奖励
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_TASK_TAKE_REWARD), task.TakeReward)              //任务奖励
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_TASK_DECREASE_MONEY), task.TaskDecrMoney)        //任务奖励
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_GET_TASK_REDDOT), task.GetTaskRedDot)            //任务红点
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_SET_TASK_REDDOT), task.SetTaskRedDot)            //任务红点
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_TAKE_GAME_REVIEWS), dzgame.TakeGameReviews)      //本局牌谱
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_WINNER_RANK), dzgame.WinnerRank)                 //本局牌谱
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_COLLECTION_SAVE), collection.CollectionSave)     //牌谱收藏
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_COLLECTION_UNSAVE), collection.CollectionUnsave) //牌谱收藏
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_COLLECTION_QUERY), collection.CollectionQuery)   //牌谱收藏
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_REPLAY_FETCH), collection.ReplayFetch)           //回放

	GetMsgRegistry().RegisterMsg(int32(pb.MessageID_DZ_MATCH_DESK_LIST), dzgame.DZMatchDeskListHandler)            //比赛场牌列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageID_DZ_MATCH_USER_LIST), dzgame.DZMatchUserListHandler)            //比赛用户列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageID_DZ_MATCH_STATUS), dzgame.DZGetMatchStatusHandler)              //比赛场状态信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageID_DZ_MATCH_RESULT_LIST), dzgame.MatchResultListHandler)          //比赛场战绩列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageID_DZ_MATCH_RESULT_USER_LIST), dzgame.MatchResultUserListHandler) //比赛场战绩玩家列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageID_DZ_MATCH_CHANGE_LIVER_ID), dzgame.MatchChangeLiverIDHandler)   //比赛场修改主播ID

	//专业版
	GetMsgRegistry().RegisterMsg(int32(pb.MID_PRO_INDEX_SCRORE_PANEL), home.IndexScorePanelHandler) //首页 成绩面板
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_HISTORY_DZ_DELETE), dzgame.DZHistoryGameResult_Del)
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_Disconnection_Report), debug.DisConnectionReport) //客户端断线重连报告。
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_Client_Version_Get_Notify), config.GetNowVersion)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageID_GAME_SETTING_DELAY_TIME), dzgame.GameSettingPlayDelay)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_Game_Setting_Replay_SetStop), dzgame.GameSettingStopReplay) //设置退出复盘

	//======
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Enter), dzclub.HallEnter)                             //进入俱乐部大厅
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Exit), dzclub.HallExit)                               //退出俱乐部大厅
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Send_Message), dzclub.HallMessageSend)                //发送聊天信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Message_Disable_Talk), dzclub.HallMessageDisableTalk) //用户禁言操作
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Mic_UpOrDown), dzclub.HallMicUpOrDown)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Mic_LockOrUnLock), dzclub.HallMicLockOrUnLock)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Mic_DisableOrRelieve), dzclub.HallMicDisableOrRelieve)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Mic_CloseOrOpen), dzclub.HallMicOpenOrClose)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Live_Up), dzclub.HallLiveUp)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Live_Down), dzclub.HallLiveDown)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Live_Apply_List), dzclub.HallLiveApplyList)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Live_Apply_Op), dzclub.HallLiveApplyOp)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Live_LockOrUnLock), dzclub.HallLiveLockOrUnLock)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_OnLine_Member_List), dzclub.HallOnLineMemberList)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Setting), dzclub.HallSetting)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ClubHall_Live_Mic_LockOrUnLock), dzclub.HallLiveMicLockOrUnLock)

	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_DZ_Club_List_All), dzclub.DZUserClubList)

	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_Club_Gift_List), dzclub.HallGiftList) //大厅礼物列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_Club_Gift_Give), dzclub.HallGiftGive) //大厅赠送礼物

	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_ORDER_CREATE), shop.CreateOrder)
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_ORDER_PAY_CONFIRM_WITH_GOOGLE), shop.PayConfirmWithGoogle)
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_ORDER_USER_UN_CONFIRM), shop.UserUnConfirmOrder)
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_ORDER_CHECK_APPLE), shop.PayCheckWithApple)

	GetMsgRegistry().RegisterMsg(int32(pb.MessageID_LIVER_LIST), dzgame.DZLiverListHandler)          //获取主播列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageID_LIVER_GAME_ID), dzgame.DZLiverGameIDHandler)     //获取主播游戏ID
	GetMsgRegistry().RegisterMsg(int32(pb.MessageID_LIVER_NEW_GAME_ID), dzgame.DZGetOneLiverHandler) //获取一个新的主播游戏ID

	//gift
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_GIFT_AMOUNT), gift.GetGiftAmountHandler)             //获取礼物数量
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_GIFT_LIST), gift.GetGiftReceiveListHandler)          //礼物记录
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_GIFT_TAKE_OUT_APPLY), gift.GetGiftTakeOutHandler)    //提现申请
	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_GIFT_TAKE_OUT_LIST), gift.GetGiftTakeOutListHandler) //提现记录

	GetMsgRegistry().RegisterMsg(int32(pb.MESSAGEID_USER_BALACE_REQ), shop.UserBalanceUpdate)               //更新资产
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDRegisterGuest), dzuser.DZRegisterGuestHandler)           //游客注册
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDRegisterLine), dzuser.DZRegisterLineHandler)             //line注册
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDCoinMatchList), dzgame.CoinMatchListHandler)             //获取金币场列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDRankingList), rankinglist.GetRankingHandler)             //获取排行榜
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDShareGame), dzuser.UserShareGameHandler)                 //分享
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDFeedback), hall.FeedbackHandler)                         //反馈
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDSignData), hall.GetSignDataHandler)                      //获取签到信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDSign), hall.SignHandler)                                 //签到
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDSignRedbag), hall.SignRedbagHandler)                     //签到礼包
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDBuyCheapSstatus), hall.CheapGoodsStatusHandler)          //购买特惠商品的状态
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDTaskStatus), hall.TaskHandler)                           //任务状态
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDTaskReward), hall.GetTaskRewardHandler)                  //获取任务奖励
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDAppleRegister), dzuser.DZRegisterAppleHandler)           //苹果用户注册
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDSetInvitor), dzuser.SetInviteUserHandler)                //设置邀请人
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDGetInviteRewardList), dzuser.InviteRewardListHandler)    //获取邀请奖励
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDSysMessage), hall.SysMessageHandler)                     //获取系统消息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDLivingGifter), live.LivingGifterRankHandler)             //获取主播打赏列表
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDGetUserDeskInfo), live.UserDeskInfoHandler)              //获取用户牌桌内信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDLiverForbidUser), live.LiverForbidUserHandler)           //主播屏蔽某用户
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDLiverGetBorbidInfo), live.LiverGetForbidInfoHandler)     //主播获取屏蔽用户信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDLiverInviteUser), live.InviteUserHandler)                //主播邀请用户上坐
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDLiverGifterRank), live.LivingLookerListHandler)          //主播打赏榜单
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDFollowLiver), hall.AddOrDeleteFollowHandler)             //关注或着取消
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDISFollowLiver), hall.IsFollowHandler)                    //是否关注了
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDFollowMe), hall.GetFollowMeHandler)                      //粉丝
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDMyFollow), hall.GetMyFollowHandler)                      //我关注的
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDSetNotice), live.SetLiverNoticeHandler)                  //设置公告
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDGetNotice), live.GetLiverNoticeHandler)                  //获取公告
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDLiverAskContact), live.AskUserContactHandler)            //询问玩家联系方式
	GetMsgRegistry().RegisterMsg(int32(pb.MesasgeIDAnswerLiverContact), live.AnswerLiverContactHandler)     //玩家应答主播的询问
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDLivingTaskStatus), live.TaskLivingHandler)               //直播间任务状态信息
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDLivingSign), live.SignLivingHandler)                     //直播间签到
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDLivingTaskReward), live.TaskRewardLivingHandler)         //直播间任务获得奖励
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDLivingSignReward), live.SignRewardLivingHandler)         //直播间签到任务总奖励
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDDeleteAccount), dzuser.DeleteAccountHandler)             //删除账号
	GetMsgRegistry().RegisterMsg(int32(pb.MessageIDCancelDeleteAccount), dzuser.CancelDeleteAccountHandler) //取消删除账号
}
