package controller

import (
	"gin-demo/model"
	"gin-demo/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// CreateUser 创建用户
func (uc *UserController) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("请求参数格式错误"))
		return
	}

	// 调用服务层创建用户
	user, err := uc.userService.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse("创建用户失败"))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse("用户创建成功", user))
}

// GetUser 获取用户
func (uc *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse("用户ID格式错误"))
		return
	}

	user, err := uc.userService.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse("用户不存在"))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse("获取用户信息成功", user))
}

// GetAllUsers 获取所有用户
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse("获取用户列表失败"))
		return
	}

	data := map[string]interface{}{
		"users": users,
		"count": len(users),
	}

	c.JSON(http.StatusOK, model.SuccessResponse("获取用户列表成功", data))
}
