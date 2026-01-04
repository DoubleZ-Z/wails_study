package main

import (
	"context"
	"embed"
	"wails_study/project"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	manager := project.ServiceManager{}
	httpService := project.NewHttpService(":8088")
	manager.AddService(httpService)

	tcpService := project.NewTcpService(":11402")
	manager.AddService(tcpService)

	if err := manager.StartAll(); err != nil {
		println("启动TCP/HTTP服务失败：" + err.Error())
	}

	// 启动 Wails 应用
	if err := wails.Run(&options.App{
		Title:  "wails_study",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:             app.getMenu(),
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown: func(ctx context.Context) {
			println("正在关闭全部服务...")
			manager.StopAll()
		},
		Bind: []interface{}{
			app,
		},
	}); err != nil {
		println("Error:", err.Error())
	}
}
