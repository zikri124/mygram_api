package repository

import (
	"context"

	"github.com/zikri124/mygram-api/internal/infrastructure"
	"github.com/zikri124/mygram-api/internal/model"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
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
