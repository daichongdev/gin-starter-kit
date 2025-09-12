//go:build wireinject
// +build wireinject

package di

import (
	"gin-demo/controller"
	"gin-demo/repository"
	"gin-demo/service"
	"github.com/google/wire"
)

// RepositorySet Repository 层的 Provider 集合
var RepositorySet = wire.NewSet(
	repository.NewUserRepository,
)

// ServiceSet Service 层的 Provider 集合
var ServiceSet = wire.NewSet(
	service.NewUserService,
	service.NewAuthService,
	service.NewEmailService,
)

// ControllerSet Controller 层的 Provider 集合
var ControllerSet = wire.NewSet(
	controller.NewUserController,
	controller.NewAuthController,
	controller.NewEmailController,
)

// AllSet 所有 Provider 的集合
var AllSet = wire.NewSet(
	RepositorySet,
	ServiceSet,
	ControllerSet,
)

// Container 应用容器
type Container struct {
	UserController  *controller.UserController
	AuthController  *controller.AuthController
	EmailController *controller.EmailController
}

// InitializeContainer 初始化应用容器
func InitializeContainer() *Container {
	wire.Build(
		AllSet,
		wire.Struct(new(Container), "*"),
	)
	return nil
}
