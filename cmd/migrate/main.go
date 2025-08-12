package main

import (
	"gin-demo/config"
	"gin-demo/database"
	"gin-demo/model"
	"gin-demo/pkg/logger"
	"gin-demo/pkg/migration"
)

func main() {
	// 初始化配置
	config.InitConfig()
	cfg := config.GetConfig()

	// 初始化日志
	logger.InitLogger(cfg)

	// 初始化数据库
	database.InitDB()

	// 创建迁移管理器
	manager := migration.NewManager(model.Registry)

	// 运行迁移命令
	manager.RunCommand()
}
