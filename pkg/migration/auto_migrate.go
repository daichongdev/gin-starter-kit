package migration

import (
	"fmt"
	"gin-demo/database"
	"gin-demo/pkg/logger"
	"reflect"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ModelInterface 模型接口
type ModelInterface interface {
	TableName() string
}

// ModelRegistry 模型注册表
type ModelRegistry struct {
	models []interface{}
}

// NewModelRegistry 创建模型注册表
func NewModelRegistry() *ModelRegistry {
	return &ModelRegistry{
		models: make([]interface{}, 0),
	}
}

// getDB 获取数据库实例
func (mr *ModelRegistry) getDB() *gorm.DB {
	if database.DB == nil {
		logger.Error("数据库未初始化")
		return nil
	}
	return database.DB
}

// Register 注册模型
func (mr *ModelRegistry) Register(model interface{}) {
	mr.models = append(mr.models, model)
	logger.Info("注册模型", zap.String("model", getModelName(model)))
}

// AutoMigrate 自动迁移所有注册的模型
func (mr *ModelRegistry) AutoMigrate() error {
	db := mr.getDB()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	logger.Info("开始自动迁移", zap.Int("models_count", len(mr.models)))

	for _, model := range mr.models {
		if err := mr.migrateModel(model, db); err != nil {
			return fmt.Errorf("迁移模型 %s 失败: %v", getModelName(model), err)
		}
	}

	logger.Info("自动迁移完成")
	return nil
}

// migrateModel 迁移单个模型
func (mr *ModelRegistry) migrateModel(model interface{}, db *gorm.DB) error {
	modelName := getModelName(model)
	tableName := getTableName(model)

	logger.Info("迁移模型",
		zap.String("model", modelName),
		zap.String("table", tableName))

	// 检查表是否存在
	if !db.Migrator().HasTable(model) {
		logger.Info("创建新表", zap.String("table", tableName))
	} else {
		logger.Info("更新现有表", zap.String("table", tableName))
	}

	// 执行自动迁移
	if err := db.AutoMigrate(model); err != nil {
		logger.Error("模型迁移失败",
			zap.String("model", modelName),
			zap.Error(err))
		return err
	}

	// 创建索引
	if err := mr.createIndexes(model, db); err != nil {
		logger.Error("创建索引失败",
			zap.String("model", modelName),
			zap.Error(err))
		return err
	}

	logger.Info("模型迁移成功", zap.String("model", modelName))
	return nil
}

// createIndexes 为模型创建索引
func (mr *ModelRegistry) createIndexes(model interface{}, db *gorm.DB) error {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		gormTag := field.Tag.Get("gorm")

		// 检查是否需要创建索引
		if strings.Contains(gormTag, "index") {
			fieldName := getFieldName(field)
			if !db.Migrator().HasIndex(model, fieldName) {
				logger.Info("创建索引",
					zap.String("table", getTableName(model)),
					zap.String("field", fieldName))

				if err := db.Migrator().CreateIndex(model, fieldName); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// CheckTableChanges 检查表结构变化
func (mr *ModelRegistry) CheckTableChanges() error {
	db := mr.getDB()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	logger.Info("检查表结构变化")

	for _, model := range mr.models {
		if err := mr.checkModelChanges(model, db); err != nil {
			return err
		}
	}

	return nil
}

// checkModelChanges 检查单个模型的变化
func (mr *ModelRegistry) checkModelChanges(model interface{}, db *gorm.DB) error {
	modelName := getModelName(model)
	tableName := getTableName(model)

	// 检查表是否存在
	if !db.Migrator().HasTable(model) {
		logger.Info("发现新表", zap.String("table", tableName))
		return nil
	}

	// 检查字段变化
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldName := getFieldName(field)

		// 跳过嵌入的gorm.Model字段
		if field.Anonymous && field.Type.Name() == "Model" {
			continue
		}

		// 检查字段是否存在
		if !db.Migrator().HasColumn(model, fieldName) {
			logger.Info("发现新字段",
				zap.String("model", modelName),
				zap.String("field", fieldName))
		}
	}

	return nil
}

// DropTables 删除所有表
func (mr *ModelRegistry) DropTables() error {
	db := mr.getDB()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	logger.Info("删除所有表")

	for _, model := range mr.models {
		tableName := getTableName(model)
		if db.Migrator().HasTable(model) {
			logger.Info("删除表", zap.String("table", tableName))
			if err := db.Migrator().DropTable(model); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetTableStatus 获取表状态
func (mr *ModelRegistry) GetTableStatus() []TableStatus {
	db := mr.getDB()
	if db == nil {
		logger.Error("数据库未初始化，无法获取表状态")
		return []TableStatus{}
	}

	var status []TableStatus

	for _, model := range mr.models {
		tableName := getTableName(model)
		exists := db.Migrator().HasTable(model)

		var recordCount int64
		if exists {
			db.Model(model).Count(&recordCount)
		}

		status = append(status, TableStatus{
			ModelName:   getModelName(model),
			TableName:   tableName,
			Exists:      exists,
			RecordCount: recordCount,
		})
	}

	return status
}

// TableStatus 表状态
type TableStatus struct {
	ModelName   string `json:"model_name"`
	TableName   string `json:"table_name"`
	Exists      bool   `json:"exists"`
	RecordCount int64  `json:"record_count"`
}

// 辅助函数

// getModelName 获取模型名称
func getModelName(model interface{}) string {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	return modelType.Name()
}

// getTableName 获取表名
func getTableName(model interface{}) string {
	if tableNamer, ok := model.(ModelInterface); ok {
		return tableNamer.TableName()
	}

	// 如果没有实现TableName方法，使用默认规则
	modelName := getModelName(model)
	return strings.ToLower(modelName) + "s"
}

// getFieldName 获取字段名
func getFieldName(field reflect.StructField) string {
	gormTag := field.Tag.Get("gorm")

	// 从gorm标签中提取column名
	if strings.Contains(gormTag, "column:") {
		parts := strings.Split(gormTag, ";")
		for _, part := range parts {
			if strings.HasPrefix(part, "column:") {
				return strings.TrimPrefix(part, "column:")
			}
		}
	}

	// 使用字段名的snake_case形式
	return toSnakeCase(field.Name)
}

// toSnakeCase 转换为snake_case
func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
