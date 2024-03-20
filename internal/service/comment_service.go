package service

import (
	"context"

	"github.com/zikri124/mygram-api/internal/model"
	"github.com/zikri124/mygram-api/internal/repository"
)

type CommentService interface {
	PostComment(ctx context.Context, comment model.Comment) (*model.CreateCommentRes, error)
}

type commentServiceImpl struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentServiceImpl{repo: repo}
}

func (c *commentServiceImpl) PostComment(ctx context.Context, comment model.Comment) (*model.CreateCommentRes, error) {
	err := c.repo.CreateComment(ctx, &comment)
	if err != nil {
		return nil, err
	}

	commentRes := model.CreateCommentRes{}
	commentRes.ID = comment.ID
	commentRes.Message = comment.Message
	commentRes.PhotoId = comment.PhotoId
	commentRes.UserId = comment.UserId
	commentRes.CreatedAt = comment.CreatedAt

	return &commentRes, nil
}