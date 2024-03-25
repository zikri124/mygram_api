package service

import (
	"context"
	"time"

	"github.com/zikri124/mygram-api/internal/model"
	"github.com/zikri124/mygram-api/internal/repository"
)

type PhotoService interface {
	PostPhoto(ctx context.Context, photo model.Photo) (*model.PhotoResCreate, error)
	GetAllPhotosByUserId(ctx context.Context, userId uint32) ([]model.PhotoView, error)
	GetPhotoById(ctx context.Context, photoId uint32) (*model.PhotoView, error)
	UpdatePhoto(ctx context.Context, photo model.Photo) (*model.PhotoResUpdate, error)
	DeletePhoto(ctx context.Context, photoId uint32) error
}

type photoServiceImpl struct {
	repo repository.PhotoRepository
}

func NewPhotoService(repo repository.PhotoRepository) PhotoService {
	return &photoServiceImpl{repo: repo}
}

func (p *photoServiceImpl) PostPhoto(ctx context.Context, photo model.Photo) (*model.PhotoResCreate, error) {
	err := p.repo.CreatePhoto(ctx, &photo)
	if err != nil {
		return nil, err
	}

	photoRes := model.PhotoResCreate{}
	photoRes.ID = photo.ID
	photoRes.Caption = photo.Caption
	photoRes.PhotoUrl = photo.PhotoUrl
	photoRes.Title = photo.Title
	photoRes.UserId = photo.UserId
	photoRes.CreatedAt = time.Now()

	return &photoRes, nil
}

func (p *photoServiceImpl) GetAllPhotosByUserId(ctx context.Context, userId uint32) ([]model.PhotoView, error) {
	photos, err := p.repo.GetAllPhotosByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return photos, nil
}

func (p *photoServiceImpl) GetPhotoById(ctx context.Context, photoId uint32) (*model.PhotoView, error) {
	photo, err := p.repo.GetPhotoById(ctx, photoId)
	if err != nil {
		return nil, err
	}

	return photo, nil
}

func (p *photoServiceImpl) UpdatePhoto(ctx context.Context, photo model.Photo) (*model.PhotoResUpdate, error) {
	err := p.repo.UpdatePhoto(ctx, &photo)
	if err != nil {
		return nil, err
	}

	photoRes := model.PhotoResUpdate{}
	photoRes.ID = photo.ID
	photoRes.Caption = photo.Caption
	photoRes.PhotoUrl = photo.PhotoUrl
	photoRes.Title = photo.Title
	photoRes.UserId = photo.UserId
	photoRes.UpdatedAt = photo.UpdatedAt

	return &photoRes, nil
}

func (p *photoServiceImpl) DeletePhoto(ctx context.Context, photoId uint32) error {
	err := p.repo.DeletePhoto(ctx, photoId)

	return err
}
