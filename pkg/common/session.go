package common

import (
	"context"
	"net"
)

type Session interface {
	Close() error
	LastActiveTime() int64
	RemoteAddr() net.Addr
	OpenStream(ctx context.Context) (Stream, error)
	AcceptStream(ctx context.Context) (Stream, error)
}
