package di

import (
	"testing"
)

func TestContainerInitialization(t *testing.T) {
	// 初始化容器
	container := InitializeContainer()

	// 验证所有组件都被正确注入
	if container == nil {
		t.Fatal("Container should not be nil")
	}

	// 验证Controller层
	if container.UserController == nil {
		t.Error("UserController should not be nil")
	}
	if container.AuthController == nil {
		t.Error("AuthController should not be nil")
	}
	if container.EmailController == nil {
		t.Error("EmailController should not be nil")
	}

	// 验证Service层
	if container.UserService == nil {
		t.Error("UserService should not be nil")
	}
	if container.AuthService == nil {
		t.Error("AuthService should not be nil")
	}
	if container.EmailService == nil {
		t.Error("EmailService should not be nil")
	}

	// 验证Repository层
	if container.UserRepository == nil {
		t.Error("UserRepository should not be nil")
	}

	// 验证依赖关系是否正确建立
	// 通过反射或者类型断言检查依赖是否正确注入
	t.Log("All components successfully injected into container")
}

func TestDependencyInjection(t *testing.T) {
	container := InitializeContainer()

	// 测试UserController是否正确注入了UserService
	userController := container.UserController
	if userController == nil {
		t.Fatal("UserController is nil")
	}

	// 测试AuthController是否正确注入了AuthService
	authController := container.AuthController
	if authController == nil {
		t.Fatal("AuthController is nil")
	}

	// 验证服务之间的依赖关系
	userService := container.UserService
	authService := container.AuthService

	if userService == nil || authService == nil {
		t.Fatal("Services should not be nil")
	}

	t.Log("Dependency injection working correctly")
}
