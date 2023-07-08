package main

import (
	"fmt"
	"github.com/google/gops/agent"
	"github.com/xiaxiangjun/relay-shell/pkg/config"
	relay_center "github.com/xiaxiangjun/relay-shell/pkg/relay-center"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func main() {
	// 输出配置文件
	if len(os.Args) > 1 && os.Args[1] == "config" {
		ShowDefaultConfig()
		return
	}

	// 加载配置文件
	app := config.NewCenterApp()
	log.Println("load config success: port=", app.Config.Listen)

	// 启用gops
	if app.Config.GOPS {
		agent.Listen(agent.Options{})
	}

	// 保存session 的app
	relay_center.SessionServe.SetApp(app)

	// 启动服务
	serverFactory := relay_center.NewServerFactory(app)
	server, err := serverFactory.CreateServer(app.Config.Mode)
	if nil != err {
		log.Panicln(err)
	}

	err = server.StartServer()
	log.Println("pkg stop: ", err)
}

func ShowDefaultConfig() {
	cfg := config.CreateDefaultCenterConfig()
	out, _ := yaml.Marshal(cfg)
	fmt.Println(string(out))
}
