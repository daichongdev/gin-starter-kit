package controller

import (
	"gin-demo/model"
	"gin-demo/model/tool"
	"gin-demo/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser 创建用户
func (uc *UserController) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, tool.ErrorResponse("请求参数格式错误"))
		return
	}

	// 调用服务层创建用户
	user, err := uc.userService.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, tool.ErrorResponse("创建用户失败"))
		return
	}

	c.JSON(http.StatusCreated, tool.SuccessResponse("用户创建成功", user))
}

// GetUser 获取用户
func (uc *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, tool.ErrorResponse("用户ID格式错误"))
		return
	}

	user, err := uc.userService.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, tool.ErrorResponse("用户不存在"))
		return
	}

	c.JSON(http.StatusOK, tool.SuccessResponse("获取用户信息成功", user))
}

// GetAllUsers 获取所有用户
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, tool.ErrorResponse("获取用户列表失败"))
		return
	}

	data := map[string]interface{}{
		"users": users,
		"count": len(users),
	}

	c.JSON(http.StatusOK, tool.SuccessResponse("获取用户列表成功", data))
}

// GetUsersWithPagination 分页获取用户列表
func (uc *UserController) GetUsersWithPagination(c *gin.Context) {
	var pagination tool.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, tool.ErrorResponse("分页参数格式错误"))
		return
	}

	// 获取搜索关键词（可选）
	keyword := c.Query("keyword")

	var result *tool.PaginateResult
	var err error

	if keyword != "" {
		// 如果有搜索关键词，使用搜索分页
		result, err = uc.userService.SearchUsersWithPagination(&pagination, keyword)
	} else {
		// 普通分页查询
		result, err = uc.userService.GetUsersWithPagination(&pagination)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, tool.ErrorResponse("获取用户列表失败"))
		return
	}

	c.JSON(http.StatusOK, tool.PaginationSuccessResponse("获取用户列表成功", result.Data, result.Meta))
}

// UpdateUser 更新用户
func (uc *UserController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, tool.ErrorResponse("用户ID格式错误"))
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, tool.ErrorResponse("请求参数格式错误"))
		return
	}

	user, err := uc.userService.UpdateUser(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, tool.ErrorResponse("用户不存在或更新失败"))
		return
	}

	c.JSON(http.StatusOK, tool.SuccessResponse("用户更新成功", user))
}

// DeleteUser 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, tool.ErrorResponse("用户ID格式错误"))
		return
	}

	err = uc.userService.DeleteUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, tool.ErrorResponse("用户不存在或删除失败"))
		return
	}

	c.JSON(http.StatusOK, tool.SuccessResponse("用户删除成功", nil))
}
