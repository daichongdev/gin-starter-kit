package model

import (
	"gin-demo/pkg/types"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"regexp"
	"time"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"` // 密码字段，JSON序列化时忽略
	Age      int    `json:"age"`
	Phone    string `json:"phone" gorm:"type:varchar(11);unique;comment:手机号码;default:''"`

	// GORM默认字段放在最后，使用自定义序列化方法
	CreatedAt types.JSONTime `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt types.JSONTime `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"comment:删除时间"`
}

func (User) TableName() string {
	return "users"
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50,alphaunicode"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=8,max=128,password"`
	Age      int    `json:"age" binding:"min=0,max=150"`
	Phone    string `json:"phone" binding:"required,len=11,numeric"`
}

// 添加自定义验证器
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("password", validatePassword)
		if err != nil {
			return
		}
	}
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// 密码必须包含大小写字母、数字和特殊字符
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"min=0"`
	Phone string `json:"phone" binding:"required,len=11"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty"`
	Age   int    `json:"age" binding:"omitempty,min=0"`
	Phone string `json:"phone" binding:"omitempty,len=11"`
}

type UserResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Age       int            `json:"age"`
	Phone     string         `json:"phone"`
	CreatedAt types.JSONTime `json:"created_at"`
	UpdatedAt types.JSONTime `json:"updated_at"`
}

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
