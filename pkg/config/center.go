package config

import (
	"github.com/xiaxiangjun/relay-shell/utils"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type CenterConfig struct {
	Mode   string            `yaml:"mode"`
	Listen int               `yaml:"listen"`
	Users  map[string]string `yaml:"users"`
	GOPS   bool              `yaml:"gops"`
	Tls    CertConfig        `yaml:"tls"`
}

type CenterApp struct {
	Config  *CenterConfig
	Context context.Context
}

// NewCenterApp 创建中心服务
func NewCenterApp() *CenterApp {
	// 读取配置文件
	buf, err := os.ReadFile(utils.GetExePath("relay-center.yaml"))
	if nil != err {
		log.Panicln(err)
	}

	cfg := &CenterConfig{
		Mode:   "quic",
		Listen: 1986,
		GOPS:   false,
	}
	err = yaml.Unmarshal(buf, cfg)
	if nil != err {
		log.Panicln(err)
	}

	return &CenterApp{
		Config:  cfg,
		Context: context.Background(),
	}
}

// 创建默认配置
func CreateDefaultCenterConfig() *CenterConfig {
	cert, key, _ := utils.CreateCertificate()

	return &CenterConfig{
		Listen: 1986,
		Users: map[string]string{
			"test": "0316@test",
		},
		Tls: CertConfig{
			Cert: cert,
			Key:  key,
		},
	}
}
