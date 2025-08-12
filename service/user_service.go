package service

import (
	"gin-demo/model"
	"gin-demo/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *UserService) CreateUser(req *model.CreateUserRequest) (*model.UserResponse, error) {
	user := &model.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

func (s *UserService) GetUser(id uint) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

func (s *UserService) GetAllUsers() ([]model.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []model.UserResponse
	for _, user := range users {
		responses = append(responses, model.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		})
	}

	return responses, nil
}

func (s *UserService) UpdateUser(id uint, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	// 先获取现有用户
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Age > 0 {
		user.Age = req.Age
	}

	// 保存更新
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

func (s *UserService) DeleteUser(id uint) error {
	// 先检查用户是否存在
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(id)
}
