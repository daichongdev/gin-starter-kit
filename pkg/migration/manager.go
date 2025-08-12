package migration

import (
	"fmt"
	"gin-demo/pkg/logger"
	"os"

	"go.uber.org/zap"
)

// Manager 迁移管理器
type Manager struct {
	registry *ModelRegistry
}

// NewManager 创建迁移管理器
func NewManager(registry *ModelRegistry) *Manager {
	return &Manager{
		registry: registry,
	}
}

// RunCommand 运行迁移命令
func (m *Manager) RunCommand() {
	if len(os.Args) < 2 {
		m.printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "migrate":
		if err := m.registry.AutoMigrate(); err != nil {
			logger.Error("自动迁移失败", zap.Error(err))
			os.Exit(1)
		}
		fmt.Println("✅ 自动迁移完成")

	case "check":
		if err := m.registry.CheckTableChanges(); err != nil {
			logger.Error("检查表变化失败", zap.Error(err))
			os.Exit(1)
		}
		fmt.Println("✅ 表结构检查完成")

	case "status":
		m.showStatus()

	case "drop":
		if err := m.registry.DropTables(); err != nil {
			logger.Error("删除表失败", zap.Error(err))
			os.Exit(1)
		}
		fmt.Println("✅ 所有表已删除")

	case "fresh":
		// 删除并重新创建所有表
		if err := m.registry.DropTables(); err != nil {
			logger.Error("删除表失败", zap.Error(err))
			os.Exit(1)
		}
		if err := m.registry.AutoMigrate(); err != nil {
			logger.Error("重新创建表失败", zap.Error(err))
			os.Exit(1)
		}
		fmt.Println("✅ 数据库已重置")

	default:
		m.printUsage()
	}
}

// showStatus 显示表状态
func (m *Manager) showStatus() {
	status := m.registry.GetTableStatus()
	
	fmt.Println("📊 数据库表状态:")
	fmt.Println("========================================")
	fmt.Printf("%-20s %-20s %-10s %-10s\n", "模型", "表名", "状态", "记录数")
	fmt.Println("----------------------------------------")
	
	for _, s := range status {
		statusStr := "❌ 不存在"
		if s.Exists {
			statusStr = "✅ 存在"
		}
		
		fmt.Printf("%-20s %-20s %-10s %-10d\n", 
			s.ModelName, 
			s.TableName, 
			statusStr, 
			s.RecordCount)
	}
}

// printUsage 打印使用说明
func (m *Manager) printUsage() {
	fmt.Println("🚀 数据库迁移工具")
	fmt.Println("========================================")
	fmt.Println("命令:")
	fmt.Println("  migrate  - 自动迁移所有模型")
	fmt.Println("  check    - 检查表结构变化")
	fmt.Println("  status   - 显示所有表状态")
	fmt.Println("  drop     - 删除所有表")
	fmt.Println("  fresh    - 删除并重新创建所有表")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  go run cmd/migrate/main.go migrate")
	fmt.Println("  go run cmd/migrate/main.go status")
}