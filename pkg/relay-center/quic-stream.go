package relay_center

import (
	"github.com/quic-go/quic-go"
	"sync/atomic"
	"time"
)

type quicStream struct {
	stream         quic.Stream
	lastActiveTime *int64
}

func (self *quicStream) Read(buf []byte) (int, error) {
	n, err := self.stream.Read(buf)
	// 保存活跃时间
	if nil == err {
		atomic.StoreInt64(self.lastActiveTime, time.Now().Unix())
	}

	return n, err
}

func (self *quicStream) Write(buf []byte) (int, error) {
	n, err := self.stream.Write(buf)
	// 保存活跃时间
	if nil == err {
		atomic.StoreInt64(self.lastActiveTime, time.Now().Unix())
	}

	return n, err
}

func (self *quicStream) Close() error {
	return self.stream.Close()
}
