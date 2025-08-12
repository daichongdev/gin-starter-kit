package test

import (
	"fmt"
	"gin-demo/pkg/auth"
	"testing"
)

func TestGetPwd(t *testing.T) {
	password := "123456"
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		fmt.Printf("加密失败: %v\n", err)
		t.FailNow()
	}

	fmt.Printf("原密码: %s\n", password)
	fmt.Printf("加密后: %s\n", hashedPassword)
}
