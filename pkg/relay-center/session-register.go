package relay_center

import (
	"github.com/xiaxiangjun/relay-shell/pkg/common"
	"sync"
)

type sessionRegister struct {
	locker   sync.Mutex
	sessions map[string]common.Session
}

var sessionReg = &sessionRegister{}

func SessionRegister() *sessionRegister {
	return sessionReg
}

func (self *sessionRegister) AddSession(id string, session common.Session) {
	self.locker.Lock()
	defer self.locker.Unlock()

	// 删除旧的会话
	oldSession, ok := self.sessions[id]
	if ok || oldSession != session {
		sessionMgr.RemoveSession(oldSession)
	}

	// 保存新会话
	self.sessions[id] = session
}

func (self *sessionRegister) RemoveSession(id string, session common.Session) {
	self.locker.Lock()
	defer self.locker.Unlock()

	// 验证当前会话是否有效
	oldSession, ok := self.sessions[id]
	if false == ok || oldSession != session {
		return
	}

	delete(self.sessions, id)
	// 执行关闭动作
	sessionMgr.RemoveSession(oldSession)
}

func (self *sessionRegister) FindSession(id string) common.Session {
	self.locker.Lock()
	defer self.locker.Unlock()

	session, ok := self.sessions[id]
	if ok {
		return session
	}

	return nil
}
