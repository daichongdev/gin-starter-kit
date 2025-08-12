package repository

import (
	"gin-demo/database"
	"gin-demo/model"
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
