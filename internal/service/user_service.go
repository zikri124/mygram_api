package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zikri124/mygram-api/internal/model"
	"github.com/zikri124/mygram-api/internal/repository"
	"github.com/zikri124/mygram-api/pkg/helper"
)

type UserService interface {
	GetUserById(ctx context.Context, userId uint32) (*model.UserView, error)
	UserRegister(ctx context.Context, userRegData model.UserSignUp) (*model.UserView, error)
	UserLogin(ctx context.Context, userData model.UserSignIn) (*model.User, error)
	GenerateAccessToken(ctx context.Context, user model.User) (token string, err error)
	EditUser(ctx context.Context, userData model.User) (*model.UserView, error)
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

func (u *userServiceImpl) UserLogin(ctx context.Context, userData model.UserSignIn) (*model.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New("invalid email or password")
	}

	isValidLogin := helper.CheckPasswordHash(userData.Password, user.Password)
	if !isValidLogin {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

func (u *userServiceImpl) GenerateAccessToken(ctx context.Context, user model.User) (token string, err error) {
	now := time.Now()

	claim := model.StandardClaim{
		Jti: fmt.Sprintf("%v", time.Now().UnixNano()),
		Iss: "MyGram",
		Aud: user.Username,
		Sub: "access-token",
		Exp: uint64(now.Add(12 * time.Hour).Unix()),
		Iat: uint64(now.Unix()),
		Nbf: uint64(now.Unix()),
	}

	userClaim := model.AccessClaim{
		StandardClaim: claim,
		UserID:        uint64(user.ID),
		Username:      user.Username,
		DOB:           user.DOB,
	}

	token, err = helper.GenerateToken(userClaim)
	return
}

func (u *userServiceImpl) EditUser(ctx context.Context, user model.User) (*model.UserView, error) {
	user.UpdatedAt = time.Now()
	userFind, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if userFind.ID != 0 && user.ID != userFind.ID {
		return nil, errors.New("email already exist")
	}

	err = u.repo.EditUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	userView := model.UserView{}
	userView.ID = user.ID
	userView.Email = user.Email
	userView.Username = user.Username
	userView.Age = helper.CountAge(user.DOB)
	return &userView, nil
}
