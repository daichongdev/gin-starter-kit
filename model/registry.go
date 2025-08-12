package model

import (
	"gin-demo/pkg/migration"
)

var Registry *migration.ModelRegistry

// init 初始化模型注册表
func init() {
	Registry = migration.NewModelRegistry()
	registerModels()
}

// registerModels 注册需要自动迁移的模型
func registerModels() {
	Registry.Register(&User{})
}
