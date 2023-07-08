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
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "config":
			CreateDefaultConfig()
		case "version":
			ShowVersion()
		default:
			ShowUsage()
		}

		return

	}

	// 加载配置文件
	app := config.NewCenterApp()
	log.Println("load config success, listen:", app.Config.Listen)

	// 启用gops
	if app.Config.GOPS {
		agent.Listen(agent.Options{})
	}

	// 保存session 的app
	relay_center.SessionServe.SetApp(app)

	// 启动服务
	serverFactory := relay_center.NewServerFactory(app)
	// 根据配置文件模式创建服务
	server, err := serverFactory.CreateServer(app.Config.Mode)
	if nil != err {
		log.Panicln(err)
	}

	log.Printf("start create %s server\n", app.Config.Mode)
	err = server.StartServer()

	log.Println("server is stop: ", err)
}

// 创建服务端配置文件
func CreateDefaultConfig() {
	cfg := config.CreateDefaultCenterConfig()
	out, _ := yaml.Marshal(cfg)
	fmt.Println(string(out))
}

// 创建客户端配置文件
func ShowUsage() {
	fmt.Printf("Usage: %s [config|help|version] \n"+
		"    config: 生成默认配置文件\n"+
		"    help: 显示帮助信息\n"+
		"    version: 显示版本信息\n",
		os.Args[0])
}

func ShowVersion() {
	fmt.Printf("%s\n", Version)
}
