package repository

import (
	"gin-demo/database"
	"gin-demo/model"
	"gin-demo/model/tool"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *model.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	err := database.DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := database.DB.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) Update(user *model.User) error {
	return database.DB.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return database.DB.Delete(&model.User{}, id).Error
}

// GetAllWithPagination 分页获取用户列表 - 使用GORM Scopes优化版本
func (r *UserRepository) GetAllWithPagination(pagination *tool.PaginationRequest) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// 获取总数
	if err := database.DB.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 使用Scopes进行分页查询 - 更简洁优雅
	err := database.DB.Scopes(pagination.Paginate()).Find(&users).Error
	return users, total, err
}

// GetAllWithPaginationAndSearch 带搜索的分页查询
func (r *UserRepository) GetAllWithPaginationAndSearch(pagination *tool.PaginationRequest, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := database.DB.Model(&model.User{})

	// 如果有搜索关键词，添加搜索条件
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR email LIKE ?", searchPattern, searchPattern)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Scopes(pagination.Paginate()).Find(&users).Error
	return users, total, err
}
