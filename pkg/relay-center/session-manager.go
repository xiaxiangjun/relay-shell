package relay_center

import (
	"github.com/xiaxiangjun/relay-shell/pkg/common"
	"log"
	"sync"
	"time"
)

type sessionManager struct {
	locker   sync.Mutex
	sessions map[string]common.Session
}

var sessionMgr = &sessionManager{}

func init() {
	go sessionMgr.runCleanSession()
}

func SessionManager() *sessionManager {
	return sessionMgr
}

// 添加会话
func (self *sessionManager) AddSession(session common.Session) {
	self.locker.Lock()
	defer self.locker.Unlock()

	// 删除旧的会话
	oldSession, ok := self.sessions[session.RemoteAddr().String()]
	if ok || oldSession != session {
		self.closeSession(oldSession)
	}

	// 保存新会话
	self.sessions[session.RemoteAddr().String()] = session
}

// 删除会话
func (self *sessionManager) RemoveSession(session common.Session) {
	self.locker.Lock()
	defer self.locker.Unlock()

	// 验证当前会话是否有效
	oldSession, ok := self.sessions[session.RemoteAddr().String()]
	if false == ok || oldSession != session {
		return
	}

	delete(self.sessions, session.RemoteAddr().String())
	// 执行关闭动作
	self.closeSession(oldSession)
}

func (self *sessionManager) runCleanSession() {
	tick := time.NewTicker(time.Second * 5)
	defer tick.Stop()

	for {
		<-tick.C
		// 执行清理
		self.doCleanSession()
	}
}

func (self *sessionManager) doCleanSession() {
	self.locker.Lock()
	defer self.locker.Unlock()

	ts := time.Now().Unix() - 15
	remove := make([]string, 0, len(self.sessions))

	// 查找过期的会话
	for addr, session := range self.sessions {
		if session.LastActiveTime() < ts {
			remove = append(remove, addr)
			// 关闭会话
			self.closeSession(session)
		}
	}

	// 删除过期的会话
	for _, addr := range remove {
		delete(self.sessions, addr)
	}
}

// 关闭session会话
func (self *sessionManager) closeSession(session common.Session) {
	if nil == session {
		return
	}

	go func() {
		addr := session.RemoteAddr().String()

		defer func() {
			if err := recover(); nil != err {
				log.Println("close session error: ", addr, err)
			}
		}()

		session.Close()
	}()
}
