# Gin Starter Kit

🚀 一个轻量级的 Gin Web 框架脚手架，开箱即用，快速构建 RESTful API

## ✨ 特性

- 🏗️ **清晰的分层架构** - Controller/Service/Repository 模式
- 🔐 **JWT 认证** - 完整的用户注册/登录系统
- 📝 **统一响应格式** - 标准化的 API 响应结构
- 🔍 **链路追踪日志** - 自动关联请求和数据库操作日志
- ⚙️ **配置管理** - 基于 YAML 的配置文件
- 🗄️ **数据库支持** - MySQL + Redis
- 🛡️ **中间件支持** - 访问日志、错误处理、CORS 等
- 📊 **健康检查** - 内置健康检查端点

## 🛠️ 技术栈

- **框架**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **数据库**: MySQL + Redis
- **认证**: JWT
- **日志**: [Zap](https://github.com/uber-go/zap)
- **配置**: [Viper](https://github.com/spf13/viper)

## 📁 项目结构
```
gin-starter-kit/
├── config/          # 配置文件和配置结构
├── controller/      # 控制器层
├── database/        # 数据库初始化
├── model/          # 数据模型
├── pkg/            # 公共包
│   ├── auth/       # 认证相关
│   ├── logger/     # 日志系统
│   └── middleware/ # 中间件
├── repository/     # 数据访问层
├── router/         # 路由配置
├── service/        # 业务逻辑层
├── logs/           # 日志文件
└── main.go         # 程序入口
```