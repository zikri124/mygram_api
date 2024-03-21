package repository

import (
	"context"

	"github.com/zikri124/mygram-api/internal/infrastructure"
	"github.com/zikri124/mygram-api/internal/model"
	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetAllCommentsByPhotoId(ctx context.Context, photoId uint32) ([]model.CommentView, error)
	GetCommentById(ctx context.Context, commentId uint32) (*model.Comment, error)
	UpdateComment(ctx context.Context, comment *model.Comment) error
	DeleteComment(ctx context.Context, commentId uint32) error
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

func (c *commentRepositoryImpl) GetAllCommentsByPhotoId(ctx context.Context, photoId uint32) ([]model.CommentView, error) {
	db := c.db.GetConnection()
	comments := []model.CommentView{}

	err := db.
		WithContext(ctx).
		Table("comments").
		Where("photo_id = ?", photoId).
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

func (c *commentRepositoryImpl) DeleteComment(ctx context.Context, commentId uint32) error {
	db := c.db.GetConnection()
	comment := model.Comment{ID: commentId}

	err := db.
		WithContext(ctx).
		Model(&comment).
		Delete(&comment).
		Error

	return err
}
