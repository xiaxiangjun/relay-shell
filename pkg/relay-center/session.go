package relay_center

import (
	"encoding/json"
	"fmt"
	"github.com/xiaxiangjun/relay-shell/pkg/common"
	"github.com/xiaxiangjun/relay-shell/pkg/config"
	"github.com/xiaxiangjun/relay-shell/utils"
	"log"
	"time"
)

type sessionServe struct {
	app *config.CenterApp
}

var SessionServe sessionServe

func (self *sessionServe) SetApp(app *config.CenterApp) {
	self.app = app
}

// 处理会话
func (self *sessionServe) Serve(session common.Session) {
	// 保存session会话
	SessionManager().AddSession(session)
	defer SessionManager().RemoveSession(session)

	for {
		// 获取一个stream
		stream, err := session.AcceptStream(self.app.Context)
		if nil != err {
			break
		}

		go self.sessionStreamServe(session, common.NewStreamReader(stream))
	}
}

// 处理一个stream请求
func (self *sessionServe) sessionStreamServe(session common.Session, stream *common.StreamReader) {
	defer stream.Close()

	// step 1: 读取一行请求
	msg, err := stream.ReadMessage(8192)
	if nil != err {
		log.Println(session.RemoteAddr().String(), "parse line error: ", err)
		return
	}

	// 判断用户是否存在
	if false == self.checkUserAuth(msg) {
		stream.WriteError(common.ErrForbidden)
		return
	}

	self.sessionDispatch(session, stream.MoveToOwner(), msg)
}

// 验证用户授权
func (self *sessionServe) checkUserAuth(msg *common.Message) bool {
	// 获取用户密码
	passwd, ok := self.app.Config.Users[msg.User]
	if false == ok {
		log.Println("user", msg.User, "is not exist")
		return false
	}

	// 判断时间是否过期, 且时间不能超过1小时，超过1小就必须禁用
	now := time.Now().Unix()
	if msg.Time < now || msg.Time > (now+3600) {
		log.Println("user", msg.User, "time is expire")
		return false
	}

	// step 2: 对数据进行签名验证
	// hmac-sha1(user+":"+ts+":"+type+":"+uuid, password)
	sign := utils.StringSign(fmt.Sprintf("%s:%d:%s:%s", msg.User, msg.Time, msg.Type, msg.UUID), passwd)
	if sign != msg.Sign {
		log.Println("sign is not same")
		return false
	}

	// 防止短期内被用户盗用
	if false == UUIDCache.Store(fmt.Sprintf("%s:%s", msg.User, msg.Msg), msg.Time) {
		log.Println("user ", msg.User, "uuid has used: ", msg.UUID)
		return false
	}

	return true
}

// 分发处理消息
func (self *sessionServe) sessionDispatch(session common.Session, stream *common.StreamReader, msg *common.Message) {
	defer stream.Close()

	switch msg.Type {
	case common.MsgConnect:
		self.onRequestConnect(session, stream.MoveToOwner(), msg)
	case common.MsgRegister:
		self.onRequestRegister(session, stream.MoveToOwner(), msg)
	default:
		// 写入失败
		stream.WriteError(common.ErrMethodNotAlowed)
	}
}

// 注册请求
func (self *sessionServe) onRequestRegister(session common.Session, stream *common.StreamReader, msg *common.Message) {
	defer stream.Close()

	// 读取设备id
	req := common.RegisterRequest{}
	json.Unmarshal(msg.Payload, req)
	if "" == req.SessionID {
		stream.WriteError(common.ErrNotAcceptable)
		return
	}

	sid := fmt.Sprintf("%s:%s", msg.User, req.SessionID)
	// 添加用户注册
	SessionRegister().AddSession(sid, session)

	// 回应客户端
	stream.WriteMessage(&common.Message{
		Type: common.MsgResponse,
		Code: common.ErrOK.Code(),
		Msg:  common.ErrOK.Msg(),
	}, &common.RegisterResponse{})
}

// 连接请求
func (self *sessionServe) onRequestConnect(session common.Session, stream *common.StreamReader, msg *common.Message) {
	defer stream.Close()

	// 读取设备id
	req := common.ConnectRequest{}
	json.Unmarshal(msg.Payload, req)
	if "" == req.TargetSID {
		stream.WriteError(common.ErrNotAcceptable)
		return
	}

	targetSid := fmt.Sprintf("%s:%s", msg.User, req.TargetSID)
	// 查找用户连接是否存在
	peer, err := self.peerConnect(targetSid, msg)
	if nil != err {
		stream.WriteError(err)
		return
	}

	// 读取回应
	buf, e := peer.ReadLine(8192)
	if nil != e {
		stream.WriteError(common.ErrInternalServerError(e.Error()))
		return
	}

	// 回应客户端
	utils.WriteAll(stream, buf)
	// 交换数据
	utils.IoSwap(stream, peer)
}

func (self *sessionServe) peerConnect(targetSid string, msg *common.Message) (*common.StreamReader, *common.Error) {
	// 查找对端的会话
	peerSession := SessionRegister().FindSession(targetSid)
	if nil == peerSession {
		return nil, common.ErrNotFound
	}

	// 创建一个新的会话
	peerStream, err := peerSession.OpenStream(self.app.Context)
	if nil != err {
		return nil, common.ErrInternalServerError(err.Error())
	}

	// 写入连接请求
	peer := common.NewStreamReader(peerStream)
	err = peer.WriteMessage(msg, nil)
	if nil != err {
		return nil, common.ErrInternalServerError(err.Error())
	}

	return peer, nil
}
