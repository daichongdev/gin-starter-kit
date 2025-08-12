package service

import (
	"context"
	"errors"
	"gin-demo/model"
	"gin-demo/pkg/auth"
	"gin-demo/pkg/logger"
	"gin-demo/repository"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// 用户注册
func (s *AuthService) Register(req *model.RegisterRequest) (*model.LoginResponse, error) {
	// 检查邮箱是否已存在
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		logger.Error("Failed to check email existence",
			logger.Err(err),
			logger.String("email", req.Email))
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password", logger.Err(err))
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Age:      req.Age,
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Error("Failed to create user",
			logger.Err(err),
			logger.String("email", req.Email))
		return nil, err
	}

	// 生成JWT Token
	token, expiresAt, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		logger.Error("Failed to generate token",
			logger.Err(err),
			logger.Uint("user_id", user.ID))
		return nil, err
	}

	logger.Info("User registered successfully",
		logger.Uint("user_id", user.ID),
		logger.String("email", user.Email))

	return &model.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: model.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		},
	}, nil
}

// 用户登录
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// 根据邮箱查找用户
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("Login attempt with non-existent email",
				logger.String("email", req.Email))
			return nil, errors.New("invalid email or password")
		}
		logger.Error("Failed to get user by email",
			logger.Err(err),
			logger.String("email", req.Email))
		return nil, err
	}

	// 验证密码
	if !auth.CheckPassword(user.Password, req.Password) {
		logger.Warn("Login attempt with wrong password",
			logger.String("email", req.Email),
			logger.Uint("user_id", user.ID))
		return nil, errors.New("invalid email or password")
	}

	// 生成JWT Token
	token, expiresAt, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		logger.Error("Failed to generate token",
			logger.Err(err),
			logger.Uint("user_id", user.ID))
		return nil, err
	}

	logger.Info("User logged in successfully",
		logger.Uint("user_id", user.ID),
		logger.String("email", user.Email))

	return &model.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: model.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		},
	}, nil
}

// 获取当前用户信息
func (s *AuthService) GetCurrentUser(ctx context.Context, userID uint) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		logger.Error("Failed to get current user",
			logger.Err(err),
			logger.Uint("user_id", userID))
		return nil, err
	}

	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}
