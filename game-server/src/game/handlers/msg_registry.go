package handlers

import (
	"game/domain/dzgame"
	"game/domain/task"
	"game/domain/user"
	domainUser "game/domain/user"
	"game/server"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

type MsgRegistry struct {
	registry           map[int32]func(msg *server.ClientMsg, sess *server.Session) []byte
	unLoginMsgRegistry map[int32]bool
	mu                 sync.RWMutex
}

func (registry *MsgRegistry) RegisterMsg(msgId int32, f func(msg *server.ClientMsg, sess *server.Session) []byte) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.registry[msgId] = f
}

func (registry *MsgRegistry) RegisterUnLoginMsg(msgId int32) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.unLoginMsgRegistry[msgId] = true
}

func (registry *MsgRegistry) isUnLoginMsg(msgId int32) bool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	_, ok := registry.unLoginMsgRegistry[msgId]
	return ok
}

func (registry *MsgRegistry) getHandler(msgId int32) func(msg *server.ClientMsg, sess *server.Session) []byte {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	return registry.registry[msgId]
}

func (registry *MsgRegistry) DispatchMsg(msg *server.ClientMsg, sess *server.Session) []byte {
	start := time.Now()
	defer func() {
		elapseTime := time.Since(start)
		if elapseTime.Seconds() > 0.1 {
			p := user.GetPlayer(sess.Data)
			if p != nil {
				//user.SaveSlowMsg(p.User.UserId, msg.MsgId, start, elapseTime.String())
			}
		}
	}()

	if !registry.isUnLoginMsg(msg.MsgId) {
		if user.GetPlayer(sess.Data) == nil {
			glog.Info("===>玩家未登录msgId:", msg.MsgId)
			return nil
		}
	}

	if player := domainUser.GetPlayer(sess.Data); player != nil {
		go task.GrantDailyReward(player.User.UserId, task.DAILY_LANDING)
		go task.RewardLandingTotal(player.User.UserId)
	}

	f := registry.getHandler(msg.MsgId)
	if f == nil {
		glog.Error("msgId:  has no handler")
		return nil
	}

	if player := domainUser.GetPlayer(sess.Data); player != nil {
		gameId := dzgame.GetEnterGameManager().GetEnterGame(player.User.UserId)
		if gameId != 0 {
			dzgame.GetGameManager().AppendReplayMessage(gameId, player.User.UserId, nil, msg.MsgId, msg.MsgBody.(string))
		}
	}

	rep := f(msg, sess)

	if player := domainUser.GetPlayer(sess.Data); player != nil {
		gameId := dzgame.GetEnterGameManager().GetEnterGame(player.User.UserId)
		if gameId != 0 && rep != nil {
			dzgame.GetGameManager().AppendReplayMessage(gameId, "", []string{player.User.UserId}, msg.MsgId, string(rep))
		}
	}
	return rep
}

func (registry *MsgRegistry) RegisterHandlers(r *mux.Router) {
	registerHandlers(r)
}

var registry *MsgRegistry

func init() {
	registry = &MsgRegistry{
		registry:           make(map[int32]func(msg *server.ClientMsg, sess *server.Session) []byte),
		unLoginMsgRegistry: make(map[int32]bool),
		mu:                 sync.RWMutex{},
	}
}

func GetMsgRegistry() *MsgRegistry {
	return registry
}
