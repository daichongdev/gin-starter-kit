package test

import (
	"fmt"
	"gin-demo/model"
	"gin-demo/service"
	"testing"
)

func TestCreateUser(t *testing.T) {
	cleanup := SetupTest(t)
	defer cleanup() // 确保测试结束后清理资源

	userService := service.NewUserService()
	userInfo, err := userService.CreateUser(&model.CreateUserRequest{
		Name:  "test",
		Email: "daichongweb@foxmail.com",
		Age:   10,
		Phone: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(userInfo)
}
