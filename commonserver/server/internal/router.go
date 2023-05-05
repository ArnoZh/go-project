package internal

import (
	"commonserver/base/chanrpc"
	"commonserver/pb"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

// Router router
type Router struct {
	hooks     sync.Map
	charProxy sync.Map // 角色代理
}

// agentHook agent msg hook
type agentHook func(agent *DBAgent, msg interface{})

func (router *Router) setHook(msg interface{}, hook agentHook) {
	msgID := chanrpc.MsgID(msg)
	_, ok := router.hooks.Load(msgID)
	if ok {
		logrus.Fatalf("router setHook: msg %v already registered", chanrpc.MsgName(msg))
		return
	}
	router.hooks.Store(msgID, hook)
}

func (router *Router) getHook(msgID uint32) agentHook {
	h, ok := router.hooks.Load(msgID)
	if !ok {
		return nil
	}
	return h.(agentHook)
}

func (router *Router) addCharProxy(srcCharID, dstCharID int64) {
	router.charProxy.Store(srcCharID, dstCharID)
}

func (router *Router) delCharProxy(srcCharID int64) {
	router.charProxy.Delete(srcCharID)
}

func (router *Router) listCharProxy() (list []string) {
	router.charProxy.Range(func(key, value interface{}) bool {
		list = append(list, fmt.Sprintf("[src:%v dst:%v]", key.(int64), value.(int64)))
		return true
	})

	return
}

func (router *Router) getCharProxy(srcCharID int64) (int64, bool) {
	id, ok := router.charProxy.Load(srcCharID)
	if !ok {
		return 0, false
	}
	return id.(int64), true
}

func (dbm *MysqlDB) initRouter() {
	dbm.router = &Router{}
	dbm.router.setHook((*pb.TestMsgAck)(nil), handleTestMsgAck)
}
