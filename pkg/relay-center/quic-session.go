package relay_center

import (
	"context"
	"github.com/quic-go/quic-go"
	"github.com/xiaxiangjun/relay-shell/pkg/common"
	"net"
)

type quicSession struct {
	session quic.Connection
}

func (self *quicSession) Close() error {
	return self.session.CloseWithError(0, "")
}

func (self *quicSession) RemoteAddr() net.Addr {
	return self.session.RemoteAddr()
}

func (self *quicSession) OpenStream(ctx context.Context) (common.Stream, error) {
	stream, err := self.session.OpenStream()
	if nil != err {
		return nil, err
	}

	return &quicStream{
		stream: stream,
	}, nil
}

func (self *quicSession) AcceptStream(ctx context.Context) (common.Stream, error) {
	stream, err := self.session.AcceptStream(ctx)
	if nil != err {
		return nil, err
	}

	return &quicStream{
		stream: stream,
	}, nil
}
