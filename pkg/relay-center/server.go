package relay_center

import (
	"errors"
	"github.com/xiaxiangjun/relay-shell/pkg/common"
	"github.com/xiaxiangjun/relay-shell/pkg/config"
)

type serverFactory struct {
	app *config.CenterApp
}

// 创建服务工厂
func NewServerFactory(app *config.CenterApp) *serverFactory {
	return &serverFactory{
		app: app,
	}
}

// 创建服务
func (self *serverFactory) CreateServer(mode string) (common.Server, error) {
	switch mode {
	case "quic":
		return &quicServer{
			app: self.app,
		}, nil
	}

	return nil, errors.New("not support")
}
