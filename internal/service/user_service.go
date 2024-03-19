package service

import (
	"context"

	"github.com/zikri124/mygram-api/internal/model"
	"github.com/zikri124/mygram-api/internal/repository"
	"github.com/zikri124/mygram-api/pkg/helper"
)

type UserService interface {
	GetUserById(ctx context.Context, userId uint32) (*model.UserView, error)
}

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (u *userServiceImpl) GetUserById(ctx context.Context, userId uint32) (*model.UserView, error) {
	user, err := u.repo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	age := helper.CountAge(user.DOB)

	userView := model.UserView{ID: user.ID, Username: user.Username, Email: user.Email, Age: age, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}

	return &userView, nil
}
