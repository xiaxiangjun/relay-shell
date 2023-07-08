package common

import (
	"context"
	"net"
)

type Session interface {
	Close() error
	RemoteAddr() net.Addr
	OpenStream(ctx context.Context) (Stream, error)
	AcceptStream(ctx context.Context) (Stream, error)
}
