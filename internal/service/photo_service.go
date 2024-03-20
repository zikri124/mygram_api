package service

import (
	"context"
	"time"

	"github.com/zikri124/mygram-api/internal/model"
	"github.com/zikri124/mygram-api/internal/repository"
)

type PhotoService interface {
	PostPhoto(ctx context.Context, photo model.Photo) (*model.PhotoRes, error)
	GetAllPhotos(ctx context.Context) ([]model.PhotoView, error)
}

type photoServiceImpl struct {
	repo repository.PhotoRepository
}

func NewPhotoService(repo repository.PhotoRepository) PhotoService {
	return &photoServiceImpl{repo: repo}
}

func (p *photoServiceImpl) PostPhoto(ctx context.Context, photo model.Photo) (*model.PhotoRes, error) {
	err := p.repo.CreatePhoto(ctx, &photo)
	if err != nil {
		return nil, err
	}

	photoRes := model.PhotoRes{}
	photoRes.ID = photo.ID
	photoRes.Caption = photo.Caption
	photoRes.PhotoUrl = photo.PhotoUrl
	photoRes.Title = photo.Title
	photoRes.UserId = photo.UserId
	photoRes.CreatedAt = time.Now()

	return &photoRes, nil
}

func (p *photoServiceImpl) GetAllPhotos(ctx context.Context) ([]model.PhotoView, error) {
	photos, err := p.repo.GetAllPhotos(ctx)
	if err != nil {
		return nil, err
	}

	return photos, nil
}
