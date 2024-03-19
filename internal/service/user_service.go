package service

import (
	"context"
	"errors"
	"time"

	"github.com/zikri124/mygram-api/internal/model"
	"github.com/zikri124/mygram-api/internal/repository"
	"github.com/zikri124/mygram-api/pkg/helper"
)

type UserService interface {
	GetUserById(ctx context.Context, userId uint32) (*model.UserView, error)
	UserRegister(ctx context.Context, userRegData model.UserSignUp) (*model.UserView, error)
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

	userView := model.UserView{ID: user.ID, Username: user.Username, Email: user.Email, Age: age}

	return &userView, nil
}

func (u *userServiceImpl) UserRegister(ctx context.Context, userRegData model.UserSignUp) (*model.UserView, error) {
	user := model.User{}
	user.Username = userRegData.Username
	user.Email = userRegData.Email
	dobTime, err := time.Parse("2006-01-02", userRegData.DOB)
	if err != nil {
		return nil, err
	}
	user.DOB = dobTime

	hashedPass, err := helper.GenerateHash(userRegData.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPass

	userAge := helper.CountAge(user.DOB)
	if userAge <= 8 {
		return nil, errors.New("user age must above 8")
	}

	err = u.repo.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	userView := model.UserView{}
	userView.ID = user.ID
	userView.Email = user.Email
	userView.Username = user.Username
	userView.Age = userAge

	return &userView, nil
}
