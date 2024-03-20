package repository

import (
	"context"

	"github.com/zikri124/mygram-api/internal/infrastructure"
	"github.com/zikri124/mygram-api/internal/model"
	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetAllComment(ctx context.Context) ([]model.CommentView, error)
	GetCommentById(ctx context.Context, commentId uint32) (*model.Comment, error)
	UpdateComment(ctx context.Context, comment *model.Comment) error
}

type commentRepositoryImpl struct {
	db infrastructure.GormPostgres
}

func NewCommentRepository(db infrastructure.GormPostgres) CommentRepository {
	return &commentRepositoryImpl{db: db}
}

func (c *commentRepositoryImpl) CreateComment(ctx context.Context, comment *model.Comment) error {
	db := c.db.GetConnection()

	err := db.
		WithContext(ctx).
		Table("comments").
		Create(&comment).
		Error

	return err
}

func (c *commentRepositoryImpl) GetAllComment(ctx context.Context) ([]model.CommentView, error) {
	db := c.db.GetConnection()
	comments := []model.CommentView{}

	err := db.
		WithContext(ctx).
		Table("comments").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email, username").Table("users").Where("deleted_at is null")
		}).
		Preload("Photo", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, title, caption, photo_url, user_id").Table("photos").Where("deleted_at is null")
		}).
		Find(&comments).
		Error

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *commentRepositoryImpl) GetCommentById(ctx context.Context, commentId uint32) (*model.Comment, error) {
	db := c.db.GetConnection()
	comment := model.Comment{}

	err := db.
		WithContext(ctx).
		Table("comments").
		Where("id = ?", commentId).
		Where("deleted_at IS NULL").
		Find(&comment).
		Error

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (c *commentRepositoryImpl) UpdateComment(ctx context.Context, comment *model.Comment) error {
	db := c.db.GetConnection()
	err := db.
		WithContext(ctx).
		Updates(&comment).
		Error

	return err
}
