package repository

import (
	"context"

	"github.com/zikri124/mygram-api/internal/infrastructure"
	"github.com/zikri124/mygram-api/internal/model"
)

type UserRepository interface {
	GetUserById(ctx context.Context, userId uint32) (model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	EditUser(ctx context.Context, user *model.User) error
}

type userRepositoryImpl struct {
	db infrastructure.GormPostgres
}

func NewUserRepository(db infrastructure.GormPostgres) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (u *userRepositoryImpl) GetUserById(ctx context.Context, userId uint32) (model.User, error) {
	db := u.db.GetConnection()

	user := model.User{}

	err := db.
		WithContext(ctx).
		Model(&user).
		Where("id = ?", userId).
		Find(&user).
		Error

	return user, err
}

func (u *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	db := u.db.GetConnection()

	user := model.User{}

	err := db.
		WithContext(ctx).
		Model(&user).
		Where("email = ?", email).
		Find(&user).
		Error

	return user, err
}

func (u *userRepositoryImpl) CreateUser(ctx context.Context, user *model.User) error {
	db := u.db.GetConnection()

	err := db.
		WithContext(ctx).
		Table("users").
		Create(&user).
		Error

	return err
}

func (u *userRepositoryImpl) EditUser(ctx context.Context, user *model.User) error {
	db := u.db.GetConnection()

	err := db.
		WithContext(ctx).
		Updates(&user).
		Error

	return err
}
