package main

import (
	"gin-demo/pkg/app"
	"gin-demo/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 创建应用实例
	appInstance := app.New()

	// 初始化应用
	if err := appInstance.Initialize(); err != nil {
		panic(err)
	}

	// 运行应用
	if err := appInstance.Run(); err != nil {
		logger.Fatal("Application failed to run", zap.Error(err))
	}
}
