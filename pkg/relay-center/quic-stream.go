package relay_center

import (
	"github.com/quic-go/quic-go"
)

type quicStream struct {
	stream quic.Stream
}

func (self *quicStream) Read(buf []byte) (int, error) {
	return self.stream.Read(buf)
}

func (self *quicStream) Write(buf []byte) (int, error) {
	return self.stream.Write(buf)
}

func (self *quicStream) Close() error {
	return self.stream.Close()
}
