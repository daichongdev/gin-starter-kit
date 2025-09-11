package controller

import (
	"gin-demo/model"
	"gin-demo/model/tool"
	"gin-demo/pkg/auth"
	"gin-demo/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: service.NewAuthService(),
	}
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	var req model.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, tool.ErrorResponse("请求参数无效: "+err.Error()))
		return
	}

	response, err := c.authService.Register(&req)
	if err != nil {
		if err.Error() == "email already exists" {
			ctx.JSON(http.StatusConflict, tool.ErrorResponse("邮箱已存在"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, tool.ErrorResponse("注册失败: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, tool.SuccessResponse("注册成功", response))
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, tool.ErrorResponse("请求参数无效: "+err.Error()))
		return
	}

	response, err := c.authService.Login(ctx, &req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			ctx.JSON(http.StatusUnauthorized, tool.ErrorResponse("邮箱或密码错误"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, tool.ErrorResponse("登录失败: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, tool.SuccessResponse("登录成功", response))
}

// GetProfile 获取当前用户信息
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, tool.ErrorResponse("用户未认证"))
		return
	}

	user, err := c.authService.GetCurrentUser(ctx, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, tool.ErrorResponse("获取用户信息失败: "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, tool.SuccessResponse("获取用户信息成功", user))
}

// RefreshToken 刷新Token
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, tool.ErrorResponse("用户未认证"))
		return
	}

	userEmail, exists := ctx.Get("user_email")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, tool.ErrorResponse("用户邮箱信息缺失"))
		return
	}

	// 生成新的Token
	token, expiresAt, err := auth.GenerateToken(userID.(uint), userEmail.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, tool.ErrorResponse("生成Token失败"))
		return
	}

	refreshData := gin.H{
		"token":      token,
		"expires_at": expiresAt,
	}

	ctx.JSON(http.StatusOK, tool.SuccessResponse("Token刷新成功", refreshData))
}
