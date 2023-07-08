package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/xiaxiangjun/relay-shell/utils"
	"io"
)

type StreamReader struct {
	stream      Stream
	cacheBuffer []byte
}

func NewStreamReader(stream Stream) *StreamReader {
	return &StreamReader{
		stream: stream,
	}
}

// 移动到新的所有者
func (self *StreamReader) MoveToOwner() *StreamReader {
	// 构建新的所有者
	other := &StreamReader{
		stream:      self.stream,
		cacheBuffer: self.cacheBuffer,
	}

	// 清空自已的所有权
	self.stream = nil
	self.cacheBuffer = nil

	return other
}

func (self *StreamReader) ReadMessage(max int) (*Message, error) {
	buf, err := self.ReadLine(max)
	if nil != err {
		return nil, err
	}

	msg := &Message{}
	err = json.Unmarshal(buf, msg)
	if nil != err {
		return nil, err
	}

	return msg, nil
}

func (self *StreamReader) WriteError(err *Error) error {
	return self.WriteMessage(&Message{
		Type: MsgResponse,
		Code: err.Code(),
		Msg:  err.Msg(),
	}, nil)
}

func (self *StreamReader) WriteMessage(msg *Message, payload interface{}) error {
	// 构建发送数据
	if nil != payload {
		msg.Payload, _ = json.Marshal(payload)
	}

	buf, _ := json.Marshal(msg)
	// 添加换行符
	buf = append(buf, '\r', '\n')

	// 写入全部数据
	return utils.WriteAll(self.stream, buf)
}

func (self *StreamReader) ReadLine(max int) ([]byte, error) {
	if nil == self.stream {
		return nil, errors.New("owner is moved")
	}

	pos := 0
	index := -1
	cache := make([]byte, max)
	for pos < max {
		// 读取一次缓存
		n, err := self.stream.Read(cache[pos:])
		if nil != err {
			return nil, err
		}

		// 游标移动
		pos += n
		// 读取到了换行符
		index = bytes.Index(cache[:pos], []byte("\r\n"))
		if index >= 0 {
			break
		}
	}

	// 判断是否找到结束符
	if index < 0 {
		return nil, io.EOF
	}

	// 复制多出的数据到缓存中
	if index+2 < pos {
		self.cacheBuffer = cache[index+2 : pos]
	}

	return cache[:index+2], nil
}

func (self *StreamReader) Read(buf []byte) (int, error) {
	if nil == self.stream {
		return 0, errors.New("owner is moved")
	}

	if len(self.cacheBuffer) > 0 {
		// 计算数据大小
		n := len(buf)
		if n > len(self.cacheBuffer) {
			n = len(self.cacheBuffer)
		}

		// 复制缓存到buf中
		copy(buf[:n], self.cacheBuffer[n:])
		// 截断缓存
		self.cacheBuffer = self.cacheBuffer[n:]

		return n, nil
	}

	return self.stream.Read(buf)
}

func (self *StreamReader) Write(buf []byte) (int, error) {
	if nil == self.stream {
		return 0, errors.New("owner is moved")
	}

	return self.stream.Write(buf)
}

func (self *StreamReader) Close() error {
	if nil == self.stream {
		return errors.New("owner is moved")
	}

	return self.stream.Close()
}
