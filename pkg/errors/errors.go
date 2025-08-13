package errors

import (
	"errors"
	"fmt"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrUserNotFound    = &AppError{Code: "USER_NOT_FOUND", Message: "用户不存在"}
	ErrEmailExists     = &AppError{Code: "EMAIL_EXISTS", Message: "邮箱已存在"}
	ErrInvalidPassword = &AppError{Code: "INVALID_PASSWORD", Message: "密码错误"}
	ErrDatabaseError   = &AppError{Code: "DATABASE_ERROR", Message: "数据库操作失败"}
	ErrInternalError   = &AppError{Code: "INTERNAL_ERROR", Message: "内部服务器错误"}
)

func NewAppError(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}
