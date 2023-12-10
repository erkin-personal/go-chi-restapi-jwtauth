package services

import (
	"restapi/internal/domain/models"

	
)

type UserService struct {
	userRepo *models.UserRepository
}

func NewUserService(userRepo *models.UserRepository) *UserService {
	
	return &UserService{
		userRepo: userRepo,
	}
    
}

func (us *UserService) GetAll() ([]*models.User, error) {
	return us.userRepo.GetAll()
}

func (us *UserService) Create(user *models.User) error {
	return us.userRepo.Create(user)
}
