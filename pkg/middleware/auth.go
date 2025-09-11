package middleware

import (
	"gin-demo/model/tool"
	"gin-demo/pkg/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header中获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, tool.ErrorResponse("请提供有效的JWT令牌"))
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, tool.ErrorResponse("授权头格式无效，应为: Bearer <token>"))
			c.Abort()
			return
		}

		// 验证Token
		token := parts[1]
		claims, err := auth.ValidateToken(token)
		if err != nil {
			// 根据错误类型返回不同的错误信息
			var message string
			if strings.Contains(err.Error(), "expired") {
				message = "JWT令牌已过期，请重新登录"
			} else {
				message = "无效的JWT令牌"
			}

			c.JSON(http.StatusUnauthorized, tool.ErrorResponse(message))
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}

// OptionalJWTAuthMiddleware 可选的JWT认证中间件（不强制要求认证）
func OptionalJWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				if claims, err := auth.ValidateToken(token); err == nil {
					c.Set("user_id", claims.UserID)
					c.Set("user_email", claims.Email)
				}
			}
		}
		c.Next()
	}
}
