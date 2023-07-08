package relay_center

import (
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/xiaxiangjun/relay-shell/pkg/config"
	"log"
)

type quicServer struct {
	app *config.CenterApp
}

func (self *quicServer) StartServer() error {
	// 构建监听地址
	addr := fmt.Sprintf("0.0.0.0:%d", self.app.Config.Listen)
	// 加载 TLS 证书和密钥
	cert, err := tls.X509KeyPair([]byte(self.app.Config.Tls.Cert), []byte(self.app.Config.Tls.Key))
	if err != nil {
		return fmt.Errorf("load tls cert error: %s", err.Error())
	}

	// 配置 TLS 握手参数
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"relay-shell"},
	}

	// 配置QUIC参数
	quicConfig := &quic.Config{
		MaxIdleTimeout: 15,
	}

	// 启动quic监听
	listener, err := quic.ListenAddr(addr, tlsConfig, quicConfig)
	if err != nil {
		return fmt.Errorf("quic listen error: %s", err.Error())
	}

	log.Println("listen QUIC success: ", addr)
	return self.Serve(listener)
}

func (self *quicServer) Serve(listener *quic.Listener) error {

	for {
		session, err := listener.Accept(self.app.Context)
		if err != nil {
			return fmt.Errorf("quic accept error: %s", err.Error())
		}

		go SessionServe.Serve(&quicSession{
			session: session,
		})
	}

	return nil
}
